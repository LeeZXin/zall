package pullrequestsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/oplogsrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/eventbus"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct {
}

// GetStats 获取统计详情
func (s *outerImpl) GetStats(ctx context.Context, reqDTO GetStatsReqDTO) (GetStatsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return GetStatsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return GetStatsRespDTO{}, err
	}
	rets, err := pullrequestmd.GroupByPrStatus(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return GetStatsRespDTO{}, err
	}
	var (
		totalCount, openCount, closedCount, mergedCount int64
	)
	for _, ret := range rets {
		totalCount += ret.TotalCount
		switch ret.PrStatus {
		case pullrequestmd.PrOpenStatus:
			openCount += ret.TotalCount
		case pullrequestmd.PrClosedStatus:
			closedCount += ret.TotalCount
		case pullrequestmd.PrMergedStatus:
			mergedCount += ret.TotalCount
		}
	}
	return GetStatsRespDTO{
		TotalCount:  totalCount,
		OpenCount:   openCount,
		ClosedCount: closedCount,
		MergedCount: mergedCount,
	}, nil
}

// ListPullRequest 查询合并请求列表
func (s *outerImpl) ListPullRequest(ctx context.Context, reqDTO ListPullRequestReqDTO) ([]PullRequestDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	requests, totalCount, err := pullrequestmd.ListPullRequest(ctx, pullrequestmd.ListPullRequestReqDTO{
		RepoId:    reqDTO.RepoId,
		SearchKey: reqDTO.SearchKey,
		Page2Req:  reqDTO.Page2Req,
		Status:    reqDTO.Status,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, err
	}
	data, _ := listutil.Map(requests, pr2Dto)
	return data, totalCount, nil
}

func pr2Dto(t pullrequestmd.PullRequest) (PullRequestDTO, error) {
	ret := PullRequestDTO{
		Id:             t.Id,
		RepoId:         t.RepoId,
		Target:         t.Target,
		TargetType:     t.TargetType,
		TargetCommitId: t.TargetCommitId,
		Head:           t.Head,
		HeadType:       t.HeadType,
		HeadCommitId:   t.HeadCommitId,
		PrStatus:       t.PrStatus,
		CreateBy:       t.CreateBy,
		CloseBy:        t.CloseBy,
		MergeBy:        t.MergeBy,
		PrTitle:        t.PrTitle,
		CommentCount:   t.CommentCount,
		Created:        t.Created,
		Closed:         t.Closed,
		Merged:         t.Merged,
	}
	return ret, nil
}

// SubmitPullRequest 创建合并请求
func (s *outerImpl) SubmitPullRequest(ctx context.Context, reqDTO SubmitPullRequestReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
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
		pr   pullrequestmd.PullRequest
		err2 error
	)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		commentCount := 0
		if reqDTO.Comment != "" {
			commentCount = 1
		}
		pr, err2 = pullrequestmd.InsertPullRequest(ctx, pullrequestmd.InsertPullRequestReqDTO{
			RepoId:       reqDTO.RepoId,
			Target:       reqDTO.Target,
			TargetType:   reqDTO.TargetType,
			Head:         reqDTO.Head,
			HeadType:     reqDTO.HeadType,
			CreateBy:     reqDTO.Operator.Account,
			Title:        reqDTO.Title,
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
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.PullRequestSrvKeysVO.SubmitPullRequest, pr.Id),
		Req:      reqDTO,
	})
	// 触发webhook
	notifyEventBus(repo, reqDTO.Operator, pr, webhook.PrSubmitAction)
	return nil
}

// GetPullRequest 查询合并请求
func (*outerImpl) GetPullRequest(ctx context.Context, reqDTO GetPullRequestReqDTO) (PullRequestDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return PullRequestDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, _, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return PullRequestDTO{}, err
	}
	return pr2Dto(pr)
}

func (*outerImpl) ClosePullRequest(ctx context.Context, reqDTO ClosePullRequestReqDTO) (statusChange bool, err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		statusChange = true
		return
	}
	var (
		b    bool
		err2 error
	)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err2 = pullrequestmd.ClosePrStatus(ctx, reqDTO.PrId, pullrequestmd.PrOpenStatus, reqDTO.Operator.Account)
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
		err = util.InternalError(err)
		return
	}
	if b {
		// 插入日志
		oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
			RepoId:   repo.Id,
			Operator: reqDTO.Operator.Account,
			Log:      oplogsrv.FormatI18n(i18n.PullRequestSrvKeysVO.ClosePullRequest, pr.Id),
			Req:      reqDTO,
		})
		// 触发webhook
		notifyEventBus(repo, reqDTO.Operator, pr, webhook.PrCloseAction)
	}
	return
}

