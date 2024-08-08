package appsrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/property/modules/model/propertymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

// CreateApp 创建应用服务
func (*outerImpl) CreateApp(ctx context.Context, reqDTO CreateAppReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验teamId 只有管理员有权限
	if err := checkAdminPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
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
	return nil
}

func (*outerImpl) DeleteApp(ctx context.Context, reqDTO DeleteAppReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	var (
		app   appmd.App
		roles []teammd.Role
	)
	app, err = checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return
	}
	needUpdateRoles := make([]teammd.Role, 0)
	roles, err = teammd.ListRole(ctx, app.TeamId)
	for _, role := range roles {
		if role.Perm != nil {
			appPermList := role.Perm.AppPermList[:]
			// 去除appId
			role.Perm.AppPermList, _ = listutil.Filter(appPermList, func(appPerm perm.AppPermWithId) (bool, error) {
				return appPerm.AppId != reqDTO.AppId, nil
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
		// 删除下线服务
		err2 = discoverymd.DeleteDownServiceByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除注册中心来源
		err2 = discoverymd.DeleteAppEtcdNodeBindByAppId(ctx, reqDTO.AppId)
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
		err = util.InternalError(err)
		return
	}

	return
}

// GetAppWithPerm 获取服务信息
func (*outerImpl) GetAppWithPerm(ctx context.Context, reqDTO GetAppWithPermReqDTO) (AppWithPermDTO, error) {
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
func (*outerImpl) ListApp(ctx context.Context, reqDTO ListAppReqDTO) ([]AppDTO, error) {
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
		appList, _ := listutil.Filter(permDetail.AppPermList, func(t perm.AppPermWithId) (bool, error) {
			return t.CanDevelop, nil
		})
		appIdList, _ := listutil.Map(appList, func(t perm.AppPermWithId) (string, error) {
			return t.AppId, nil
		})
		if len(appIdList) > 0 {
			apps, err = appmd.GetByAppIdList(ctx, appIdList)
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
	ret, _ := listutil.Map(apps, func(t appmd.App) (AppDTO, error) {
		return AppDTO{
			AppId: t.AppId,
			Name:  t.Name,
		}, nil
	})
	return ret, nil
}

// ListAllAppByAdmin 所有应用服务列表 管理员权限
func (*outerImpl) ListAllAppByAdmin(ctx context.Context, reqDTO ListAppReqDTO) ([]AppDTO, error) {
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
	ret, _ := listutil.Map(apps, func(t appmd.App) (AppDTO, error) {
		return AppDTO{
			AppId: t.AppId,
			Name:  t.Name,
		}, nil
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

func (*outerImpl) UpdateApp(ctx context.Context, reqDTO UpdateAppReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验teamId
	if _, err = checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return
	}
	_, err = appmd.UpdateApp(ctx, appmd.UpdateAppReqDTO{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// TransferTeam 迁移team
func (*outerImpl) TransferTeam(ctx context.Context, reqDTO TransferTeamReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才能迁移team
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := teammd.ExistByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
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
	if app.TeamId == reqDTO.TeamId {
		return
	}
	_, err = appmd.TransferTeam(ctx, reqDTO.AppId, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func checkAdminPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return app, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return app, util.InternalError(err)
	}
	if !b || !p.IsAdmin {
		return app, util.UnauthorizedError()
	}
	return app, nil
}

func checkAdminPermByTeamId(ctx context.Context, operator apisession.UserInfo, teamId int64) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
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
