package smartsrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

const (
	accessRepo = iota
	updateRepo
)

type outerImpl struct{}

func (s *outerImpl) UploadPack(ctx context.Context, reqDTO UploadPackReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.AccessRepo),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	cfg, b := cfgsrv.Inner.GetGitCfg(ctx)
	if !b {
		return util.InternalError(errors.New("can not get git config"))
	}
	if !reqDTO.FromAccessToken {
		// 获取权限
		err = getPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
		if err != nil {
			return
		}
	}
	err = client.UploadPack(reqvo.UploadPackReq{
		RepoPath: reqDTO.Repo.Path,
		C:        reqDTO.C,
	}, reqDTO.Repo.Id, reqDTO.Operator.Account, reqDTO.Operator.Email, cfg.AppUrl)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (s *outerImpl) ReceivePack(ctx context.Context, reqDTO ReceivePackReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.PushRepo),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	cfg, b := cfgsrv.Inner.GetGitCfg(ctx)
	if !b {
		return util.InternalError(errors.New("can not get git config"))
	}
	// 获取权限
	err = getPerm(ctx, reqDTO.Repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return
	}
	err = client.ReceivePack(reqvo.ReceivePackReq{
		RepoPath: reqDTO.Repo.Path,
		C:        reqDTO.C,
	}, reqDTO.Repo.Id, reqDTO.Operator.Account, reqDTO.Operator.Email, cfg.AppUrl)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (s *outerImpl) InfoRefs(ctx context.Context, reqDTO InfoRefsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.FromAccessToken {
		// 获取权限
		err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
		if err != nil {
			return err
		}
	}
	err := client.InfoRefs(reqvo.InfoRefsReq{
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

func getPerm(ctx context.Context, repo repomd.RepoInfo, operator usermd.UserInfo, permCode int) error {
	if operator.IsAdmin {
		return nil
	}
	// 获取权限
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
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
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanUpdateRepo
	}
	if !pass {
		return util.UnauthorizedError()
	}
	return nil
}
