package deploysrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"strconv"
	"strings"
	"time"
)

func checkDeployConfigPerm(ctx context.Context, appId string, operator apisession.UserInfo) error {
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
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanHandleDeployConfig {
		return nil
	}
	return util.UnauthorizedError()
}

func checkAccessConfigPerm(ctx context.Context, operator apisession.UserInfo, configId int64, env string) (deploymd.Config, int64, error) {
	config, b, err := deploymd.GetConfigById(ctx, configId, env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, 0, util.InternalError(err)
	}
	if !b {
		return deploymd.Config{}, 0, util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, config.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return deploymd.Config{}, 0, util.InternalError(err)
	}
	if !b {
		return deploymd.Config{}, 0, util.ThereHasBugErr()
	}
	return config, app.TeamId, checkAppDevelopPerm(ctx, operator, app)
}

func checkAppDevelopPerm(ctx context.Context, operator apisession.UserInfo, apps ...appmd.App) error {
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, apps[0].TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	contains, _ := listutil.Contains(p.PermDetail.DevelopAppList, func(developAppId string) (bool, error) {
		return listutil.Contains(apps, func(app appmd.App) (bool, error) {
			return app.AppId == developAppId, nil
		})
	})
	if contains {
		return nil
	}
	return util.UnauthorizedError()
}

