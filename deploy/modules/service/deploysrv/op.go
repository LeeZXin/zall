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

func checkAppDevelopPermByPipelineId(ctx context.Context, operator apisession.UserInfo, pipelineId int64) (deploymd.Pipeline, error) {
	pipeline, b, err := deploymd.GetPipelineById(ctx, pipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Pipeline{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Pipeline{}, util.InvalidArgsError()
	}
	return pipeline, checkAppDevelopPermByAppId(ctx, operator, pipeline.AppId)
}

func checkAppDevelopPermByStageId(ctx context.Context, operator apisession.UserInfo, stageId int64) (deploymd.Stage, deploymd.Plan, deploymd.Pipeline, error) {
	stage, b, err := deploymd.GetStageByStageId(ctx, stageId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Pipeline{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Pipeline{}, util.InvalidArgsError()
	}
	plan, pipeline, err := checkAppDevelopPermByPlanId(ctx, operator, stage.PlanId)
	if err != nil {
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Pipeline{}, err
	}
	return stage, plan, pipeline, nil
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

func checkDeployPlanPermByPipelineId(ctx context.Context, operator apisession.UserInfo, pipelineId int64) (deploymd.Pipeline, error) {
	pipeline, b, err := deploymd.GetPipelineById(ctx, pipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Pipeline{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Pipeline{}, util.InvalidArgsError()
	}
	return pipeline, checkDeployPlanPermByAppId(ctx, operator, pipeline.AppId)
}

func checkAppDevelopPermByPlanId(ctx context.Context, operator apisession.UserInfo, planId int64) (deploymd.Plan, deploymd.Pipeline, error) {
	plan, b, err := deploymd.GetPlanById(ctx, planId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Plan{}, deploymd.Pipeline{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Plan{}, deploymd.Pipeline{}, util.InvalidArgsError()
	}
	pipeline, err := checkAppDevelopPermByPipelineId(ctx, operator, plan.PipelineId)
	if err != nil {
		return deploymd.Plan{}, deploymd.Pipeline{}, err
	}
	return plan, pipeline, err
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
	return checkCreateDeployPlanPerm(ctx, operator, &app)
}

func checkCreateDeployPlanPerm(ctx context.Context, operator apisession.UserInfo, app *appmd.App) error {
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
	if !p.PermDetail.DevelopAppList.Contains(app.AppId) ||
		!p.PermDetail.TeamPerm.CanCreateDeployPlan {
		return util.UnauthorizedError()
	}
	return nil
}

func checkPipelinePerm(ctx context.Context, appId string, operator apisession.UserInfo) error {
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
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || (!p.IsAdmin && !p.PermDetail.TeamPerm.CanManagePipeline) {
		return util.UnauthorizedError()
	}
	return nil
}

func checkPipelinePermByPipelineId(ctx context.Context, pipelineId int64, operator apisession.UserInfo) error {
	pipeline, b, err := deploymd.GetPipelineById(ctx, pipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkPipelinePerm(ctx, pipeline.AppId, operator)
}

func checkManageServiceSourcePerm(ctx context.Context, operator apisession.UserInfo, app *appmd.App) error {
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
	if !p.PermDetail.DevelopAppList.Contains(app.AppId) ||
		!p.PermDetail.TeamPerm.CanManageServiceSource {
		return util.UnauthorizedError()
	}
	return nil
}

func checkManageServiceSourcePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkManageServiceSourcePerm(ctx, operator, &app)
}

func checkManageServiceSourcePermBySourceId(ctx context.Context, operator apisession.UserInfo, sourceId int64) (deploymd.ServiceSource, error) {
	ret, b, err := deploymd.GetServiceSourceById(ctx, sourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.ServiceSource{}, util.InternalError(err)
	}
	if !b {
		return deploymd.ServiceSource{}, util.InvalidArgsError()
	}
	return ret, checkManageServiceSourcePermByAppId(ctx, operator, ret.AppId)
}