// CanMergePullRequest 是否可合并
func (s *outerImpl) CanMergePullRequest(ctx context.Context, reqDTO CanMergePullRequestReqDTO) (CanMergePullRequestRespDTO, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CanMergePullRequestRespDTO{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
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
	reviewCanMerge, isProtectedBranch, cfg, reviewCount, err = s.detectCanMergePullRequest(ctx, pr)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ret, false, util.InternalError(err)
	}
	ret.ProtectedBranchCfg = cfg
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
func (s *outerImpl) CanReviewPullRequest(ctx context.Context, reqDTO CanReviewPullRequestReqDTO) (CanReviewPullRequestRespDTO, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CanReviewPullRequestRespDTO{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, ret, statusChange, err := s.canReview(ctx, reqDTO.PrId, reqDTO.Operator)
	return ret, statusChange, err
}

func (s *outerImpl) canReview(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, CanReviewPullRequestRespDTO, bool, error) {
	// 校验权限
	pr, repo, err := checkPerm(ctx, prId, operator)
	if err != nil {
		return pr, repo, CanReviewPullRequestRespDTO{}, false, err
	}
	// 只允许从open
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return pr, repo, CanReviewPullRequestRespDTO{}, true, nil
	}
	protectedBranches, err := branchmd.ListProtectedBranch(ctx, pr.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pr, repo, CanReviewPullRequestRespDTO{}, false, util.InternalError(err)
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
		return pr, repo, CanReviewPullRequestRespDTO{}, false, util.InternalError(err)
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
	canReview := !hasAgree && (!isProtectedBranch || len(reviewerList) == 0 || isInReviewerList)
	return pr, repo, CanReviewPullRequestRespDTO{
		CanReview:         canReview,
		IsProtectedBranch: isProtectedBranch,
		ReviewerList:      reviewerList,
		IsInReviewerList:  isInReviewerList,
		HasAgree:          hasAgree,
	}, false, nil
}

// MergePullRequest 提交合并代码
func (s *outerImpl) MergePullRequest(ctx context.Context, reqDTO MergePullRequestReqDTO) (statusChange bool, err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PullRequestSrvKeysVO.MergePullRequest),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return
	}
	// 只允许从open -> merged
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		statusChange = true
		return
	}
	var canMerge bool
	canMerge, _, _, _, err = s.detectCanMergePullRequest(ctx, pr)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !canMerge {
		err = util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestReviewerCountLowerThanCfg)
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		return s.mergeWithTx(ctx, pr, repo, reqDTO.Operator)
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return
		}
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.PullRequestSrvKeysVO.MergePullRequest, pr.Id),
		Req:      reqDTO,
	})
	// 触发webhook和工作流
	notifyEventBus(repo, reqDTO.Operator, pr, webhook.PrMergeAction)
	return
}

func (*outerImpl) detectCanMergePullRequest(ctx context.Context, pr pullrequestmd.PullRequest) (bool, bool, branch.ProtectedBranchCfg, int, error) {
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

func (s *outerImpl) mergeWithTx(ctx context.Context, pr pullrequestmd.PullRequest, repo repomd.Repo, operator apisession.UserInfo) error {
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
			PusherEmail   string `json:"pusherEmail"`
			Message       string `json:"message"`
			AppUrl        string `json:"appUrl"`
		}{
			RepoId:        repo.Id,
			PrId:          pr.Id,
			PusherAccount: operator.Account,
			PusherEmail:   operator.Email,
			Message:       fmt.Sprintf("merge %s from %s", pr.Head, pr.Target),
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

func (s *outerImpl) AgreeReviewPullRequest(ctx context.Context, reqDTO AgreeReviewPullRequestReqDTO) (bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pr, repo, review, statusChange, err := s.canReview(ctx, reqDTO.PrId, reqDTO.Operator)
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
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.PullRequestSrvKeysVO.ReviewPullRequest, pr.Id),
		Req:      reqDTO,
	})
	// 触发webhook
	notifyEventBus(repo, reqDTO.Operator, pr, webhook.PrReviewAction)
	return false, nil
}

