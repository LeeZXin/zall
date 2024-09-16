package pullrequestsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/git/modules/service/workflowsrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.PullRequestTopic, func(data any) {
			req, ok := data.(event.PullRequestEvent)
			if ok {
				ctx, closer := xormstore.Context(context.Background())
				// 触发webhook
				hookList, err := webhookmd.ListWebhookByRepoId(ctx, req.RepoId)
				closer.Close()
				if err == nil && len(hookList) > 0 {
					for _, hook := range hookList {
						if hook.GetEvents().PullRequest {
							webhook.TriggerWebhook(hook.HookUrl, hook.Secret, &req)
						}
					}
				}
				if req.Action == event.PullRequestEventMergeAction {
					// 触发工作流
					workflowsrv.FindAndExecute(workflowsrv.FindAndExecuteWorkflowReqDTO{
						RepoId:      req.RepoId,
						RepoPath:    req.RepoPath,
						Operator:    req.Operator,
						TriggerType: workflowmd.HookTriggerType,
						Branch:      req.HeadRef,
						Source:      workflowmd.PullRequestTriggerSource,
						PrId:        req.PrId,
						PrIndex:     req.PrIndex,
					})
				}
				// 触发teamhook
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.PullRequestEventSubmitAction:
						return events.PullRequest.Submit
					case event.PullRequestEventCloseAction:
						return events.PullRequest.Close
					case event.PullRequestEventMergeAction:
						return events.PullRequest.Merge
					case event.PullRequestEventReviewAction:
						return events.PullRequest.Review
					case event.PullRequestEventAddCommentAction:
						return events.PullRequest.AddComment
					case event.PullRequestEventDeleteCommentAction:
						return events.PullRequest.DeleteComment
					default:
						return false
					}
				})
			}
		})
	})
}

// ListPullRequest 查询合并请求列表
func ListPullRequest(ctx context.Context, reqDTO ListPullRequestReqDTO) ([]PullRequestDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkAccessRepoPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	prs, totalCount, err := pullrequestmd.ListPullRequest(ctx, pullrequestmd.ListPullRequestReqDTO{
		RepoId:    reqDTO.RepoId,
		SearchKey: reqDTO.SearchKey,
		Status:    reqDTO.Status,
		PageNum:   reqDTO.PageNum,
		PageSize:  10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}
	accounts := listutil.MapNe(prs, func(t pullrequestmd.PullRequest) string {
		return t.CreateBy
	})
	userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, accounts...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}
	data := listutil.MapNe(prs, func(t pullrequestmd.PullRequest) PullRequestDTO {
		return PullRequestDTO{
			Id:             t.Id,
			RepoId:         t.RepoId,
			Target:         t.Target,
			TargetType:     t.TargetType,
			TargetCommitId: t.TargetCommitId,
			Head:           t.Head,
			HeadType:       t.HeadType,
			HeadCommitId:   t.HeadCommitId,
			PrStatus:       t.PrStatus,
			CreateBy:       userMap[t.CreateBy],
			PrTitle:        t.PrTitle,
			PrIndex:        t.PrIndex,
			CommentCount:   t.CommentCount,
			Created:        t.Created,
		}
	})
	return data, totalCount, nil
}

