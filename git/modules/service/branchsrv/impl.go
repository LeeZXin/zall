package branchsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) InsertProtectedBranch(ctx context.Context, reqDTO InsertProtectedBranchReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.BranchSrvKeysVO.InsertProtectedBranch),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return
	}
	for _, account := range reqDTO.Cfg.ReviewerList {
		_, b := usersrv.Inner.GetByAccount(ctx, account)
		// 评审账号不合法
		if !b {
			err = util.NewBizErr(apicode.InvalidArgsCode, i18n.UserAccountNotFoundWarnFormat, account)
			return
		}
		// 检查评审者是否有访问代码的权限
		detail, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, account)
		if !b || !detail.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo {
			err = util.NewBizErr(apicode.InvalidArgsCode, i18n.UserAccountUnauthorizedReviewCodeWarnFormat, account)
			return
		}
	}
	if err = branchmd.InsertProtectedBranch(ctx, branchmd.InsertProtectedBranchReqDTO{
		RepoId: reqDTO.RepoId,
		Branch: reqDTO.Branch,
		Cfg:    reqDTO.Cfg,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteProtectedBranch(ctx context.Context, reqDTO DeleteProtectedBranchReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.BranchSrvKeysVO.DeleteProtectedBranch),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pb, b, err := branchmd.GetById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	_, err = checkPerm(ctx, pb.RepoId, reqDTO.Operator)
	if err != nil {
		return
	}
	_, err = branchmd.DeleteById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) ListProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) ([]ProtectedBranchDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	branchList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret, _ := listutil.Map(branchList, func(t branchmd.ProtectedBranchDTO) (ProtectedBranchDTO, error) {
		return ProtectedBranchDTO{
			Id:     t.Id,
			RepoId: t.RepoId,
			Branch: t.Branch,
			Cfg:    t.Cfg,
		}, nil
	})
	return ret, nil
}

// checkPerm 检查权限 只需检查是否是项目管理员
func checkPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.RepoInfo, error) {
	// 检查仓库是否存在
	repo, b := reposrv.Inner.GetByRepoId(ctx, repoId)
	if !b {
		return repomd.RepoInfo{}, util.InvalidArgsError()
	}
	// 如果不是 检查用户组权限
	detail, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return repo, util.UnauthorizedError()
	}
	if !detail.PermDetail.GetRepoPerm(repoId).CanHandleProtectedBranch {
		return repo, util.UnauthorizedError()
	}
	return repo, nil
}