// ListTimeline 展示时间轴
func (*outerImpl) ListTimeline(ctx context.Context, reqDTO ListTimelineReqDTO) ([]TimelineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, _, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator); err != nil {
		return nil, err
	}
	timelines, err := pullrequestmd.ListTimeline(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(timelines, func(t pullrequestmd.Timeline) (TimelineDTO, error) {
		ret := TimelineDTO{
			Id:      t.Id,
			PrId:    t.PrId,
			Account: t.Account,
			Created: t.Created,
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
func (*outerImpl) AddComment(ctx context.Context, reqDTO AddCommentReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		pr pullrequestmd.PullRequest
	)
	pr, _, err = checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return
	}
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		err = util.InvalidArgsError()
		return
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
			err = util.InternalError(err)
			return
		}
		if !b || timeline.PrId != reqDTO.PrId ||
			timeline.Action == nil ||
			(timeline.Action.ActionType != pullrequestmd.CommentType &&
				timeline.Action.ActionType != pullrequestmd.ReplyType) {
			err = util.InvalidArgsError()
			return
		}
		switch timeline.Action.ActionType {
		case pullrequestmd.CommentType:
			if timeline.Action.Comment == nil {
				err = util.InvalidArgsError()
				return
			}
			action = pullrequestmd.NewReplyAction(reqDTO.ReplyFrom, timeline.Account, timeline.Action.Comment.Comment, reqDTO.Comment)
		case pullrequestmd.ReplyType:
			if timeline.Action.Reply == nil {
				err = util.InvalidArgsError()
				return
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
		err = util.InternalError(err)
		return
	}
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   pr.RepoId,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.PullRequestSrvKeysVO.AddComment, reqDTO.Comment),
		Req:      reqDTO,
	})
	return
}

// DeleteComment 删除评论
func (*outerImpl) DeleteComment(ctx context.Context, reqDTO DeleteCommentReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		timeline pullrequestmd.Timeline
		b        bool
	)
	timeline, b, err = pullrequestmd.GetTimelineById(ctx, reqDTO.CommentId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b ||
		timeline.Account != reqDTO.Operator.Account ||
		timeline.Action == nil ||
		!timeline.Action.ActionType.IsRelatedToComment() {
		err = util.InvalidArgsError()
		return
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
		err = util.InternalError(err)
		return
	}
	// 插入日志
	{
		ctx2, closer2 := xormstore.Context(ctx)
		defer closer2.Close()
		pr, b, err := pullrequestmd.GetPullRequestById(ctx2, timeline.PrId)
		if err != nil {
			logger.Logger.WithContext(ctx2).Error(err)
		} else if b {
			oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
				RepoId:   pr.RepoId,
				Operator: reqDTO.Operator.Account,
				Log:      oplogsrv.FormatI18n(i18n.PullRequestSrvKeysVO.DeleteComment, timeline.Action.GetCommentText()),
				Req:      reqDTO,
			})
		}
	}
	return
}

// ListReview 评审记录
func (*outerImpl) ListReview(ctx context.Context, reqDTO ListReviewReqDTO) ([]ReviewDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, _, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator); err != nil {
		return nil, err
	}
	reviews, err := pullrequestmd.ListReview(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(reviews, func(t pullrequestmd.Review) (ReviewDTO, error) {
		return ReviewDTO{
			Id:           t.Id,
			Reviewer:     t.Reviewer,
			ReviewStatus: t.ReviewStatus,
			Updated:      t.Updated,
		}, nil
	})
}

// checkPerm 校验权限
func checkPerm(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.Repo, error) {
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pullrequestmd.PullRequest{}, repomd.Repo{}, util.InternalError(err)
	}
	if !b {
		return pullrequestmd.PullRequest{}, repomd.Repo{}, util.InvalidArgsError()
	}
	repo, err := checkPermByRepoId(ctx, pr.RepoId, operator)
	return pr, repo, err
}

// checkPermByRepoId 校验权限
func checkPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, util.InternalError(err)
	}
	if !b {
		return repo, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return repo, nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return repo, util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.GetRepoPerm(repoId).CanSubmitPullRequest {
		return repo, util.UnauthorizedError()
	}
	return repo, nil
}

func notifyEventBus(repo repomd.Repo, operator apisession.UserInfo, pr pullrequestmd.PullRequest, action webhook.PullRequestAction) {
	psub.Publish(eventbus.PullRequestEventTopic, eventbus.PullRequestEvent{
		PrId:      pr.Id,
		PrTitle:   pr.PrTitle,
		Action:    string(action),
		RepoId:    repo.Id,
		RepoPath:  repo.Path,
		RepoName:  repo.Name,
		Account:   operator.Account,
		Ref:       pr.Head,
		EventTime: time.Now(),
	})
}
