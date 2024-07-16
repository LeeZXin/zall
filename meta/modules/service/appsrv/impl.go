package appsrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
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
	if err := checkPermAdmin(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
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
	err = checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除应用
		_, err2 := appmd.DeleteByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除配置来源
		err2 = propertymd.DeleteEtcdNodeByAppId(ctx, reqDTO.AppId)
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
		// 删除服务状态来源
		err2 = deploymd.DeleteServiceSourceByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除下线服务
		err2 = discoverymd.DeleteDownServiceByAppId(ctx, reqDTO.AppId)
		if err2 != nil {
			return err2
		}
		// 删除注册中心来源
		return discoverymd.DeleteEtcdNodeByAppId(ctx, reqDTO.AppId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}

	return
}

// GetApp 获取服务信息
func (*outerImpl) GetApp(ctx context.Context, reqDTO GetAppReqDTO) (AppDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return AppDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return AppDTO{}, err
	}
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return AppDTO{}, util.InternalError(err)
	}
	if !b {
		return AppDTO{}, util.InvalidArgsError()
	}
	return AppDTO{
		AppId: app.AppId,
		Name:  app.Name,
	}, nil
}

// ListApp 应用服务列表
func (*outerImpl) ListApp(ctx context.Context, reqDTO ListAppReqDTO) ([]AppDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	appIdList, isAdmin, err := checkAppList(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, err
	}
	var (
		apps []appmd.App
	)
	if len(appIdList) > 0 {
		apps, err = appmd.GetByAppIdList(ctx, appIdList)
	} else if isAdmin {
		// 管理员可访问所有app
		apps, err = appmd.ListAppByTeamId(ctx, reqDTO.TeamId)
	} else {
		apps = make([]appmd.App, 0)
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

func checkAppList(ctx context.Context, operator apisession.UserInfo, teamId int64) ([]string, bool, error) {
	if operator.IsAdmin {
		return nil, true, nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return nil, false, util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil, true, nil
	}
	return p.PermDetail.DevelopAppList, false, nil
}

func (*outerImpl) UpdateApp(ctx context.Context, reqDTO UpdateAppReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验teamId
	if err = checkAdminPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
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
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AppSrvKeysVO.TransferTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
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
	_, b, err := teammd.GetByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	_, b, err = appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
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

func checkPermAdmin(ctx context.Context, operator apisession.UserInfo, teamId int64) error {
	if operator.IsAdmin {
		_, b, err := teammd.GetByTeamId(ctx, teamId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, teamId, operator.Account)
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}

func checkAdminPermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) error {
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
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}
