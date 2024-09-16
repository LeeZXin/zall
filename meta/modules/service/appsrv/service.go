package appsrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppTopic, func(data any) {
			req, ok := data.(event.AppEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.AppCreateAction:
						return events.App.Create
					case event.AppUpdateAction:
						return events.App.Update
					case event.AppDeleteAction:
						return events.App.Delete
					case event.AppTransferAction:
						return events.App.Transfer
					default:
						return false
					}
				})
			}
		})
	})
}

// CreateApp 创建应用服务
func CreateApp(ctx context.Context, reqDTO CreateAppReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验teamId 只有管理员有权限
	team, err := checkAdminPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return err
	}
	b, err := appmd.ExistByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = appmd.InsertApp(ctx, appmd.InsertAppReqDTO{
		AppId:  reqDTO.AppId,
		TeamId: reqDTO.TeamId,
		Name:   reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(
		reqDTO.Operator,
		team,
		appmd.App{
			AppId: reqDTO.AppId,
			Name:  reqDTO.Name,
		},
		event.AppCreateAction,
	)
	return nil
}

func DeleteApp(ctx context.Context, reqDTO DeleteAppReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	var (
		app   appmd.App
		roles []teammd.Role
	)
	app, team, err := checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	needUpdateRoles := make([]teammd.Role, 0)
	roles, err = teammd.ListRole(ctx, app.TeamId)
	for _, role := range roles {
		if role.Perm != nil {
			appPermList := role.Perm.AppPermList[:]
			// 去除appId
			role.Perm.AppPermList = listutil.FilterNe(appPermList, func(appPerm perm.AppPermWithId) bool {
				return appPerm.AppId != reqDTO.AppId
			})
			if len(role.Perm.AppPermList) != len(appPermList) {
				needUpdateRoles = append(needUpdateRoles, role)
			}
		}
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除应用
		_, err2 := appmd.DeleteByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除配置来源绑定
		err2 = propertymd.DeleteAppEtcdNodeBindByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除配置文件
		err2 = propertymd.DeleteFileByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除部署记录
		err2 = propertymd.DeleteDeployByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除部署流水线
		err2 = deploymd.DeletePipelineByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除部署流水线变量
		err2 = deploymd.DeletePipelineVarsByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除发布计划
		err2 = deploymd.DeletePlanByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除发布计划步骤
		err2 = deploymd.DeleteStageByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除服务状态来源绑定
		err2 = deploymd.DeleteAppServiceSourceBindByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除注册中心来源
		err2 = discoverymd.DeleteAppEtcdNodeBindByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除prometheus抓取任务
		err2 = prommd.DeleteScrapeByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 团队权限检查
		if len(needUpdateRoles) > 0 {
			for _, role := range needUpdateRoles {
				_, err2 = teammd.UpdateRoleById(ctx, teammd.UpdateRoleReqDTO{
					RoleId: role.Id,
					Name:   role.Name,
					Perm:   *role.Perm,
				})
				if err2 != nil {
					return err2
				}
			}
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppDeleteAction,
	)
	return nil
}

// GetAppWithPerm 获取服务信息
func GetAppWithPerm(ctx context.Context, reqDTO GetAppWithPermReqDTO) (AppWithPermDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return AppWithPermDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	appPerm, app, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return AppWithPermDTO{}, err
	}
	return AppWithPermDTO{
		AppDTO: AppDTO{
			AppId: app.AppId,
			Name:  app.Name,
		},
		Perm: appPerm,
	}, nil
}

// ListApp 应用服务列表
func ListApp(ctx context.Context, reqDTO ListAppReqDTO) ([]AppDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	permDetail, isAdmin, err := checkAppList(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, err
	}
	var (
		apps []appmd.App
	)
	if isAdmin {
		// 管理员可访问所有app
		apps, err = appmd.ListAppByTeamId(ctx, reqDTO.TeamId)
	} else if len(permDetail.AppPermList) > 0 {
		appList := listutil.FilterNe(permDetail.AppPermList, func(t perm.AppPermWithId) bool {
			return t.CanDevelop
		})
		appIdList := listutil.MapNe(appList, func(t perm.AppPermWithId) string {
			return t.AppId
		})
		if len(appIdList) > 0 {
			apps, err = appmd.GetByAppIdList(ctx, appIdList, []string{"app_id", "name"})
		} else {
			apps = []appmd.App{}
		}
	} else if permDetail.DefaultAppPerm.CanDevelop {
		// 默认访问权限可访问所有app
		apps, err = appmd.ListAppByTeamId(ctx, reqDTO.TeamId)
	} else {
		apps = []appmd.App{}
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := listutil.MapNe(apps, func(t appmd.App) AppDTO {
		return AppDTO{
			AppId: t.AppId,
			Name:  t.Name,
		}
	})
	return ret, nil
}

// ListAllAppByAdmin 所有应用服务列表 管理员权限
func ListAllAppByAdmin(ctx context.Context, reqDTO ListAllAppByAdminReqDTO) ([]AppDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, isAdmin, err := checkAppList(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, util.UnauthorizedError()
	}
	// 管理员可访问所有app
	apps, err := appmd.ListAppByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := listutil.MapNe(apps, func(t appmd.App) AppDTO {
		return AppDTO{
			AppId: t.AppId,
			Name:  t.Name,
		}
	})
	return ret, nil
}

// ListAllAppBySa 所有应用服务列表 超级管理员权限
func ListAllAppBySa(ctx context.Context, reqDTO ListAllAppBySaReqDTO) ([]AppDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	apps, err := appmd.ListAllApp(ctx, []string{"app_id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := listutil.MapNe(apps, func(t appmd.App) AppDTO {
		return AppDTO{
			AppId: t.AppId,
			Name:  t.Name,
		}
	})
	return ret, nil
}

func checkAppList(ctx context.Context, operator apisession.UserInfo, teamId int64) (perm.Detail, bool, error) {
	if operator.IsAdmin {
		return perm.Detail{}, true, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		return perm.Detail{}, false, util.InternalError(err)
	}
	if !b {
		return perm.Detail{}, false, util.UnauthorizedError()
	}
	if p.IsAdmin {
		return perm.Detail{}, true, nil
	}
	return p.PermDetail, false, nil
}

func UpdateApp(ctx context.Context, reqDTO UpdateAppReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验teamId
	app, team, err := checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	_, err = appmd.UpdateApp(ctx, appmd.UpdateAppReqDTO{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppUpdateAction,
	)
	return nil
}

// TransferTeam 迁移team
func TransferTeam(ctx context.Context, reqDTO TransferTeamReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 只有系统管理员才能迁移team
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	transferTeam, b, err := teammd.GetByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if app.TeamId == reqDTO.TeamId {
		return nil
	}
	appTeam, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	_, err = appmd.TransferTeam(ctx, reqDTO.AppId, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTransferEvent(
		reqDTO.Operator,
		appTeam,
		app,
		transferTeam,
	)
	return nil
}

func notifyEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, action event.AppEventAction) {
	initPsub()
	psub.Publish(event.AppTopic, event.AppEvent{
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action: action,
	})
}

func notifyTransferEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, transferTeam teammd.Team) {
	initPsub()
	psub.Publish(event.AppTopic, event.AppEvent{
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, event.AppTransferAction.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, event.AppTransferAction.GetI18nValue()),
		},
		Action: event.AppTransferAction,
		TransferTeam: &event.BaseTeam{
			TeamId:   transferTeam.Id,
			TeamName: transferTeam.Name,
		},
	})
}

func checkAdminPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkAdminPermByTeamId(ctx, operator, app.TeamId)
	if err != nil {
		return appmd.App{}, teammd.Team{}, err
	}
	return app, team, nil
}

func checkAdminPermByTeamId(ctx context.Context, operator apisession.UserInfo, teamId int64) (teammd.Team, error) {
	team, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Team{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return team, util.InternalError(err)
	}
	if !b || !p.IsAdmin {
		return team, util.UnauthorizedError()
	}
	return team, nil
}

func checkAppDevelopPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (perm.AppPerm, appmd.App, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return perm.AppPerm{}, appmd.App{}, util.InternalError(err)
	}
	if !b {
		return perm.AppPerm{}, appmd.App{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return perm.DefaultAppPerm, app, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return perm.AppPerm{}, appmd.App{}, util.InternalError(err)
	}
	if !b || (!p.IsAdmin && !p.PermDetail.GetAppPerm(appId).CanDevelop) {
		return perm.AppPerm{}, app, util.UnauthorizedError()
	}
	return p.PermDetail.GetAppPerm(appId), app, nil
}
