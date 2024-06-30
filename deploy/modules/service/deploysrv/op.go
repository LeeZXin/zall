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

func checkAppDevelopPermByStageId(ctx context.Context, operator apisession.UserInfo, stageId int64) (deploymd.Stage, deploymd.Plan, deploymd.Service, error) {
	stage, b, err := deploymd.GetStageByStageId(ctx, stageId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Service{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Service{}, util.InvalidArgsError()
	}
	plan, service, err := checkAppDevelopPermByPlanId(ctx, operator, stage.PlanId)
	if err != nil {
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Service{}, err
	}
	return stage, plan, service, nil
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

func checkAppDevelopPermByPlanId(ctx context.Context, operator apisession.UserInfo, planId int64) (deploymd.Plan, deploymd.Service, error) {
	plan, b, err := deploymd.GetPlanById(ctx, planId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Plan{}, deploymd.Service{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Plan{}, deploymd.Service{}, util.InvalidArgsError()
	}
	service, err := checkAppDevelopPermByServiceId(ctx, operator, plan.ServiceId)
	if err != nil {
		return deploymd.Plan{}, deploymd.Service{}, err
	}
	return plan, service, err
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
