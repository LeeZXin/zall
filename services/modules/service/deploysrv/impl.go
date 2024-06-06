package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/fileserv/modules/model/productmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct{}

func (*outerImpl) ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]ConfigDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkDeployConfigPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	configs, err := deploymd.ListConfigByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(configs, func(t deploymd.Config) (ConfigDTO, error) {
		return ConfigDTO{
			Id:        t.Id,
			AppId:     t.AppId,
			Name:      t.Name,
			Content:   t.Content,
			Env:       t.Env,
			IsEnabled: t.IsEnabled,
		}, nil
	})
}

// UpdateConfig 编辑部署配置
func (*outerImpl) UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkDeployConfigPermByConfigId(ctx, reqDTO.ConfigId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = deploymd.UpdateConfig(ctx, deploymd.UpdateConfigReqDTO{
		ConfigId: reqDTO.ConfigId,
		Name:     reqDTO.Name,
		Content:  reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CreateConfig 新增部署配置
func (*outerImpl) CreateConfig(ctx context.Context, reqDTO CreateConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkDeployConfigPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = deploymd.InsertConfig(ctx, deploymd.InsertConfigReqDTO{
		AppId:     reqDTO.AppId,
		Name:      reqDTO.Name,
		Content:   reqDTO.Content,
		Env:       reqDTO.Env,
		IsEnabled: false,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// EnableConfig 启动配置
func (*outerImpl) EnableConfig(ctx context.Context, reqDTO EnableConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	config, err := checkDeployConfigPermByConfigId(ctx, reqDTO.ConfigId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		configs, err2 := deploymd.ListConfigForUpdate(ctx, config.AppId, config.Env)
		if err2 != nil {
			return err2
		}
		shouldDisabled := make([]int64, 0)
		for _, config := range configs {
			if config.Id != reqDTO.ConfigId {
				shouldDisabled = append(shouldDisabled, config.Id)
			}
		}
		err2 = deploymd.BatchSetConfigIsEnabled(ctx, shouldDisabled, false)
		if err2 != nil {
			return err2
		}
		return deploymd.BatchSetConfigIsEnabled(ctx, []int64{reqDTO.ConfigId}, true)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DisableConfig 关闭配置
func (*outerImpl) DisableConfig(ctx context.Context, reqDTO DisableConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkDeployConfigPermByConfigId(ctx, reqDTO.ConfigId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = deploymd.BatchSetConfigIsEnabled(ctx, []int64{reqDTO.ConfigId}, false)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
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
	// 检查deployItems
	if reqDTO.PlanType == deploymd.AddServiceBeforePlanCreatingType {
		if err = checkDeployItems(ctx, reqDTO.TeamId, reqDTO.DeployItems, reqDTO.Operator, reqDTO.Env); err != nil {
			return
		}
	}
	// 检查权限
	if err = checkDeployPlanPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		// 检查用户是否归属在项目组下
		var b bool
		_, b, err = teammd.GetTeamUser(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if !b {
			err = util.UnauthorizedError()
			return
		}
		if reqDTO.PlanType == deploymd.AddServiceBeforePlanCreatingType {

			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				err = util.InternalError(err)
				return
			}
			return
		}
	}
	switch reqDTO.PlanType {
	case deploymd.AddServiceAfterPlanCreatingType:
		_, err = deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
			Name:       reqDTO.Name,
			PlanStatus: deploymd.RunningPlanStatus,
			PlanType:   reqDTO.PlanType,
			TeamId:     reqDTO.TeamId,
			Creator:    reqDTO.Operator.Account,
			Env:        reqDTO.Env,
			Expired:    time.Now().Add(time.Duration(reqDTO.ExpireHours) * time.Hour),
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	case deploymd.AddServiceBeforePlanCreatingType:
		err = xormstore.WithTx(ctx, func(ctx context.Context) error {
			plan, err := deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
				Name:       reqDTO.Name,
				PlanStatus: deploymd.RunningPlanStatus,
				PlanType:   reqDTO.PlanType,
				TeamId:     reqDTO.TeamId,
				Creator:    reqDTO.Operator.Account,
				Env:        reqDTO.Env,
				Expired:    time.Now().Add(time.Duration(reqDTO.ExpireHours) * time.Hour),
			})
			if err != nil {
				return err
			}
			reqs, _ := listutil.Map(reqDTO.DeployItems, func(t deploymd.DeployItem) (deploymd.InsertPlanItemReqDTO, error) {
				return deploymd.InsertPlanItemReqDTO{
					PlanId:         plan.Id,
					ConfigId:       t.ConfigId,
					ProductVersion: t.ProductVersion,
					ItemStatus:     deploymd.WaitItemStatus,
				}, nil
			})
			return deploymd.BatchInsertPlanItem(ctx, reqDTO.Env, reqs...)
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	}
	return
}

// ClosePlan 关闭发布计划
func (*outerImpl) ClosePlan(ctx context.Context, reqDTO ClosePlanReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.ClosePlan),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		plan deploymd.Plan
		b    bool
	)
	plan, b, err = deploymd.GetPlanById(ctx, reqDTO.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || plan.IsExpired() || plan.PlanStatus != deploymd.RunningPlanStatus {
		err = util.InvalidArgsError()
		return
	}
	err = checkDeployPlanPerm(ctx, plan.TeamId, reqDTO.Operator)
	// 发布计划创建人和管理员有权限关闭
	if err != nil && plan.Creator != reqDTO.Operator.Account {
		return
	}
	_, err = deploymd.UpdatePlanStatusById(ctx, reqDTO.PlanId, reqDTO.Env, deploymd.ClosedPlanStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// InsertPlanItem 添加发布计划部署服务
func (*outerImpl) InsertPlanItem(ctx context.Context, reqDTO InsertPlanItemReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.InsertPlanItem),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		plan deploymd.Plan
		b    bool
		reqs []deploymd.InsertPlanItemReqDTO
	)
	plan, b, err = deploymd.GetPlanById(ctx, reqDTO.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || plan.PlanType != deploymd.AddServiceAfterPlanCreatingType || plan.IsClosed() {
		err = util.InvalidArgsError()
		return
	}
	if err = checkDeployItems(ctx, plan.TeamId, reqDTO.DeployItems, reqDTO.Operator, reqDTO.Env); err != nil {
		return
	}
	var items []deploymd.PlanItem
	items, err = deploymd.ListPlanItemByPlanId(ctx, reqDTO.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	existsConfigIdList, _ := listutil.Map(items, func(t deploymd.PlanItem) (int64, error) {
		return t.ConfigId, nil
	})
	addConfigIdList, _ := listutil.Map(reqDTO.DeployItems, func(t deploymd.DeployItem) (int64, error) {
		return t.ConfigId, nil
	})
	if hashset.NewHashSet(existsConfigIdList...).
		Intersect(hashset.NewHashSet(addConfigIdList...)).
		Size() > 0 {
		err = util.AlreadyExistsError()
		return
	}
	reqs, err = listutil.Map(reqDTO.DeployItems, func(t deploymd.DeployItem) (deploymd.InsertPlanItemReqDTO, error) {
		service, b, err := deploymd.GetServiceByConfigId(ctx, t.ConfigId, reqDTO.Env)
		if err != nil {
			return deploymd.InsertPlanItemReqDTO{}, err
		}
		var lastProductVersion string
		if b {
			lastProductVersion = service.CurrProductVersion
		}
		return deploymd.InsertPlanItemReqDTO{
			PlanId:             reqDTO.PlanId,
			ConfigId:           t.ConfigId,
			ProductVersion:     t.ProductVersion,
			ItemStatus:         deploymd.WaitItemStatus,
			LastProductVersion: lastProductVersion,
		}, nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	err = deploymd.BatchInsertPlanItem(ctx, reqDTO.Env, reqs...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ClosePlanItem 关闭发布计划单项服务
func (*outerImpl) ClosePlanItem(ctx context.Context, reqDTO ClosePlanItemReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.ClosePlanItem),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		planItem deploymd.PlanItem
		b        bool
	)
	planItem, b, err = deploymd.GetPlanItemById(ctx, reqDTO.ItemId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || planItem.ItemStatus == deploymd.ClosedItemStatus {
		err = util.InvalidArgsError()
		return
	}
	// 校验权限
	_, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, planItem.ConfigId)
	if err != nil {
		return
	}
	_, err = deploymd.UpdateItemStatusWithOldStatus(ctx, reqDTO.ItemId, reqDTO.Env, deploymd.ClosedItemStatus, planItem.ItemStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListPlanItem 展示发布计划的服务
func (*outerImpl) ListPlanItem(ctx context.Context, reqDTO ListPlanItemReqDTO) ([]PlanItemDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	_, b, err = teammd.GetTeamUser(ctx, plan.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.UnauthorizedError()
	}
	items, err := deploymd.ListPlanItemByPlanId(ctx, reqDTO.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	configIdList, _ := listutil.Map(items, func(t deploymd.PlanItem) (int64, error) {
		return t.ConfigId, nil
	})
	configs, err := deploymd.BatchGetSimpleConfigById(ctx, configIdList, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	configMap, _ := listutil.CollectToMap(configs, func(t deploymd.Config) (int64, error) {
		return t.Id, nil
	}, func(t deploymd.Config) (deploymd.Config, error) {
		return t, nil
	})
	return listutil.Map(items, func(t deploymd.PlanItem) (PlanItemDTO, error) {
		return PlanItemDTO{
			Id:                 t.Id,
			AppId:              configMap[t.ConfigId].AppId,
			ConfigId:           t.ConfigId,
			ConfigName:         configMap[t.ConfigId].Name,
			ProductVersion:     t.ProductVersion,
			LastProductVersion: t.LastProductVersion,
			ItemStatus:         t.ItemStatus,
			Created:            t.Created,
		}, nil
	})
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
	// 检查权限
	config, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, reqDTO.ConfigId)
	if err != nil {
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
	config, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, reqDTO.ConfigId)
	if err != nil {
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
	config, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, reqDTO.ConfigId)
	if err != nil {
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
	app, b, err := appmd.GetByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 检查权限
	if err = checkAppDevelopPerm(ctx, reqDTO.Operator, app); err != nil {
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
			ActiveStatus:       t.ActiveStatus,
			StartTime:          t.StartTime,
			ProbeTime:          t.ProbeTime,
			Created:            t.Created,
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
	// 检查权限
	if _, _, err := checkAccessConfigPerm(ctx, reqDTO.Operator, reqDTO.ConfigId); err != nil {
		return nil, 0, err
	}
	logs, err := deploymd.ListDeployLog(ctx, deploymd.ListDeployLogReqDTO{
		ConfigId: reqDTO.ConfigId,
		Cursor:   reqDTO.Cursor,
		Limit:    reqDTO.Limit,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(logs, func(t deploymd.DeployLog) (DeployLogDTO, error) {
		return DeployLogDTO{
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
	// 检查权限
	if _, _, err := checkAccessConfigPerm(ctx, reqDTO.Operator, reqDTO.ConfigId); err != nil {
		return nil, 0, err
	}
	logs, err := deploymd.ListOpLog(ctx, deploymd.ListOpLogReqDTO{
		ConfigId: reqDTO.ConfigId,
		Cursor:   reqDTO.Cursor,
		Limit:    reqDTO.Limit,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
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

// DeployServiceWithPlan 通过发布计划部署服务
func (*outerImpl) DeployServiceWithPlan(ctx context.Context, reqDTO DeployServiceWithPlanReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.DeployServiceWithPlan),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		plan   deploymd.Plan
		config deploymd.Config
		item   deploymd.PlanItem
		b      bool
	)
	item, b, err = deploymd.GetPlanItemById(ctx, reqDTO.ItemId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || (item.ItemStatus != deploymd.WaitItemStatus && item.ItemStatus != deploymd.RollbackItemStatus) {
		err = util.InvalidArgsError()
		return
	}
	plan, b, err = deploymd.GetPlanById(ctx, item.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.ThereHasBugErr()
		return
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus || plan.IsClosed() {
		err = util.InvalidArgsError()
		return
	}
	config, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, item.ConfigId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	b, err = deploymd.UpdateItemStatusWithOldStatus(ctx,
		reqDTO.ItemId, reqDTO.Env,
		deploymd.DeployedItemStatus, deploymd.WaitItemStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = deployService(&config, item.ProductVersion, reqDTO.Env, reqDTO.Operator.Account, item.PlanId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	}
	return
}

// RollbackServiceWithPlan 通过发布计划回滚服务
func (*outerImpl) RollbackServiceWithPlan(ctx context.Context, reqDTO RollbackServiceWithPlanReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.DeploySrvKeysVO.RollbackServiceWithPlan),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		plan   deploymd.Plan
		config deploymd.Config
		item   deploymd.PlanItem
		b      bool
	)
	item, b, err = deploymd.GetPlanItemById(ctx, reqDTO.ItemId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b || item.ItemStatus != deploymd.DeployedItemStatus {
		err = util.InvalidArgsError()
		return
	}
	plan, b, err = deploymd.GetPlanById(ctx, item.PlanId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.ThereHasBugErr()
		return
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus || plan.IsClosed() {
		err = util.InvalidArgsError()
		return
	}
	config, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, item.ConfigId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	b, err = deploymd.UpdateItemStatusWithOldStatus(ctx,
		reqDTO.ItemId, reqDTO.Env,
		deploymd.RollbackItemStatus, deploymd.DeployedItemStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = deployService(&config, item.LastProductVersion, reqDTO.Env, reqDTO.Operator.Account, item.PlanId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	}
	return
}

// ListPlan 发布计划列表
func (*outerImpl) ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]PlanDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		_, b, err := teammd.GetTeamUser(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		if !b {
			return nil, 0, util.UnauthorizedError()
		}
	}
	plans, err := deploymd.ListPlan(ctx, deploymd.ListPlanReqDTO{
		TeamId: reqDTO.TeamId,
		Cursor: reqDTO.Cursor,
		Limit:  reqDTO.Limit,
		Env:    reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(plans) == reqDTO.Limit {
		next = plans[len(plans)-1].Id
	}
	data, _ := listutil.Map(plans, func(t deploymd.Plan) (PlanDTO, error) {
		planStatus := t.PlanStatus
		if planStatus == deploymd.RunningPlanStatus && t.IsExpired() {
			planStatus = deploymd.ClosedPlanStatus
			deploymd.UpdatePlanStatusById(ctx, t.Id, reqDTO.Env, deploymd.ClosedPlanStatus)
		}
		return PlanDTO{
			Id:         t.Id,
			Name:       t.Name,
			PlanType:   t.PlanType,
			PlanStatus: planStatus,
			TeamId:     t.TeamId,
			Creator:    t.Creator,
			Expired:    t.Expired,
			Created:    t.Created,
		}, nil
	})
	return data, next, nil
}

type innerImpl struct{}

// DeployServiceWithoutPlan 不通过发布计划部署服务
func (*innerImpl) DeployServiceWithoutPlan(ctx context.Context, reqDTO DeployServiceWithoutPlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查制品
	_, b, err := productmd.GetProduct(ctx, productmd.GetProductReqDTO{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Product,
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
	config, b, err := deploymd.GetEnabledConfigByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	err = deployService(&config, reqDTO.Product, reqDTO.Env, reqDTO.Operator, 0)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}
