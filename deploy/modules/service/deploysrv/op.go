package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
)

func checkAppDevelopPermByServiceId(ctx context.Context, operator apisession.UserInfo, serviceId int64) (deploymd.Service, error) {
	service, b, err := deploymd.GetServiceById(ctx, serviceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Service{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Service{}, util.InvalidArgsError()
	}
	return service, checkAppDevelopPermByAppId(ctx, operator, service.AppId)
}

func checkAppDevelopPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkAppDevelopPerm(ctx, operator, &app)
}

func checkAppDevelopPerm(ctx context.Context, operator apisession.UserInfo, app *appmd.App) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
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
	if !p.PermDetail.DevelopAppList.Contains(app.AppId) {
		return util.UnauthorizedError()
	}
	return nil
}

func checkDeployPlanPermByServiceId(ctx context.Context, operator apisession.UserInfo, serviceId int64) (deploymd.Service, error) {
	service, b, err := deploymd.GetServiceById(ctx, serviceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Service{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Service{}, util.InvalidArgsError()
	}
	return service, checkDeployPlanPermByAppId(ctx, operator, service.AppId)
}

func checkDeployPlanPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkDeployPlanPerm(ctx, operator, &app)
}

func checkDeployPlanPerm(ctx context.Context, operator apisession.UserInfo, app *appmd.App) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
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
	if !p.PermDetail.DevelopAppList.Contains(app.AppId) || !p.PermDetail.TeamPerm.CanCreateDeployPlan {
		return util.UnauthorizedError()
	}
	return nil
}
