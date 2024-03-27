package deploysrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"strings"
)

var (
	httpClient = httputil.NewRetryableHttpClient()
)

type outerImpl struct{}

func (*outerImpl) ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]ConfigDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err = checkDeployConfigPerm(ctx, app.TeamId, reqDTO.Operator)
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
	app, b, err := appmd.GetByAppId(ctx, cfg.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		logger.Logger.WithContext(ctx).Errorf("configId: %v 's appId: %v does not exists", reqDTO.ConfigId, cfg.AppId)
		err = util.ThereHasBugErr()
		return
	}
	// 校验权限
	err = checkDeployConfigPerm(ctx, app.TeamId, reqDTO.Operator)
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
	// 校验权限
	err = checkDeployConfigPerm(ctx, app.TeamId, reqDTO.Operator)
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
		PlanStatus: deploymd.Created,
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

func checkDeployConfigPerm(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
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
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanHandleDeployConfig {
		return nil
	}
	return util.UnauthorizedError()
}

func checkDeployPlanPerm(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
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
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanHandleDeployPlan {
		return nil
	}
	return util.UnauthorizedError()
}

type innerImpl struct{}

// DeployServiceWithoutPlan 不通过发布计划部署服务
func (*innerImpl) DeployServiceWithoutPlan(ctx context.Context, reqDTO DeployServiceWithoutPlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
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
	return deployService(&config, reqDTO.ProductVersion, reqDTO.Env, reqDTO.Operator)
}

func deployService(config *deploymd.Config, productVersion, env, operator string) error {
	switch config.ServiceType {
	case deploy.ProcessServiceType:
		var p deploy.ProcessConfig
		err := json.Unmarshal([]byte(config.Content), &p)
		if err != nil {
			logger.Logger.Errorf("configId: %v unmarshal processConfig err: %v", config.Id, err)
			return err
		}
		return deployProcessService(config, &p, productVersion, env, operator)
	case deploy.K8sServiceType:
		return nil
	}
	return fmt.Errorf("configId: %v, unknown service type: %v ", config.Id, config.ServiceType)
}

func deployProcessService(config *deploymd.Config, p *deploy.ProcessConfig, productVersion, env, operator string) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	service, b, err := deploymd.GetServiceByConfigId(ctx, config.Id, env)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	if !b {
		// 插入服务列表
		err = deploymd.InsertService(ctx, deploymd.InsertServiceReqDTO{
			ConfigId:           config.Id,
			CurrProductVersion: productVersion,
			ServiceType:        config.ServiceType,
			ServiceConfig:      config.Content,
			Env:                env,
		})
		if err != nil {
			logger.Logger.Error(err)
			return err
		}
	} else {
		// 更新服务列表
		_, err = deploymd.UpdateService(ctx, deploymd.UpdateServiceReqDTO{
			ConfigId:           config.Id,
			CurrProductVersion: productVersion,
			LastProductVersion: service.CurrProductVersion,
			ServiceConfig:      config.Content,
			Env:                env,
		})
		if err != nil {
			logger.Logger.Error(err)
			return err
		}
	}
	// 执行部署脚本
	go func() {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		script := p.DeployScript
		script = strings.ReplaceAll(script, "{{productVersion}}", productVersion)
		command := action.NewServiceCommand(p.AgentHost, p.AgentToken, config.AppId)
		result, err := command.Execute(strings.NewReader(script), nil)
		var deployOutput string
		if err != nil {
			deployOutput = err.Error()
		} else {
			deployOutput = result
		}
		// 插入日志
		err = deploymd.InsertLog(ctx, deploymd.InsertLogReqDTO{
			ConfigId:       config.Id,
			AppId:          config.AppId,
			ServiceType:    config.ServiceType,
			ServiceConfig:  config.Content,
			ProductVersion: productVersion,
			Env:            env,
			DeployOutput:   deployOutput,
			Operator:       operator,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}()
	return nil
}
