package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
)

func checkDeployConfigPermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManageDeployConfig {
		return nil
	}
	return util.UnauthorizedError()
}

func checkDeployConfigPermByConfigId(ctx context.Context, configId int64, operator apisession.UserInfo) (deploymd.Config, error) {
	config, b, err := deploymd.GetConfigById(ctx, configId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Config{}, util.InvalidArgsError()
	}
	return config, checkDeployConfigPermByAppId(ctx, config.AppId, operator)
}

func checkAccessConfigPerm(ctx context.Context, operator apisession.UserInfo, configId int64) (deploymd.Config, int64, error) {
	config, b, err := deploymd.GetConfigById(ctx, configId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, 0, util.InternalError(err)
	}
	if !b {
		return deploymd.Config{}, 0, util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, config.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, 0, util.InternalError(err)
	}
	if !b {
		return deploymd.Config{}, 0, util.ThereHasBugErr()
	}
	return config, app.TeamId, checkAppDevelopPerm(ctx, operator, &app)
}

func checkAppDevelopPerm(ctx context.Context, operator apisession.UserInfo, app *appmd.App) error {
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	if !p.PermDetail.DevelopAppList.Contains(app.AppId) {
		return util.UnauthorizedError()
	}
	return nil
}

func checkDeployPlanPerm(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanCreateDeployPlan {
		return nil
	}
	return util.UnauthorizedError()
}

func checkPlanService(ctx context.Context, plan *deploymd.Plan, configId int64, operator apisession.UserInfo) (deploymd.Config, error) {
	config, b, err := deploymd.GetConfigById(ctx, configId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, util.InternalError(err)
	}
	if !b || config.Env != plan.Env {
		return deploymd.Config{}, util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, config.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Config{}, util.ThereHasBugErr()
	}
	if app.TeamId != plan.TeamId {
		return deploymd.Config{}, util.InvalidArgsError()
	}
	return config, checkAppDevelopPerm(ctx, operator, &app)
}
