package pullrequestsrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
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
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"time"
)

type outerImpl struct {
}

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
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return
	}
	// 检查是否有已经有合并的pr
	b, err := pullrequestmd.ExistsOpenStatusPrByRepoIdAndRef(ctx, pullrequestmd.ExistsOpenStatusPrByRepoIdAndRefReqDTO{
		RepoId: reqDTO.RepoId,
		Head:   reqDTO.Head,
		Target: reqDTO.Target,
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
	info, err := client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath: repo.Path,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
	}, repo.NodeId)
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 不可合并
	if !info.CanMerge {
		return util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
	}
	pr, err := pullrequestmd.InsertPullRequest(ctx, pullrequestmd.InsertPullRequestReqDTO{
		RepoId:   reqDTO.RepoId,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		CreateBy: reqDTO.Operator.Account,
		PrStatus: pullrequestmd.PrOpenStatus,
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
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		err = util.InvalidArgsError()
		return
	}
	_, err = pullrequestmd.UpdatePrStatus(ctx, reqDTO.Id, pullrequestmd.PrOpenStatus, pullrequestmd.PrClosedStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 触发webhook
	triggerWebhook(repo, reqDTO.Operator, pr, "close")
	return
}

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
		return err
	}
	gitCfg, b := cfgsrv.Inner.GetGitCfg(ctx)
	if !b {
		err = util.InternalError(errors.New("can not get git config"))
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	pr, repo, err := checkPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 只允许从open -> closed
	if pr.PrStatus != pullrequestmd.PrOpenStatus {
		err = util.InvalidArgsError()
		return
	}
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, pr.RepoId, pr.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if isProtectedBranch {
		// 检查评审配置 评审者数量大于0
		if cfg.ReviewCountWhenCreatePr > 0 {
			var reviewCount int64
			reviewCount, err = pullrequestmd.CountReview(ctx, reqDTO.Id, pullrequestmd.AgreeMergeStatus)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				err = util.InternalError(err)
				return
			}
			// 小于配置数量 不可合并
			if int(reviewCount) < cfg.ReviewCountWhenCreatePr {
				err = util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestReviewerCountLowerThanCfg)
				return
			}
		}
	}
	info, err := client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath: repo.Path,
		Target:   pr.Target,
		Head:     pr.Head,
	}, repo.NodeId)
	// 不可合并
	if !info.CanMerge {
		err = util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
		return
	}
	err = mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		return s.mergeWithTx(ctx, pr, info, repo, reqDTO.Operator, "", gitCfg)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 触发webhook
	triggerWebhook(repo, reqDTO.Operator, pr, "merge")
	return nil
}

func (s *outerImpl) mergeWithTx(ctx context.Context, pr pullrequestmd.PullRequest, info reqvo.DiffRefsResp, repo repomd.RepoInfo, operator apisession.UserInfo, message string, cfg cfgsrv.GitCfg) error {
	b, err := pullrequestmd.UpdatePrStatusAndCommitId(
		ctx,
		pr.Id,
		pullrequestmd.PrOpenStatus,
		pullrequestmd.PrMergedStatus,
		info.TargetCommit.CommitId,
		info.HeadCommit.CommitId,
	)
	if err != nil {
		return err
	}
	if b {
		if message == "" {
			message = fmt.Sprintf(i18n.GetByKey(i18n.PullRequestMergeMessage), pr.Id, pr.CreateBy, operator.Account)
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
				RepoId:        repo.RepoId,
				PrId:          pr.Id,
				PusherAccount: operator.Account,
				PusherEmail:   operator.Email,
				Message:       message,
				AppUrl:        cfg.AppUrl,
			},
		}, repo.NodeId)
		if err != nil {
			if bizerr.IsBizErr(err) {
				return err
			}
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
	}
	return nil
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
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	// 检查是否重复提交
	_, b, err := pullrequestmd.GetReview(ctx, reqDTO.Id, reqDTO.Operator.Account)
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
	pr, b, err := pullrequestmd.GetById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	repo, b, err := repomd.GetByRepoId(ctx, pr.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, reqDTO.Operator.Account)
	if !b {
		err = util.UnauthorizedError()
		return
	}
	if !p.PermDetail.GetRepoPerm(pr.RepoId).CanAccessRepo {
		err = util.UnauthorizedError()
		return
	}
	// 检查是否是保护分支
	cfg, isProtectedBranch, err := branchmd.IsProtectedBranch(ctx, repo.RepoId, pr.Head)
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
		PrId:      reqDTO.Id,
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

// checkPerm 校验权限
func checkPerm(ctx context.Context, prId int64, operator apisession.UserInfo) (pullrequestmd.PullRequest, repomd.RepoInfo, error) {
	pr, b, err := pullrequestmd.GetById(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return pullrequestmd.PullRequest{}, repomd.RepoInfo{}, util.InternalError(err)
	}
	if !b {
		return pullrequestmd.PullRequest{}, repomd.RepoInfo{}, util.InvalidArgsError()
	}
	repo, err := checkPermByRepoId(ctx, pr.RepoId, operator)
	return pr, repo, err
}

// checkPermByRepoId 校验权限
func checkPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.RepoInfo, error) {
	repo, b := reposrv.Inner.GetByRepoId(ctx, repoId)
	if !b {
		return repomd.RepoInfo{}, util.InvalidArgsError()
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return repo, util.UnauthorizedError()
	}
	if !p.PermDetail.GetRepoPerm(repoId).CanHandlePullRequest {
		return repo, util.UnauthorizedError()
	}
	return repo, nil
}

func triggerWebhook(repo repomd.RepoInfo, operator apisession.UserInfo, pr pullrequestmd.PullRequest, actionType string) {
	go func() {
		ctx, closer := mysqlstore.Context(context.Background())
		defer closer.Close()
		// 触发webhook
		hookList, err := webhookmd.ListWebhook(ctx, repo.RepoId, webhookmd.PullRequestHook)
		if err == nil {
			req := webhook.PullRequestActionHook{
				PrId:       pr.Id,
				RepoId:     repo.RepoId,
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
