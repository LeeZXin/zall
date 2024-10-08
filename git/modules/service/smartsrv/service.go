package smartsrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

const (
	accessRepo = iota
	updateRepo
)

func UploadPack(ctx context.Context, reqDTO UploadPackReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	cfg, err := cfgsrv.GetGitCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(errors.New("can not get git config"))
	}
	// 获取权限
	err = checkRepoPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	err = client.UploadPack(reqvo.UploadPackReq{
		RepoPath: reqDTO.Repo.Path,
		C:        reqDTO.C,
	}, reqDTO.Repo.Id, reqDTO.Operator.Account, reqDTO.Operator.Email, cfg.HttpUrl)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func ReceivePack(ctx context.Context, reqDTO ReceivePackReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	cfg, err := cfgsrv.GetGitCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(errors.New("can not get git config"))
	}
	// 获取权限
	err = checkRepoPerm(ctx, reqDTO.Repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return err
	}
	err = client.ReceivePack(reqvo.ReceivePackReq{
		RepoPath: reqDTO.Repo.Path,
		C:        reqDTO.C,
	}, reqDTO.Repo.Id, reqDTO.Operator.Account, reqDTO.Operator.Email, cfg.HttpUrl)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func InfoRefs(ctx context.Context, reqDTO InfoRefsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 获取权限
	err := checkRepoPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	err = client.InfoRefs(reqvo.InfoRefsReq{
		Service:  reqDTO.C.Query("service"),
		RepoPath: reqDTO.Repo.Path,
		C:        reqDTO.C,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkRepoPerm(ctx context.Context, repo repomd.Repo, operator usermd.UserInfo, permCode int) error {
	if operator.IsAdmin {
		return nil
	}
	// 获取权限
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case accessRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo
	case updateRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanPushRepo
	}
	if !pass {
		return util.UnauthorizedError()
	}
	return nil
}
