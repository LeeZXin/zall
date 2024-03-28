package deploysrv

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]ConfigDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkDeployConfigPerm(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	configs, err := deploymd.ListConfigByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(configs, func(t deploymd.Config) (ConfigDTO, error) {
		ret := ConfigDTO{
			Id:          t.Id,
			AppId:       t.AppId,
			Name:        t.Name,
			ServiceType: t.ServiceType,
			Created:     t.Created,
		}
		switch t.ServiceType {
		case deploy.ProcessServiceType:
			ret.ProcessConfig = new(deploy.ProcessConfig)
			json.Unmarshal([]byte(t.Content), ret.ProcessConfig)
		case deploy.K8sServiceType:
			ret.K8sConfig = new(deploy.K8sConfig)
			json.Unmarshal([]byte(t.Content), ret.K8sConfig)
		}
		return ret, nil
	})
}

// UpdateConfig 编辑部署配置
func (*outerImpl) UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.UpdateConfig),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	cfg, b, err := deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
	}
	// 校验权限
	err = checkDeployConfigPerm(ctx, cfg.AppId, reqDTO.Operator)
	if err != nil {
		return
	}
	var cfgJson []byte
	switch cfg.ServiceType {
	case deploy.ProcessServiceType:
		if reqDTO.ProcessConfig == nil || !reqDTO.ProcessConfig.IsValid() {
			err = util.InvalidArgsError()
			return
		}
		cfgJson, _ = json.Marshal(reqDTO.ProcessConfig)
	case deploy.K8sServiceType:
		if reqDTO.K8sConfig == nil || !reqDTO.K8sConfig.IsValid() {
			err = util.InvalidArgsError()
			return
		}
		cfgJson, _ = json.Marshal(reqDTO.K8sConfig)
	default:
		logger.Logger.WithContext(ctx).Errorf("configId: %v 's serviceType: %v is unknown", reqDTO.ConfigId, cfg.ServiceType)
		err = util.ThereHasBugErr()
		return
	}
	_, err = deploymd.UpdateConfig(ctx, deploymd.UpdateConfigReqDTO{
		ConfigId: reqDTO.ConfigId,
		Name:     reqDTO.Name,
		Content:  string(cfgJson),
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// InsertConfig 新增部署配置
func (*outerImpl) InsertConfig(ctx context.Context, reqDTO InsertConfigReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.InsertConfig),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err = checkDeployConfigPerm(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return
	}
	var cfgJson []byte
	switch reqDTO.ServiceType {
	case deploy.ProcessServiceType:
		cfgJson, _ = json.Marshal(reqDTO.ProcessConfig)
	case deploy.K8sServiceType:
		cfgJson, _ = json.Marshal(reqDTO.ProcessConfig)
	default:
		err = util.ThereHasBugErr()
		return
	}
	err = deploymd.InsertConfig(ctx, deploymd.InsertConfigReqDTO{
		AppId:       reqDTO.AppId,
		Name:        reqDTO.Name,
		ServiceType: reqDTO.ServiceType,
		Content:     string(cfgJson),
		Env:         reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// InsertPlan 创建发布计划
func (*outerImpl) InsertPlan(ctx context.Context, reqDTO InsertPlanReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.InsertPlan),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	hasPerm := checkDeployPlanPerm(ctx, reqDTO.TeamId, reqDTO.Operator)
	if hasPerm != nil {
		var (
			teamConfig deploymd.TeamConfig
			b          bool
		)
		teamConfig, b, err = deploymd.GetTeamConfig(ctx, reqDTO.TeamId, reqDTO.Env)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if b && teamConfig.Content.InsertPlanWithApprovalFlow {
			// todo 创建审批流
			return
		}
		err = hasPerm
		return
	}
	err = deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
		Name:       reqDTO.Name,
		PlanStatus: deploymd.CreatedPlanStatus,
		TeamId:     reqDTO.TeamId,
		Creator:    reqDTO.Operator.Account,
		Env:        reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// DeployService 部署服务
func (*outerImpl) DeployService(ctx context.Context, reqDTO DeployServiceReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.ReDeployService),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查参数
	var (
		config deploymd.Config
		b      bool
	)
	config, b, err = deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 检查权限
	if err = checkAppDevelopPerm(ctx, config.AppId, reqDTO.Operator); err != nil {
		return
	}
	// 检查制品
	_, b, err = productmd.GetProduct(ctx, productmd.GetProductReqDTO{
		AppId: config.AppId,
		Name:  reqDTO.ProductVersion,
		Env:   reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	err = deployService(&config, reqDTO.ProductVersion, reqDTO.Env, reqDTO.Operator.Account, 0)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// StopService 下线服务
func (*outerImpl) StopService(ctx context.Context, reqDTO StopServiceReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.StopService),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查参数
	var (
		config  deploymd.Config
		b       bool
		service deploymd.Service
	)
	config, b, err = deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 检查权限
	if err = checkAppDevelopPerm(ctx, config.AppId, reqDTO.Operator); err != nil {
		return
	}
	service, b, err = deploymd.GetServiceByConfigId(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	err = stopService(&config, &service, reqDTO.Env, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// RestartService 重启服务
func (*outerImpl) RestartService(ctx context.Context, reqDTO RestartServiceReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.RestartService),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查参数
	var (
		config  deploymd.Config
		b       bool
		service deploymd.Service
	)
	config, b, err = deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 检查权限
	if err = checkAppDevelopPerm(ctx, config.AppId, reqDTO.Operator); err != nil {
		return
	}
	service, b, err = deploymd.GetServiceByConfigId(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	err = restartService(&config, &service, reqDTO.Env, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListService 服务列表
func (*outerImpl) ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]ServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err := checkAppDevelopPerm(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	// 获取所有配置
	configs, err := deploymd.ListConfigByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	configIdList, _ := listutil.Map(configs, func(t deploymd.Config) (int64, error) {
		return t.Id, nil
	})
	services, err := deploymd.ListServiceByConfigIdList(ctx, configIdList, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(services, func(t deploymd.Service) (ServiceDTO, error) {
		ret := ServiceDTO{
			CurrProductVersion: t.CurrProductVersion,
			LastProductVersion: t.LastProductVersion,
			ServiceType:        t.ServiceType,
			ActiveStatus:       t.ActiveStatus,
			StartTime:          t.StartTime,
			ProbeTime:          t.ProbeTime,
			Created:            t.Created,
		}
		switch t.ServiceType {
		case deploy.ProcessServiceType:
			ret.ProcessConfig = new(deploy.ProcessConfig)
			json.Unmarshal([]byte(t.ServiceConfig), ret.ProcessConfig)
		case deploy.K8sServiceType:
			ret.K8sConfig = new(deploy.K8sConfig)
			json.Unmarshal([]byte(t.ServiceConfig), ret.K8sConfig)
		}
		return ret, nil
	})
}

// ListDeployLog 查看部署日志
func (*outerImpl) ListDeployLog(ctx context.Context, reqDTO ListDeployLogReqDTO) ([]DeployLogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	config, b, err := deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 检查权限
	if err = checkAppDevelopPerm(ctx, config.AppId, reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	logs, err := deploymd.ListDeployLog(ctx, deploymd.ListDeployLogReqDTO{
		ConfigId: reqDTO.ConfigId,
		Cursor:   reqDTO.Cursor,
		Limit:    reqDTO.Limit,
		Env:      reqDTO.Env,
	})
	ret, _ := listutil.Map(logs, func(t deploymd.DeployLog) (DeployLogDTO, error) {
		return DeployLogDTO{
			ServiceType:    t.ServiceType,
			ServiceConfig:  t.ServiceConfig,
			ProductVersion: t.ProductVersion,
			Operator:       t.Operator,
			DeployOutput:   t.DeployOutput,
			Created:        t.Created,
			PlanId:         t.PlanId,
		}, nil
	})
	if len(logs) == reqDTO.Limit {
		return ret, logs[len(logs)-1].Id, nil
	}
	return ret, 0, nil
}

// ListOpLog 查看操作日志
func (*outerImpl) ListOpLog(ctx context.Context, reqDTO ListOpLogReqDTO) ([]OpLogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	config, b, err := deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 检查权限
	if err = checkAppDevelopPerm(ctx, config.AppId, reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	logs, err := deploymd.ListOpLog(ctx, deploymd.ListOpLogReqDTO{
		ConfigId: reqDTO.ConfigId,
		Cursor:   reqDTO.Cursor,
		Limit:    reqDTO.Limit,
		Env:      reqDTO.Env,
	})
	ret, _ := listutil.Map(logs, func(t deploymd.OpLog) (OpLogDTO, error) {
		return OpLogDTO{
			Op:             t.Op,
			Operator:       t.Operator,
			ScriptOutput:   t.ScriptOutput,
			ProductVersion: t.ProductVersion,
			Created:        t.Created,
		}, nil
	})
	if len(logs) == reqDTO.Limit {
		return ret, logs[len(logs)-1].Id, nil
	}
	return ret, 0, nil
}

type innerImpl struct{}

// DeployServiceWithoutPlan 不通过发布计划部署服务
func (*innerImpl) DeployServiceWithoutPlan(ctx context.Context, reqDTO DeployServiceWithoutPlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 如果有configId 优先configId 否则通过appId获取全部的部署配置
	if reqDTO.ConfigId > 0 {
		// 检查配置
		config, b, err := deploymd.GetConfigById(ctx, reqDTO.ConfigId, reqDTO.Env)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
		// 检查制品
		_, b, err = productmd.GetProduct(ctx, productmd.GetProductReqDTO{
			AppId: config.AppId,
			Name:  reqDTO.ProductVersion,
			Env:   reqDTO.Env,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
		err = deployService(&config, reqDTO.ProductVersion, reqDTO.Env, reqDTO.Operator, 0)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
	} else if reqDTO.AppId != "" {
		// 检查制品
		_, b, err := productmd.GetProduct(ctx, productmd.GetProductReqDTO{
			AppId: reqDTO.AppId,
			Name:  reqDTO.ProductVersion,
			Env:   reqDTO.Env,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
		// 获取所有部署配置
		configs, err := deploymd.ListConfigByAppId(ctx, reqDTO.AppId, reqDTO.Env)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		for i := range configs {
			err = deployService(&configs[i], reqDTO.ProductVersion, reqDTO.Env, reqDTO.Operator, 0)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return util.InternalError(err)
			}
		}
	} else {
		return util.InvalidArgsError()
	}
	return nil
}
