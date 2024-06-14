package servicesrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/servicemd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct {
}

// CreateService 创建服务
func (*outerImpl) CreateService(ctx context.Context, reqDTO CreateServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePerm(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return err
	}
	err := servicemd.InsertService(ctx, servicemd.InsertServiceReqDTO{
		AppId:       reqDTO.AppId,
		ServiceType: reqDTO.service.Type,
		Config:      reqDTO.Config,
		Env:         reqDTO.Env,
		Name:        reqDTO.Name,
		Probed:      0,
		IsEnabled:   false,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateService 编辑服务
func (*outerImpl) UpdateService(ctx context.Context, reqDTO UpdateServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePermByServiceId(ctx, reqDTO.ServiceId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := servicemd.UpdateService(ctx, servicemd.UpdateServiceReqDTO{
		ServiceId:   reqDTO.ServiceId,
		Name:        reqDTO.Name,
		Config:      reqDTO.Config,
		ServiceType: reqDTO.service.Type,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteService 删除服务
func (*outerImpl) DeleteService(ctx context.Context, reqDTO DeleteServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePermByServiceId(ctx, reqDTO.ServiceId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := servicemd.DeleteService(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListService 服务列表
func (*outerImpl) ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]ServiceDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePerm(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	const pageSize = 10
	probes, total, err := servicemd.ListService(ctx, servicemd.ListServiceReqDTO{
		AppId:    reqDTO.AppId,
		Env:      reqDTO.Env,
		PageNum:  reqDTO.PageNum,
		PageSize: pageSize,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(probes, func(t servicemd.Service) (ServiceDTO, error) {
		return ServiceDTO{
			Id:          t.Id,
			AppId:       t.AppId,
			Config:      t.Config,
			ServiceType: t.ServiceType,
			Env:         t.Env,
			Name:        t.Name,
			IsEnabled:   t.IsEnabled,
		}, nil
	})
	return data, total, nil
}

// EnableService 启动服务
func (*outerImpl) EnableService(ctx context.Context, reqDTO EnableServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePermByServiceId(ctx, reqDTO.ServiceId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := servicemd.UpdateServiceIsEnabled(ctx, reqDTO.ServiceId, true)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DisableService 关闭服务
func (*outerImpl) DisableService(ctx context.Context, reqDTO DisableServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePermByServiceId(ctx, reqDTO.ServiceId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := servicemd.UpdateServiceIsEnabled(ctx, reqDTO.ServiceId, false)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkServicePerm(ctx context.Context, appId string, operator apisession.UserInfo) error {
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
	if !b || (!p.IsAdmin && !p.PermDetail.TeamPerm.CanManageService) {
		return util.UnauthorizedError()
	}
	return nil
}

func checkServicePermByServiceId(ctx context.Context, probeId int64, operator apisession.UserInfo) error {
	probe, b, err := servicemd.GetServiceById(ctx, probeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkServicePerm(ctx, probe.AppId, operator)
}
