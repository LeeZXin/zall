package pullrequestsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
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
func (s *outerImpl) SubmitPullRequest(ctx context.Context, reqDTO SubmitPullRequestReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PullRequestSrvKeysVO.SubmitPullRequest),
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
	repo, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return
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
		err = util.InternalError(err)
		return
	}
	// 已存在pr
	if b {
		err = util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestAlreadyExists)
		return
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
		err = util.InternalError(err)
		return
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
		err = util.InternalError(err)
		return
	}
	// 触发webhook
	triggerWebhook(repo, reqDTO.Operator, pr, "submit")
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

func (*outerImpl) ClosePullRequest(ctx context.Context, reqDTO ClosePullRequestReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PullRequestSrvKeysVO.ClosePullRequest),
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
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		err = util.InvalidArgsError()
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
		// 触发webhook
		triggerWebhook(repo, reqDTO.Operator, pr, "close")
	}
	return
}

// CanMergePullRequest 是否可合并
func (s *outerImpl) CanMergePullRequest(ctx context.Context, reqDTO CanMergePullRequestReqDTO) (CanMergePullRequestRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CanMergePullRequestRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return CanMergePullRequestRespDTO{}, err
	}
	// 只允许从open
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		return CanMergePullRequestRespDTO{}, util.InvalidArgsError()
	}
	var (
		canMerge, isProtectedBranch bool
		cfg                         branchmd.ProtectedBranchCfg
		reviewCount                 int
	)
	ret := CanMergePullRequestRespDTO{}
	canMerge, isProtectedBranch, cfg, reviewCount, err = s.detectCanMergePullRequest(ctx, pr, reqDTO.Operator)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ret, util.InternalError(err)
	}
	ret.ReviewerList = cfg.ReviewerList
	ret.DirectPushList = cfg.DirectPushList
	ret.ReviewCountWhenCreatePr = cfg.ReviewCountWhenCreatePr
	ret.ReviewCount = reviewCount
	ret.IsProtectedBranch = isProtectedBranch
	if canMerge {
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
			return ret, util.InternalError(err)
		}
		canMerge = info.CanMerge
		ret.GitCanMerge = info.CanMerge
		ret.GitCommitCount = len(info.Commits)
		ret.GitConflictFiles = info.ConflictFiles
	}
	ret.CanMerge = canMerge
	return ret, nil
}

// MergePullRequest 提交合并代码
func (s *outerImpl) MergePullRequest(ctx context.Context, reqDTO MergePullRequestReqDTO) (err error) {
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
		err = util.InvalidArgsError()
		return
	}
	var canMerge bool
	canMerge, _, _, _, err = s.detectCanMergePullRequest(ctx, pr, reqDTO.Operator)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !canMerge {
		err = util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestReviewerCountLowerThanCfg)
		return
	}
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
		err = util.InternalError(err)
		return
	}
	// 不可合并
	if !info.CanMerge {
		err = util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
		return
	}
	var merged bool
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		merged, err2 = s.mergeWithTx(ctx, pr, info, repo, reqDTO.Operator, fmt.Sprintf("merge %s from %s", info.Head, info.Target))
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if merged {
		// 触发webhook
		triggerWebhook(repo, reqDTO.Operator, pr, "merge")
	}
	return nil
}

func (*outerImpl) detectCanMergePullRequest(ctx context.Context, pr pullrequestmd.PullRequest, operator apisession.UserInfo) (bool, bool, branchmd.ProtectedBranchCfg, int, error) {
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, pr.Id, pr.Head)
	if err != nil {
		return false, false, cfg, 0, err
	}
	if isProtectedBranch {
		// 判断是否可直接推送
		contains, _ := listutil.Contains(cfg.DirectPushList, func(account string) (bool, error) {
			return account == operator.Account, nil
		})
		if !contains {
			// 检查评审配置 评审者数量大于0
			if cfg.ReviewCountWhenCreatePr > 0 {
				var reviewCount int64
				reviewCount, err = pullrequestmd.CountReview(ctx, pr.Id, pullrequestmd.AgreeMergeStatus)
				if err != nil {
					return false, true, cfg, 0, err
				}
				// 小于配置数量 不可合并
				return int(reviewCount) > cfg.ReviewCountWhenCreatePr, true, cfg, int(reviewCount), nil
			}
		}
	}
	return true, isProtectedBranch, cfg, 0, nil
}