// SubmitPullRequest 创建合并请求
func SubmitPullRequest(ctx context.Context, reqDTO SubmitPullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, team, err := checkSubmitPullRequestPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 检查是否有已经有合并的pr
	b, err := pullrequestmd.ExistsPrByRepoIdAndRef(ctx, pullrequestmd.ExistsPrByRepoIdAndRefReqDTO{
		RepoId:     reqDTO.RepoId,
		Head:       reqDTO.Head,
		HeadType:   reqDTO.HeadType,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Status:     pullrequestmd.PrOpenStatus,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 已存在pr
	if b {
		return util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestAlreadyExists)
	}
	// 判断是否能合并
	canMerge, err := client.CanMerge(ctx, reqvo.CanMergeReq{
		RepoPath:   repo.Path,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Head:       reqDTO.Head,
		HeadType:   reqDTO.HeadType,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 不可合并
	if !canMerge {
		return util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
	}
	var (
		pr    pullrequestmd.PullRequest
		err2  error
		index int
	)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		commentCount := 0
		if reqDTO.Comment != "" {
			commentCount = 1
		}
		index, err2 = pullrequestmd.GetNextMaxIndex(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		pr, err2 = pullrequestmd.InsertPullRequest(ctx, pullrequestmd.InsertPullRequestReqDTO{
			RepoId:       reqDTO.RepoId,
			Target:       reqDTO.Target,
			TargetType:   reqDTO.TargetType,
			Head:         reqDTO.Head,
			HeadType:     reqDTO.HeadType,
			CreateBy:     reqDTO.Operator.Account,
			Title:        reqDTO.Title,
			Index:        index,
			PrStatus:     pullrequestmd.PrOpenStatus,
			CommentCount: commentCount,
		})
		if err2 != nil {
			return err2
		}
		if reqDTO.Comment != "" {
			return pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
				{
					PrId:    pr.Id,
					Action:  pullrequestmd.NewPrAction(pr.Id, pullrequestmd.PrOpenStatus),
					Account: reqDTO.Operator.Account,
				}, {
					PrId:    pr.Id,
					Action:  pullrequestmd.NewCommentAction(reqDTO.Comment),
					Account: reqDTO.Operator.Account,
				},
			})
		} else {
			return pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
				{
					PrId:    pr.Id,
					Action:  pullrequestmd.NewPrAction(pr.Id, pullrequestmd.PrOpenStatus),
					Account: reqDTO.Operator.Account,
				},
			})
		}
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 触发webhook
	notifyEvent(repo, team, reqDTO.Operator, pr, event.PullRequestEventSubmitAction)
	return nil
}

// GetPullRequest 查询合并请求
func GetPullRequest(ctx context.Context, reqDTO GetPullRequestReqDTO) (PullRequestDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return PullRequestDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkAccessRepoPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return PullRequestDTO{}, err
	}
	t, b, err := pullrequestmd.GetPullRequestByRepoIdAndIndex(ctx, reqDTO.RepoId, reqDTO.Index)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return PullRequestDTO{}, util.InternalError(err)
	}
	if !b {
		return PullRequestDTO{}, util.InvalidArgsError()
	}
	userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, t.CreateBy)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return PullRequestDTO{}, util.InternalError(err)
	}
	return PullRequestDTO{
		Id:             t.Id,
		RepoId:         t.RepoId,
		Target:         t.Target,
		TargetType:     t.TargetType,
		TargetCommitId: t.TargetCommitId,
		Head:           t.Head,
		HeadType:       t.HeadType,
		HeadCommitId:   t.HeadCommitId,
		PrStatus:       t.PrStatus,
		CreateBy:       userMap[t.CreateBy],
		PrTitle:        t.PrTitle,
		PrIndex:        t.PrIndex,
		CommentCount:   t.CommentCount,
		Created:        t.Created,
	}, nil
}

func ClosePullRequest(ctx context.Context, reqDTO ClosePullRequestReqDTO) (bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, team, err := checkSubmitPullRequestPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return false, err
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return true, nil
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := pullrequestmd.ClosePrStatus(ctx, reqDTO.PrId, pullrequestmd.PrOpenStatus, reqDTO.Operator.Account)
		if err2 != nil {
			return err2
		}
		return pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
			{
				PrId:    pr.Id,
				Action:  pullrequestmd.NewPrAction(pr.Id, pullrequestmd.PrClosedStatus),
				Account: reqDTO.Operator.Account,
			},
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false, util.InternalError(err)
	}
	// 触发webhook
	notifyEvent(repo, team, reqDTO.Operator, pr, event.PullRequestEventCloseAction)
	return false, nil
}