func checkDeployPlanPerm(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
	if operator.IsAdmin {
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

func checkDeployItems(ctx context.Context, teamId int64, items deploymd.DeployItems, operator apisession.UserInfo, env string) error {
	prdVerMap := make(map[int64]string)
	for _, item := range items {
		prdVerMap[item.ConfigId] = item.ProductVersion
	}
	// 检查重复的configId
	configIdList, _ := listutil.Map(items, func(t deploymd.DeployItem) (int64, error) {
		return t.ConfigId, nil
	})
	if hashset.NewHashSet(configIdList...).Size() != len(items) {
		return util.InvalidArgsError()
	}
	configs, err := deploymd.BatchGetConfigById(ctx, configIdList, env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(configs) != len(configIdList) {
		return util.InvalidArgsError()
	}
	appIdList, _ := listutil.Map(configs, func(t deploymd.Config) (string, error) {
		return t.AppId, nil
	})
	appIdList = listutil.Distinct(appIdList...)
	var apps []appmd.App
	apps, err = appmd.GetByAppIdList(ctx, appIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(apps) != len(appIdList) {
		return util.ThereHasBugErr()
	}
	for _, app := range apps {
		if app.TeamId != teamId {
			return util.InvalidArgsError()
		}
	}
	// configId -> productVersion ===> appId -> productVersion
	appIdMap := make(map[string]string)
	for _, config := range configs {
		appIdMap[config.AppId] = prdVerMap[config.Id]
	}
	// 检查制品
	for appId, productVersion := range appIdMap {
		_, b, err := productmd.GetProduct(ctx, productmd.GetProductReqDTO{
			AppId: appId,
			Name:  productVersion,
			Env:   env,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
	}
	return checkAppDevelopPerm(ctx, operator, apps...)
}

func deployService(config *deploymd.Config, productVersion, env, operator string, planId int64) error {
	switch config.ServiceType {
	case deploy.ProcessServiceType:
		var p deploy.ProcessConfig
		err := json.Unmarshal([]byte(config.Content), &p)
		if err != nil {
			return fmt.Errorf("configId: %v unmarshal processConfig err: %v", config.Id, err)
		}
		if !p.IsValid() {
			return fmt.Errorf("configId: %v invalid processConfig", config.Id)
		}
		return deployProcessService(config, &p, productVersion, env, operator, planId)
	case deploy.K8sServiceType:
		return nil
	}
	return fmt.Errorf("configId: %v, unknown service type: %v ", config.Id, config.ServiceType)
}

func deployProcessService(config *deploymd.Config, p *deploy.ProcessConfig, productVersion, env, operator string, planId int64) error {
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
			StartTime:          time.Now().UnixMilli(),
			ActiveStatus:       deploymd.StartingStatus, // 启动中
		})
		if err != nil {
			return err
		}
	} else {
		// 更新服务列表
		b, err = deploymd.UpdateServiceWithOldStatus(ctx, service.ActiveStatus, deploymd.UpdateServiceReqDTO{
			ConfigId:           config.Id,
			CurrProductVersion: productVersion,
			LastProductVersion: service.CurrProductVersion,
			ServiceConfig:      config.Content,
			Env:                env,
			ActiveStatus:       deploymd.StartingStatus, // 启动中
			ProbeTime:          0,
			StartTime:          time.Now().UnixMilli(),
		})
		if err != nil {
			return err
		}
		if !b {
			return nil
		}
	}
	// 执行部署脚本
	go func() {
		script := p.DeployScript
		script = strings.ReplaceAll(script, "{{configId}}", strconv.FormatInt(config.Id, 10))
		script = strings.ReplaceAll(script, "{{appId}}", config.AppId)
		script = strings.ReplaceAll(script, "{{productVersion}}", productVersion)
		command := zssh.NewServiceCommand(p.AgentHost, p.AgentToken, config.AppId)
		result, err := command.Execute(strings.NewReader(script), nil)
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		var deployOutput string
		if err != nil {
			// 启动失败
			_, err = deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.AbnormalStatus, deploymd.StartingStatus)
			if err != nil {
				logger.Logger.Error(err)
			}
			deployOutput = err.Error()
		} else {
			// 启动成功
			_, err = deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.StartedStatus, deploymd.StartingStatus)
			if err != nil {
				logger.Logger.Error(err)
			}
			deployOutput = result
		}
		// 插入日志
		err = deploymd.InsertDeployLog(ctx, deploymd.InsertDeployLogReqDTO{
			ConfigId:       config.Id,
			AppId:          config.AppId,
			ServiceType:    config.ServiceType,
			ServiceConfig:  config.Content,
			ProductVersion: productVersion,
			Env:            env,
			DeployOutput:   deployOutput,
			Operator:       operator,
			PlanId:         planId,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}()
	return nil
}

func stopService(config *deploymd.Config, service *deploymd.Service, env, operator string) error {
	switch config.ServiceType {
	case deploy.ProcessServiceType:
		var p deploy.ProcessConfig
		err := json.Unmarshal([]byte(config.Content), &p)
		if err != nil {
			return fmt.Errorf("configId: %v unmarshal processConfig err: %v", config.Id, err)
		}
		if !p.IsValid() {
			return fmt.Errorf("configId: %v invalid processConfig", config.Id)
		}
		return stopProcessService(config, &p, service, env, operator)
	case deploy.K8sServiceType:
		return nil
	}
	return fmt.Errorf("configId: %v, unknown service type: %v ", config.Id, config.ServiceType)

}

func stopProcessService(config *deploymd.Config, p *deploy.ProcessConfig, service *deploymd.Service, env, operator string) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	// 停止中
	b, err := deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.StoppingStatus, service.ActiveStatus)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	if !b {
		return nil
	}
	// 执行停止脚本
	go func() {
		script := p.StopScript
		script = strings.ReplaceAll(script, "{{configId}}", strconv.FormatInt(config.Id, 10))
		script = strings.ReplaceAll(script, "{{appId}}", config.AppId)
		script = strings.ReplaceAll(script, "{{productVersion}}", service.CurrProductVersion)
		command := zssh.NewServiceCommand(p.AgentHost, p.AgentToken, config.AppId)
		result, err := command.Execute(strings.NewReader(script), nil)
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		if err != nil {
			// 停止失败
			_, err = deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.AbnormalStatus, deploymd.StoppingStatus)
			if err != nil {
				logger.Logger.Error(err)
			}
		} else {
			// 停止成功
			_, err = deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.StoppedStatus, deploymd.StoppingStatus)
			if err != nil {
				logger.Logger.Error(err)
			}
		}
		err = deploymd.InsertOpLog(ctx, deploymd.InsertOpLogReqDTO{
			ConfigId:       config.Id,
			Operator:       operator,
			ScriptOutput:   result,
			Env:            env,
			ProductVersion: service.CurrProductVersion,
			Op:             deploymd.StopServiceOp,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}()
	return nil
}

func restartService(config *deploymd.Config, service *deploymd.Service, env, operator string) error {
	switch config.ServiceType {
	case deploy.ProcessServiceType:
		var p deploy.ProcessConfig
		err := json.Unmarshal([]byte(config.Content), &p)
		if err != nil {
			return fmt.Errorf("configId: %v unmarshal processConfig err: %v", config.Id, err)
		}
		if !p.IsValid() {
			return fmt.Errorf("configId: %v invalid processConfig", config.Id)
		}
		return restartProcessService(config, &p, service, env, operator)
	case deploy.K8sServiceType:
		return nil
	}
	return fmt.Errorf("configId: %v, unknown service type: %v ", config.Id, config.ServiceType)

}

func restartProcessService(config *deploymd.Config, p *deploy.ProcessConfig, service *deploymd.Service, env, operator string) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	// 重启中
	b, err := deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.StartingStatus, service.ActiveStatus)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	if !b {
		return nil
	}
	// 执行重启脚本
	go func() {
		script := p.RestartScript
		script = strings.ReplaceAll(script, "{{configId}}", strconv.FormatInt(config.Id, 10))
		script = strings.ReplaceAll(script, "{{appId}}", config.AppId)
		script = strings.ReplaceAll(script, "{{productVersion}}", service.CurrProductVersion)
		command := zssh.NewServiceCommand(p.AgentHost, p.AgentToken, config.AppId)
		result, err := command.Execute(strings.NewReader(script), nil)
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		if err != nil {
			// 重启失败
			_, err = deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.AbnormalStatus, deploymd.StartingStatus)
			if err != nil {
				logger.Logger.Error(err)
			}
		} else {
			// 重启成功
			_, err = deploymd.UpdateServiceActiveStatusWithOldStatus(ctx, config.Id, env, deploymd.StartedStatus, deploymd.StartingStatus)
			if err != nil {
				logger.Logger.Error(err)
			}
		}
		err = deploymd.InsertOpLog(ctx, deploymd.InsertOpLogReqDTO{
			ConfigId:       config.Id,
			Operator:       operator,
			ScriptOutput:   result,
			Env:            env,
			ProductVersion: service.CurrProductVersion,
			Op:             deploymd.RestartServiceOp,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}()
	return nil
}