func (s *outerImpl) mergeWithTx(ctx context.Context, pr pullrequestmd.PullRequest, info reqvo.DiffRefsResp, repo repomd.Repo, operator apisession.UserInfo, message string) (bool, error) {
	b, err := pullrequestmd.MergePrStatus(
		ctx,
		pr.Id,
		pullrequestmd.PrOpenStatus,
		info.TargetCommit.CommitId,
		info.HeadCommit.CommitId,
		operator.Account,
	)
	if err != nil {
		return false, err
	}
	if b {
		err = pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
			{
				PrId:    pr.Id,
				Action:  pullrequestmd.NewPrAction(pr.Id, pullrequestmd.PrMergedStatus),
				Account: operator.Account,
			},
		})
		if err != nil {
			return false, err
		}
		err = client.Merge(ctx, reqvo.MergeReq{
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
				Message:       message,
			},
		})
		if err != nil {
			return false, err
		}
	}
	return b, nil
}

func (*outerImpl) ReviewPullRequest(ctx context.Context, reqDTO ReviewPullRequestReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PullRequestSrvKeysVO.ReviewPullRequest),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查是否重复提交
	_, b, err := pullrequestmd.GetReview(ctx, reqDTO.PrId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.NewBizErr(apicode.DataAlreadyExistsCode, i18n.ReviewAlreadyExists)
		return
	}
	// 检查评审者是否有访问代码的权限
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	repo, b, err := repomd.GetByRepoId(ctx, pr.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, reqDTO.Operator.Account)
	if !b {
		err = util.UnauthorizedError()
		return
	}
	if !p.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo {
		err = util.UnauthorizedError()
		return
	}
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, repo.Id, pr.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if isProtectedBranch {
		// 看看是否在评审名单里面 如果设置了评审名单
		if len(cfg.ReviewerList) > 0 {
			contains, _ := listutil.Contains(cfg.ReviewerList, func(account string) (bool, error) {
				return account == reqDTO.Operator.Account, nil
			})
			if !contains {
				err = util.UnauthorizedError()
				return
			}
		}
	}
	err = pullrequestmd.InsertReview(ctx, pullrequestmd.InsertReviewReqDTO{
		PrId:      reqDTO.PrId,
		ReviewMsg: reqDTO.ReviewMsg,
		Status:    reqDTO.Status,
		Reviewer:  reqDTO.Operator.Account,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 触发webhook
	triggerWebhook(repo, reqDTO.Operator, pr, "review")
	return
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
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PullRequestSrvKeysVO.AddComment),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
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
	return
}

// DeleteComment 删除评论
func (*outerImpl) DeleteComment(ctx context.Context, reqDTO DeleteCommentReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PullRequestSrvKeysVO.DeleteComment),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
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
	if !b || timeline.Account != reqDTO.Operator.Account ||
		timeline.Action == nil ||
		(timeline.Action.ActionType != pullrequestmd.CommentType &&
			timeline.Action.ActionType != pullrequestmd.ReplyType) {
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
	return
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
	repo, b := reposrv.Inner.GetByRepoId(ctx, repoId)
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return repo, util.UnauthorizedError()
	}
	if !p.PermDetail.GetRepoPerm(repoId).CanHandlePullRequest {
		return repo, util.UnauthorizedError()
	}
	return repo, nil
}

func triggerWebhook(repo repomd.Repo, operator apisession.UserInfo, pr pullrequestmd.PullRequest, actionType string) {
	go func() {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		// 触发webhook
		hookList, err := webhookmd.ListWebhook(ctx, repo.Id, webhookmd.PullRequestHook)
		if err == nil {
			req := webhook.PullRequestHook{
				PrId:       pr.Id,
				RepoId:     repo.Id,
				RepoName:   repo.Name,
				TargetRef:  pr.Target,
				HeadRef:    pr.Head,
				EventTime:  time.Now().UnixMilli(),
				ActionType: actionType,
				Operator: git.User{
					Account: operator.Account,
					Email:   operator.Email,
				},
			}
			for _, hook := range hookList {
				webhook.TriggerPrHook(hook.HookUrl, hook.HttpHeaders, req)
			}
		} else {
			logger.Logger.Error(err)
		}
	}()

}