// CanMergePullRequest 是否可合并
func CanMergePullRequest(ctx context.Context, reqDTO CanMergePullRequestReqDTO) (CanMergePullRequestRespDTO, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CanMergePullRequestRespDTO{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, _, err := checkAccessRepoPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return CanMergePullRequestRespDTO{}, false, err
	}
	// 只允许从open
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return CanMergePullRequestRespDTO{}, true, nil
	}
	var (
		reviewCanMerge, gitCanMerge, isProtectedBranch bool

		cfg         branch.ProtectedBranchCfg
		reviewCount int
	)
	ret := CanMergePullRequestRespDTO{}
	reviewCanMerge, isProtectedBranch, cfg, reviewCount, err = detectCanMergePullRequest(ctx, pr)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ret, false, util.InternalError(err)
	}
	reviewList, err := usersrv.GetUsersNameAndAvatar(ctx, cfg.ReviewerList...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ret, false, util.InternalError(err)
	}
	ret.ProtectedBranchCfg = ProtectedBranchCfgDTO{
		PushOption:              cfg.PushOption,
		ReviewCountWhenCreatePr: cfg.ReviewCountWhenCreatePr,
		ReviewerList:            reviewList,
	}
	ret.IsProtectedBranch = isProtectedBranch
	ret.ReviewCount = reviewCount
	var info reqvo.DiffRefsResp
	info, err = client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath:   repo.Path,
		Target:     pr.Target,
		TargetType: pr.TargetType,
		Head:       pr.Head,
		HeadType:   pr.HeadType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ret, false, util.InternalError(err)
	}
	gitCanMerge = info.CanMerge
	ret.GitCanMerge = info.CanMerge
	ret.GitCommitCount = len(info.Commits)
	ret.GitConflictFiles = info.ConflictFiles
	ret.CanMerge = reviewCanMerge && gitCanMerge
	return ret, false, nil
}

// CanReviewPullRequest 是否可评审代码
func CanReviewPullRequest(ctx context.Context, reqDTO CanReviewPullRequestReqDTO) (CanReviewPullRequestRespDTO, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CanReviewPullRequestRespDTO{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, ret, statusChange, err := canReview(ctx, reqDTO.PrId, reqDTO.Operator)
	return ret, statusChange, err
}

func canReview(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, teammd.Team, CanReviewPullRequestRespDTO, bool, error) {
	// 校验权限
	pr, repo, team, err := checkAccessRepoPermByPrId(ctx, prId, operator)
	if err != nil {
		return pr, repo, team, CanReviewPullRequestRespDTO{}, false, err
	}
	// 只允许从open
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return pr, repo, team, CanReviewPullRequestRespDTO{}, true, nil
	}
	protectedBranches, err := branchmd.ListProtectedBranch(ctx, pr.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pr, repo, team, CanReviewPullRequestRespDTO{}, false, util.InternalError(err)
	}
	isProtectedBranch, protectedBranch := protectedBranches.IsProtectedBranch(pr.Head)
	reviewerList := protectedBranch.GetCfg().ReviewerList
	if reviewerList == nil {
		reviewerList = []string{}
	}
	isInReviewerList := reviewerList.Contains(operator.Account)
	// 检查是否重复提交
	review, b, err := pullrequestmd.GetReview(ctx, prId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pr, repo, team, CanReviewPullRequestRespDTO{}, false, util.InternalError(err)
	}
	hasAgree := b && review.ReviewStatus == pullrequestmd.AgreeReviewStatus
	/*
		权限：本人拥有在该仓库"处理合并请求（handlePullRequest）"的权限
		是否可评审取决于：
		1、未重复提交
		2、当不是保护分支， 可评审
		3、否则
			当指定了评审白名单，则判断自己是否在白名单里面，若在，可评审 若不在 不可评审
			当未指定白名单，则可评审
	*/
	return pr, repo, team, CanReviewPullRequestRespDTO{
		CanReview:         !hasAgree && (!isProtectedBranch || len(reviewerList) == 0 || isInReviewerList),
		IsProtectedBranch: isProtectedBranch,
		ReviewerList:      reviewerList,
		IsInReviewerList:  isInReviewerList,
		HasAgree:          hasAgree,
	}, false, nil
}

