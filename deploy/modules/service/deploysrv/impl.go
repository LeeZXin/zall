package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"time"
)

type outerImpl struct{}

func NewOuterService() OuterService {
	initRunner()
	return new(outerImpl)
}

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
			Id:      t.Id,
			AppId:   t.AppId,
			Name:    t.Name,
			Content: t.Content,
			Env:     t.Env,
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

// DeleteConfig 删除配置
func (*outerImpl) DeleteConfig(ctx context.Context, reqDTO DeleteConfigReqDTO) error {
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
	_, err = deploymd.DeleteConfigById(ctx, reqDTO.ConfigId)
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
		AppId:   reqDTO.AppId,
		Name:    reqDTO.Name,
		Content: reqDTO.Content,
		Env:     reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CreatePlan 创建发布计划
func (*outerImpl) CreatePlan(ctx context.Context, reqDTO CreatePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err := checkDeployPlanPerm(ctx, reqDTO.TeamId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
		Name:     reqDTO.Name,
		IsClosed: false,
		TeamId:   reqDTO.TeamId,
		Creator:  reqDTO.Operator.Account,
		Env:      reqDTO.Env,
		Expired:  time.Now().Add(time.Duration(reqDTO.ExpireHours) * time.Hour),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ClosePlan 关闭发布计划
func (*outerImpl) ClosePlan(ctx context.Context, reqDTO ClosePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		plan deploymd.Plan
		b    bool
	)
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || plan.IsExpired() || plan.IsClosed {
		return util.InvalidArgsError()
	}
	if plan.Creator != reqDTO.Operator.Account {
		if err = checkDeployPlanPerm(ctx, plan.TeamId, reqDTO.Operator); err != nil {
			return err
		}
	}
	_, err = deploymd.ClosePlan(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// AddPlanService 添加发布计划部署服务
func (*outerImpl) AddPlanService(ctx context.Context, reqDTO AddPlanServiceReqDTO) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || plan.IsInvalid() {
		return util.InvalidArgsError()
	}
	// 校验是否在其他发布计划存在
	b, err = deploymd.ExistPlanServiceByConfigId(
		ctx,
		reqDTO.ConfigId,
		[]deploymd.ServiceStatus{
			deploymd.RunningServiceStatus,
			deploymd.PendingServiceStatus,
		},
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	// 校验权限
	config, err := checkPlanService(ctx, &plan, reqDTO.ConfigId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = deploymd.InsertPlanService(ctx, deploymd.InsertPlanServiceReqDTO{
		ConfigId:           config.Id,
		CurrProductVersion: reqDTO.CurrProductVersion,
		LastProductVersion: reqDTO.LastProductVersion,
		DeployConfig:       config.Content,
		Status:             deploymd.PendingServiceStatus,
		PlanId:             plan.Id,
	})
	if err != nil {
		return util.InternalError(err)
	}
	return nil
}

// DeletePendingPlanService 删除未执行发布计划单项服务
func (*outerImpl) DeletePendingPlanService(ctx context.Context, reqDTO DeletePendingPlanServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	service, b, err := deploymd.GetServiceById(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.ServiceStatus != deploymd.PendingServiceStatus {
		return util.InvalidArgsError()
	}
	plan, b, err := deploymd.GetPlanById(ctx, service.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	if plan.IsInvalid() {
		return util.InvalidArgsError()
	}
	// 校验权限
	_, _, err = checkAccessConfigPerm(ctx, reqDTO.Operator, service.ConfigId)
	if err != nil {
		return err
	}
	_, err = deploymd.DeletePendingPlanServiceById(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListPlanService 展示发布计划的服务
func (*outerImpl) ListPlanService(ctx context.Context, reqDTO ListPlanServiceReqDTO) ([]PlanServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
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
	services, err := deploymd.ListPlanServiceByPlanId(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	configIdList, _ := listutil.Map(services, func(t deploymd.PlanService) (int64, error) {
		return t.ConfigId, nil
	})
	configs, err := deploymd.BatchGetSimpleConfigById(ctx, configIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	configMap, _ := listutil.CollectToMap(configs, func(t deploymd.Config) (int64, error) {
		return t.Id, nil
	}, func(t deploymd.Config) (deploymd.Config, error) {
		return t, nil
	})
	return listutil.Map(services, func(t deploymd.PlanService) (PlanServiceDTO, error) {
		return PlanServiceDTO{
			Id:                 t.Id,
			AppId:              configMap[t.ConfigId].AppId,
			ConfigId:           t.ConfigId,
			ConfigName:         configMap[t.ConfigId].Name,
			CurrProductVersion: t.CurrProductVersion,
			LastProductVersion: t.LastProductVersion,
			ServiceStatus:      t.ServiceStatus,
			Created:            t.Created,
		}, nil
	})
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
	const pageSize = 10
	plans, total, err := deploymd.ListPlan(ctx, deploymd.ListPlanReqDTO{
		TeamId:   reqDTO.TeamId,
		PageNum:  reqDTO.PageNum,
		PageSize: pageSize,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(plans, func(t deploymd.Plan) (PlanDTO, error) {
		if !t.IsClosed && t.IsExpired() {
			deploymd.ClosePlan(ctx, t.Id)
		}
		return PlanDTO{
			Id:       t.Id,
			Name:     t.Name,
			IsClosed: t.IsInvalid(),
			TeamId:   t.TeamId,
			Creator:  t.Creator,
			Expired:  t.Expired,
			Created:  t.Created,
		}, nil
	})
	return data, total, nil
}

// StartPlanService 启动部署服务流水线
func (*outerImpl) StartPlanService(ctx context.Context, reqDTO StartPlanServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	service, b, err := deploymd.GetServiceById(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.ServiceStatus != deploymd.PendingServiceStatus {
		return util.InvalidArgsError()
	}
	plan, b, err := deploymd.GetPlanById(ctx, service.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	if plan.IsInvalid() {
		return util.InvalidArgsError()
	}
	// 校验权限
	config, _, err := checkAccessConfigPerm(ctx, reqDTO.Operator, service.ConfigId)
	if err != nil {
		return err
	}
	var dp deploy.Deploy
	err = yaml.Unmarshal([]byte(config.Content), &dp)
	if err != nil || !dp.IsValid() {
		return util.ThereHasBugErr()
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b2, err2 := deploymd.UpdateServiceStatusByIdWithOldStatus(
			ctx,
			service.Id,
			deploymd.RunningServiceStatus,
			deploymd.PendingServiceStatus,
		)
		if err2 != nil {
			return err2
		}
		if b2 {
			reqList := make([]deploymd.InsertDeployStepReqDTO, 0)
			for index, stage := range dp.Deploy {
				if len(stage.Agents) > 0 {
					for _, agent := range stage.Agents {
						reqList = append(reqList, deploymd.InsertDeployStepReqDTO{
							ServiceId:  service.Id,
							StepIndex:  index,
							Agent:      agent,
							StepStatus: deploymd.PendingStepStatus,
						})
					}
				} else {
					for _, agent := range dp.Agents {
						reqList = append(reqList, deploymd.InsertDeployStepReqDTO{
							ServiceId:  service.Id,
							StepIndex:  index,
							Agent:      agent.Id,
							StepStatus: deploymd.PendingStepStatus,
						})
					}
				}
			}
			err2 = deploymd.BatchInsertDeployStep(ctx, reqList...)
			if err2 != nil {
				return err2
			}
			return executeDeployOnStartPlanService(service.Id, config.AppId, dp, map[string]string{
				deploy.CurrentProductVersionKey: service.CurrProductVersion,
				deploy.LastProductVersionKey:    service.LastProductVersion,
				deploy.OperatorKey:              reqDTO.Operator.Account,
			})
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// FinishPlanService 完成服务流水线
func (*outerImpl) FinishPlanService(ctx context.Context, reqDTO FinishPlanServiceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	service, b, err := deploymd.GetServiceById(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.ServiceStatus != deploymd.PendingServiceStatus {
		return util.InvalidArgsError()
	}
	plan, b, err := deploymd.GetPlanById(ctx, service.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	if plan.IsInvalid() {
		return util.InvalidArgsError()
	}
	// 校验权限
	config, _, err := checkAccessConfigPerm(ctx, reqDTO.Operator, service.ConfigId)
	if err != nil {
		return err
	}
	var dp deploy.Deploy
	err = yaml.Unmarshal([]byte(config.Content), &dp)
	if err != nil || !dp.IsValid() {
		return util.ThereHasBugErr()
	}
	_, err = deploymd.UpdateServiceStatusByIdWithOldStatus(
		ctx,
		service.Id,
		deploymd.FinishServiceStatus,
		deploymd.RunningServiceStatus,
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ConfirmPlanServiceStep 执行流水线其中一个环节
func (*outerImpl) ConfirmPlanServiceStep(ctx context.Context, reqDTO ConfirmPlanServiceStepReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	service, b, err := deploymd.GetServiceById(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.ServiceStatus != deploymd.RunningServiceStatus {
		return util.InvalidArgsError()
	}
	plan, b, err := deploymd.GetPlanById(ctx, service.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	if plan.IsInvalid() {
		return util.InvalidArgsError()
	}
	// 校验权限
	config, _, err := checkAccessConfigPerm(ctx, reqDTO.Operator, service.ConfigId)
	if err != nil {
		return err
	}
	var dp deploy.Deploy
	err = yaml.Unmarshal([]byte(config.Content), &dp)
	if err != nil || !dp.IsValid() || reqDTO.Index >= len(dp.Deploy) {
		return util.ThereHasBugErr()
	}
	// 找到下标之前的任务 如果还有未成功的任务 则不允许执行
	steps, err := deploymd.ListStepByServiceIdAndLessThanIndex(ctx, reqDTO.ServiceId, reqDTO.Index)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	for _, step := range steps {
		if step.StepStatus != deploymd.SuccessStepStatus {
			return util.OperationFailedError()
		}
	}
	input := make(map[string]string, len(reqDTO.Input)+3)
	input[deploy.CurrentProductVersionKey] = service.CurrProductVersion
	input[deploy.LastProductVersionKey] = service.LastProductVersion
	input[deploy.OperatorKey] = reqDTO.Operator.Account
	for k, v := range reqDTO.Input {
		input[k] = v
	}
	err = executeDeployByIndex(service.Id, config.AppId, dp, reqDTO.Index, input)
	if err != nil {
		// 超出协程池能力范围
		return util.OperationFailedError()
	}
	return nil
}

// RollbackPlanServiceStep 回滚执行流水线其中一个环节
func (*outerImpl) RollbackPlanServiceStep(ctx context.Context, reqDTO RollbackPlanServiceStepReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	step, b, err := deploymd.GetStepByStepId(ctx, reqDTO.StepId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || step.StepStatus != deploymd.SuccessStepStatus {
		return util.InvalidArgsError()
	}
	service, b, err := deploymd.GetServiceById(ctx, step.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || service.ServiceStatus != deploymd.RunningServiceStatus {
		return util.InvalidArgsError()
	}
	plan, b, err := deploymd.GetPlanById(ctx, service.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	if plan.IsInvalid() {
		return util.InvalidArgsError()
	}
	// 校验权限
	config, _, err := checkAccessConfigPerm(ctx, reqDTO.Operator, service.ConfigId)
	if err != nil {
		return err
	}
	var dp deploy.Deploy
	err = yaml.Unmarshal([]byte(config.Content), &dp)
	if err != nil || !dp.IsValid() {
		return util.ThereHasBugErr()
	}
	agentMap := dp.GetAgentMap()
	agent, b := agentMap[step.Agent]
	if !b {
		return util.ThereHasBugErr()
	}
	// 找到下标之前的任务 如果还有未成功的任务 则不允许执行
	steps, err := deploymd.ListStepByServiceIdAndLessThanIndex(ctx, step.ServiceId, step.StepIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	for _, s := range steps {
		if s.StepStatus != deploymd.SuccessStepStatus {
			return util.OperationFailedError()
		}
	}
	input := step.GetInputArgs()
	input[deploy.CurrentProductVersionKey] = service.CurrProductVersion
	input[deploy.LastProductVersionKey] = service.LastProductVersion
	input[deploy.OperatorKey] = reqDTO.Operator.Account
	err = rollbackStage(service.Id, config.AppId, dp, step.StepIndex, agent, input)
	if err != nil {
		// 超出协程池能力范围
		return util.OperationFailedError()
	}
	return nil
}
