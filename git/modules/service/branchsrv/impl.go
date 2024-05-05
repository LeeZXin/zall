package branchsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
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
	err = checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return
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
	err = checkPerm(ctx, pb.RepoId, reqDTO.Operator)
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
	err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	branchList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret, _ := listutil.Map(branchList, func(t branchmd.ProtectedBranch) (ProtectedBranchDTO, error) {
		d := ProtectedBranchDTO{
			Id:     t.Id,
			RepoId: t.RepoId,
			Branch: t.Branch,
		}
		if t.Cfg == nil {
			d.Cfg = branchmd.ProtectedBranchCfg{}
		} else {
			d.Cfg = *t.Cfg
		}
		return d, nil
	})
	return ret, nil
}

// checkPerm 检查权限 只需检查是否是项目管理员
func checkPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) error {
	// 检查仓库是否存在
	repo, b := reposrv.Inner.GetByRepoId(ctx, repoId)
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	// 如果不是 检查用户组权限
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.GetRepoPerm(repoId).CanHandleProtectedBranch {
		return util.UnauthorizedError()
	}
	return nil
}