// MergePullRequest 提交合并代码
func MergePullRequest(ctx context.Context, reqDTO MergePullRequestReqDTO) (bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, team, err := checkSubmitPullRequestPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return false, err
	}
	// 只允许从open -> merged
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return true, nil
	}
	canMerge, _, _, _, err := detectCanMergePullRequest(ctx, pr)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false, util.InternalError(err)
	}
	if !canMerge {
		return false, util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestReviewerCountLowerThanCfg)
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		return mergeWithTx(ctx, pr, repo, reqDTO.Operator)
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return false, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return false, util.InternalError(err)
	}
	// 触发webhook和工作流
	notifyEvent(repo, team, reqDTO.Operator, pr, event.PullRequestEventMergeAction)
	return false, nil
}

func detectCanMergePullRequest(ctx context.Context, pr pullrequestmd.PullRequest) (bool, bool, branch.ProtectedBranchCfg, int, error) {
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, pr.RepoId, pr.Head)
	if err != nil {
		return false, false, cfg, 0, err
	}
	if isProtectedBranch {
		// 检查评审配置 评审者数量大于0
		if cfg.ReviewCountWhenCreatePr > 0 {
			var reviewCount int64
			reviewCount, err = pullrequestmd.CountReview(ctx, pr.Id, pullrequestmd.AgreeReviewStatus)
			if err != nil {
				return false, true, cfg, 0, err
			}
			// 小于配置数量 不可合并
			return int(reviewCount) >= cfg.ReviewCountWhenCreatePr, true, cfg, int(reviewCount), nil
		}
	}
	return true, isProtectedBranch, cfg, 0, nil
}

func mergeWithTx(ctx context.Context, pr pullrequestmd.PullRequest, repo repomd.Repo, operator apisession.UserInfo) error {
	err := pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
		{
			PrId:    pr.Id,
			Action:  pullrequestmd.NewPrAction(pr.Id, pullrequestmd.PrMergedStatus),
			Account: operator.Account,
		},
	})
	if err != nil {
		return err
	}
	info, err := client.Merge(ctx, reqvo.MergeReq{
		RepoPath: repo.Path,
		Target:   pr.Target,
		Head:     pr.Head,
		MergeOpts: struct {
			RepoId        int64  `json:"repoId"`
			PrId          int64  `json:"prId"`
			PusherAccount string `json:"pusherAccount"`
			PusherName    string `json:"pusherName"`
			PusherEmail   string `json:"pusherEmail"`
			Message       string `json:"message"`
			AppUrl        string `json:"appUrl"`
		}{
			RepoId:        repo.Id,
			PrId:          pr.Id,
			PusherAccount: operator.Account,
			PusherName:    operator.Name,
			PusherEmail:   operator.Email,
			Message:       fmt.Sprintf("merge %s from %s with #%d", pr.Head, pr.Target, pr.PrIndex),
		},
	})
	if err != nil {
		return err
	}
	_, err = pullrequestmd.MergePrStatus(
		ctx,
		pr.Id,
		pullrequestmd.PrOpenStatus,
		info.TargetCommit.CommitId,
		info.HeadCommit.CommitId,
		operator.Account,
	)
	return err
}

