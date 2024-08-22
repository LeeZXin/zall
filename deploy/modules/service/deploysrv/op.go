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

func checkAppDevelopPermByPipelineId(ctx context.Context, operator apisession.UserInfo, pipelineId int64) (deploymd.Pipeline, appmd.App, teammd.Team, error) {
	pipeline, b, err := deploymd.GetPipelineById(ctx, pipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, operator, pipeline.AppId)
	return pipeline, app, team, err
}

func checkAppDevelopPermByStageId(ctx context.Context, operator apisession.UserInfo, stageId int64) (deploymd.Stage, deploymd.Plan, deploymd.Pipeline, appmd.App, teammd.Team, error) {
	stage, b, err := deploymd.GetStageByStageId(ctx, stageId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Stage{}, deploymd.Plan{}, deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	plan, pipeline, app, team, err := checkAppDevelopPermByPlanId(ctx, operator, stage.PlanId)
	return stage, plan, pipeline, app, team, err
}

func checkAppDevelopPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	return app, team, checkAppDevelopPerm(ctx, operator, &app)
}

func checkAppDevelopPermByBindId(ctx context.Context, operator apisession.UserInfo, bindId int64) (appmd.App, teammd.Team, deploymd.ServiceSource, error) {
	bind, b, err := deploymd.GetAppServiceSourceBindById(ctx, bindId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, deploymd.ServiceSource{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, deploymd.ServiceSource{}, util.InvalidArgsError()
	}
	source, b, err := deploymd.GetServiceSourceById(ctx, bind.SourceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, deploymd.ServiceSource{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, deploymd.ServiceSource{}, util.ThereHasBugErr()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, operator, bind.AppId)
	return app, team, source, err
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
	if !p.PermDetail.GetAppPerm(app.AppId).CanDevelop {
		return util.UnauthorizedError()
	}
	return nil
}

func checkDeployPlanPermByPipelineId(ctx context.Context, operator apisession.UserInfo, pipelineId int64) (deploymd.Pipeline, appmd.App, teammd.Team, error) {
	pipeline, b, err := deploymd.GetPipelineById(ctx, pipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkDeployPlanPermByAppId(ctx, operator, pipeline.AppId)
	return pipeline, app, team, err
}

func checkAppDevelopPermByPlanId(ctx context.Context, operator apisession.UserInfo, planId int64) (deploymd.Plan, deploymd.Pipeline, appmd.App, teammd.Team, error) {
	plan, b, err := deploymd.GetPlanById(ctx, planId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Plan{}, deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Plan{}, deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	pipeline, app, team, err := checkAppDevelopPermByPipelineId(ctx, operator, plan.PipelineId)
	return plan, pipeline, app, team, err
}

func checkDeployPlanPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	return app, team, checkCreateDeployPlanPerm(ctx, operator, &app)
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
	if !p.PermDetail.GetAppPerm(app.AppId).CanCreateDeployPlan {
		return util.UnauthorizedError()
	}
	return nil
}

func checkManagePipelinePermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return app, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return app, team, util.InternalError(err)
	}
	if !b {
		return app, team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.GetAppPerm(appId).CanManagePipeline {
		return app, team, nil
	}
	return app, team, util.UnauthorizedError()
}

func checkManagePipelinePermByVarsId(ctx context.Context, varsId int64, operator apisession.UserInfo) (deploymd.PipelineVars, appmd.App, teammd.Team, error) {
	vars, b, err := deploymd.GetPipelineVarsById(ctx, varsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.PipelineVars{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return deploymd.PipelineVars{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkManagePipelinePermByAppId(ctx, vars.AppId, operator)
	return vars, app, team, err
}

func checkPipelinePermByPipelineId(ctx context.Context, pipelineId int64, operator apisession.UserInfo) (deploymd.Pipeline, appmd.App, teammd.Team, error) {
	pipeline, b, err := deploymd.GetPipelineById(ctx, pipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return deploymd.Pipeline{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkManagePipelinePermByAppId(ctx, pipeline.AppId, operator)
	return pipeline, app, team, err
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
	if !p.PermDetail.GetAppPerm(app.AppId).CanManageServiceSource {
		return util.UnauthorizedError()
	}
	return nil
}

func checkManageServiceSourcePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	return app, team, checkManageServiceSourcePerm(ctx, operator, &app)
}
