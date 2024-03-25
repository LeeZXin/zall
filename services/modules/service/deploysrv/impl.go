package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) GetDeploy(ctx context.Context, reqDTO GetDeployReqDTO) (deploy.Config, error) {
	if err := reqDTO.IsValid(); err != nil {
		return deploy.Config{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploy.Config{}, util.InternalError(err)
	}
	if !b {
		return deploy.Config{}, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPerm(ctx, app.AppId, app.TeamId, reqDTO.Operator)
	if err != nil {
		return deploy.Config{}, err
	}
	d, b, err := deploymd.GetByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploy.Config{}, util.InternalError(err)
	}
	if !b {
		ret := deploy.Config{}
		err = deploymd.InsertDeploy(ctx, deploymd.InsertDeployReqDTO{
			AppId:  reqDTO.AppId,
			Config: ret,
			Env:    reqDTO.Env,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return ret, util.InternalError(err)
		}
		return ret, nil
	}
	return *d.Config, nil
}

// UpdateDeploy 编辑部署配置
func (*outerImpl) UpdateDeploy(ctx context.Context, reqDTO UpdateDeployReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.UpdateDeploy),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 校验权限
	err = checkPerm(ctx, app.AppId, app.TeamId, reqDTO.Operator)
	if err != nil {
		return
	}
	_, err = deploymd.UpdateDeploy(ctx, deploymd.UpdateDeployReqDTO{
		AppId:  reqDTO.AppId,
		Config: reqDTO.Config,
		Env:    reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func checkPerm(ctx context.Context, appId string, teamId int64, operator apisession.UserInfo) error {
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	contains, _ := listutil.Contains(p.PermDetail.DevelopAppList, func(s string) (bool, error) {
		return s == appId, nil
	})
	if contains {
		return nil
	}
	return util.UnauthorizedError()
}