func AgreeReviewPullRequest(ctx context.Context, reqDTO AgreeReviewPullRequestReqDTO) (bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pr, repo, team, review, statusChange, err := canReview(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil || statusChange {
		return statusChange, err
	}
	if !review.CanReview {
		return false, util.InvalidArgsError()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		action, err := pullrequestmd.InsertReview(ctx, pullrequestmd.InsertReviewReqDTO{
			PrId:     reqDTO.PrId,
			Status:   pullrequestmd.AgreeReviewStatus,
			Reviewer: reqDTO.Operator.Account,
		})
		if err != nil {
			return err
		}
		return pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
			{
				PrId:    reqDTO.PrId,
				Action:  pullrequestmd.NewReviewAction(action.Id),
				Account: reqDTO.Operator.Account,
			},
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false, util.InternalError(err)
	}
	// 触发webhook
	notifyEvent(repo, team, reqDTO.Operator, pr, event.PullRequestEventReviewAction)
	return false, nil
}

// ListTimeline 展示时间轴
func ListTimeline(ctx context.Context, reqDTO ListTimelineReqDTO) ([]TimelineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, err := checkAccessRepoPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	timelines, err := pullrequestmd.ListTimeline(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	// 查找头像和姓名
	accountList := listutil.MapNe(timelines, func(t pullrequestmd.Timeline) string {
		return t.Account
	})
	users, err := usermd.ListUserByAccounts(ctx, listutil.Distinct(accountList...), []string{"account", "avatar_url", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	avatarMap := make(map[string]usermd.User, len(users))
	for _, user := range users {
		avatarMap[user.Account] = user
	}
	return listutil.Map(timelines, func(t pullrequestmd.Timeline) (TimelineDTO, error) {
		ret := TimelineDTO{
			Id:        t.Id,
			PrId:      t.PrId,
			Account:   t.Account,
			Created:   t.Created,
			AvatarUrl: avatarMap[t.Account].AvatarUrl,
			Name:      avatarMap[t.Account].Name,
		}
		if t.Action == nil {
			ret.Action = pullrequestmd.Action{}
		} else {
			ret.Action = *t.Action
		}
		return ret, nil
	})
}

// AddComment 添加评论
func AddComment(ctx context.Context, reqDTO AddCommentReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		pr pullrequestmd.PullRequest
	)
	pr, repo, team, err := checkAddCommentPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return err
	}
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return util.InvalidArgsError()
	}
	var action pullrequestmd.Action
	if reqDTO.HasReply {
		var (
			timeline pullrequestmd.Timeline
			b        bool
		)
		timeline, b, err = pullrequestmd.GetTimelineById(ctx, reqDTO.ReplyFrom)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b || timeline.PrId != reqDTO.PrId ||
			timeline.Action == nil ||
			(timeline.Action.ActionType != pullrequestmd.CommentType &&
				timeline.Action.ActionType != pullrequestmd.ReplyType) {
			return util.InvalidArgsError()
		}
		switch timeline.Action.ActionType {
		case pullrequestmd.CommentType:
			if timeline.Action.Comment == nil {
				return util.InvalidArgsError()
			}
			action = pullrequestmd.NewReplyAction(reqDTO.ReplyFrom, timeline.Account, timeline.Action.Comment.Comment, reqDTO.Comment)
		case pullrequestmd.ReplyType:
			if timeline.Action.Reply == nil {
				return util.InvalidArgsError()
			}
			action = pullrequestmd.NewReplyAction(reqDTO.ReplyFrom, timeline.Account, timeline.Action.Reply.ReplyComment, reqDTO.Comment)
		}
	} else {
		action = pullrequestmd.NewCommentAction(reqDTO.Comment)
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err2 := pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
			{
				PrId:    reqDTO.PrId,
				Action:  action,
				Account: reqDTO.Operator.Account,
			},
		})
		if err2 != nil {
			return err2
		}
		_, err2 = pullrequestmd.IncrCommentCount(ctx, reqDTO.PrId)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, pr, event.PullRequestEventAddCommentAction)
	return nil
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, reqDTO DeleteCommentReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	timeline, b, err := pullrequestmd.GetTimelineById(ctx, reqDTO.CommentId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b ||
		timeline.Account != reqDTO.Operator.Account ||
		timeline.Action == nil ||
		!timeline.Action.ActionType.IsRelatedToComment() {
		return util.InvalidArgsError()
	}
	pr, repo, team, err := checkAccessRepoPermByPrId(ctx, timeline.PrId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := pullrequestmd.DeleteTimelineById(ctx, reqDTO.CommentId)
		if err2 != nil {
			return err2
		}
		_, err2 = pullrequestmd.DecrCommentCount(ctx, timeline.PrId)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, pr, event.PullRequestEventDeleteCommentAction)
	return nil
}

// ListReview 评审记录
func ListReview(ctx context.Context, reqDTO ListReviewReqDTO) ([]ReviewDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, err := checkAccessRepoPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	reviews, err := pullrequestmd.ListReview(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	accounts := listutil.MapNe(reviews, func(t pullrequestmd.Review) string {
		return t.Reviewer
	})
	userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, accounts...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(reviews, func(t pullrequestmd.Review) (ReviewDTO, error) {
		return ReviewDTO{
			Id:           t.Id,
			Reviewer:     t.Reviewer,
			AvatarUrl:    userMap[t.Reviewer].AvatarUrl,
			Name:         userMap[t.Reviewer].Name,
			ReviewStatus: t.ReviewStatus,
			Updated:      t.Updated,
		}, nil
	})
}

// checkSubmitPullRequestPermByPrId 校验权限
func checkSubmitPullRequestPermByPrId(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, teammd.Team, error) {
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pullrequestmd.PullRequest{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return pullrequestmd.PullRequest{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkSubmitPullRequestPermByRepoId(ctx, pr.RepoId, operator)
	return pr, repo, team, err
}

// checkAccessRepoPermByPrId 校验权限
func checkAccessRepoPermByPrId(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, teammd.Team, error) {
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pullrequestmd.PullRequest{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return pullrequestmd.PullRequest{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkAccessRepoPermByRepoId(ctx, pr.RepoId, operator)
	return pr, repo, team, err
}

// checkAddCommentPermByPrId 校验权限
func checkAddCommentPermByPrId(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, teammd.Team, error) {
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pullrequestmd.PullRequest{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return pullrequestmd.PullRequest{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkAddCommentPermByRepoId(ctx, pr.RepoId, operator)
	return pr, repo, team, err
}

// checkSubmitPullRequestPermByRepoId 校验权限
func checkSubmitPullRequestPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return repo, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, team, util.InternalError(err)
	}
	if !b {
		return repo, team, util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.GetRepoPerm(repoId).CanSubmitPullRequest {
		return repo, team, util.UnauthorizedError()
	}
	return repo, team, nil
}

// checkAccessRepoPermByRepoId 校验权限
func checkAccessRepoPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return repo, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, team, util.InternalError(err)
	}
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanAccessRepo) {
		return repo, team, nil
	}
	return repo, team, util.UnauthorizedError()
}

// checkAddCommentPermByRepoId 校验权限
func checkAddCommentPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return repo, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, team, util.InternalError(err)
	}
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanAddCommentInPullRequest) {
		return repo, team, nil
	}
	return repo, team, util.UnauthorizedError()
}

func notifyEvent(repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, pr pullrequestmd.PullRequest, action event.PullRequestEventAction) {
	initPsub()
	psub.Publish(event.PullRequestTopic, event.PullRequestEvent{
		PrId:          pr.Id,
		PrTitle:       pr.PrTitle,
		PrIndex:       pr.PrIndex,
		TargetRef:     pr.Target,
		TargetRefType: pr.TargetType.String(),
		HeadRef:       pr.Head,
		HeadRefType:   pr.HeadType.String(),
		Action:        action,
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoPath: repo.Path,
			RepoId:   repo.Id,
			RepoName: repo.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
		},
	})
}
