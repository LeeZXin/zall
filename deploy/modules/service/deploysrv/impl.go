package deploysrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"math"
)

type outerImpl struct{}

func NewOuterService() OuterService {
	initRunner()
	return new(outerImpl)
}

// CreatePlan 创建发布计划
func (*outerImpl) CreatePlan(ctx context.Context, reqDTO CreatePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	service, err := checkDeployPlanPermByServiceId(ctx, reqDTO.Operator, reqDTO.ServiceId)
	if err != nil {
		return err
	}
	// 检查是否有其他发布计划在
	b, err := deploymd.ExistPendingOrRunningPlanByServiceId(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	_, err = deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
		Name:           reqDTO.Name,
		PlanStatus:     deploymd.PendingPlanStatus,
		AppId:          service.AppId,
		ServiceId:      reqDTO.ServiceId,
		ProductVersion: reqDTO.ProductVersion,
		Creator:        reqDTO.Operator.Account,
		Env:            service.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// StartPlan 开始发布计划
func (*outerImpl) StartPlan(ctx context.Context, reqDTO StartPlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || plan.PlanStatus != deploymd.PendingPlanStatus {
		return util.InvalidArgsError()
	}
	service, err := checkAppDevelopPermByServiceId(ctx, reqDTO.Operator, plan.ServiceId)
	if err != nil {
		return err
	}
	var serviceConfig deploy.Service
	err = yaml.Unmarshal([]byte(service.Config), &serviceConfig)
	if err != nil || !serviceConfig.IsValid() {
		return util.ThereHasBugErr()
	}
	insertStageList, taskIdMapList := serviceConfig2DeployStageReq(reqDTO.PlanId, &serviceConfig)
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b2, err2 := deploymd.StartPlan(ctx, reqDTO.PlanId, service.Config)
		if err2 != nil {
			return err2
		}
		if !b2 {
			return fmt.Errorf("start plan failed: %v", reqDTO.PlanId)
		}
		// 新增部署步骤
		err2 = deploymd.BatchInsertDeployStage(ctx, insertStageList...)
		if err2 != nil {
			return err2
		}
		return executeDeployOnStartPlan(
			reqDTO.PlanId,
			service.AppId,
			serviceConfig,
			map[string]string{
				deploy.CurrentProductVersionKey: plan.ProductVersion,
				deploy.OperatorAccountKey:       reqDTO.Operator.Account,
				deploy.AppKey:                   plan.AppId,
			},
			taskIdMapList,
		)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	return nil
}

// RedoAgentStage 重新执行agent
func (*outerImpl) RedoAgentStage(ctx context.Context, reqDTO RedoAgentStageReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	stage, plan, _, err := checkAppDevelopPermByStageId(ctx, reqDTO.Operator, reqDTO.StageId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus ||
		stage.StageStatus != deploymd.FailedStageStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Service
	err = yaml.Unmarshal([]byte(plan.ServiceConfig), &dp)
	if err != nil {
		return nil
	}
	b, err := deploymd.UpdateStageStatusWithOldStatusById(ctx, stage.Id, deploymd.PendingStageStatus, stage.StageStatus)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		args := make(map[string]string)
		if stage.InputArgs != nil && stage.InputArgs.Data != nil {
			for k, v := range stage.InputArgs.Data {
				args[k] = v
			}
		}
		args[deploy.CurrentProductVersionKey] = plan.ProductVersion
		args[deploy.OperatorAccountKey] = reqDTO.Operator.Account
		args[deploy.AppKey] = plan.AppId
		err = runner.Execute(func() {
			redoAgentStage(stage.PlanId, plan.AppId, dp, stage.StageIndex, stage.Agent, args, stage.TaskId)
		})
		if err != nil {
			return util.OperationFailedError()
		}
	}
	return nil
}

// ForceRedoNotSuccessfulAgentStages 强制重新执行未完成的任务
func (*outerImpl) ForceRedoNotSuccessfulAgentStages(ctx context.Context, reqDTO ForceRedoNotSuccessfulAgentStagesReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Service
	err = yaml.Unmarshal([]byte(plan.ServiceConfig), &dp)
	if err != nil {
		return nil
	}
	if reqDTO.StageIndex >= len(dp.Deploy) {
		return util.InvalidArgsError()
	}
	stages, err := deploymd.ListStagesByPlanIdAndLETIndex(ctx, reqDTO.PlanId, reqDTO.StageIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 不太可能会发生
	if len(stages) == 0 {
		return nil
	}
	notSuccessfulStages := make([]deploymd.Stage, 0)
	for _, s := range stages {
		// 前面阶段都必须成功
		if s.StageIndex < reqDTO.StageIndex && s.StageStatus != deploymd.SuccessfulStageStatus {
			return util.InvalidArgsError()
		}
		// 计算未成功数量
		if s.StageIndex == reqDTO.StageIndex && s.StageStatus != deploymd.SuccessfulStageStatus {
			notSuccessfulStages = append(notSuccessfulStages, s)
		}
	}
	if len(notSuccessfulStages) == 0 {
		return util.InvalidArgsError()
	}
	stage := dp.Deploy[reqDTO.StageIndex]
	args, err := deploymd.MergeInputArgsByPlanIdAndLTIndex(ctx, reqDTO.PlanId, reqDTO.StageIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	args[deploy.CurrentProductVersionKey] = plan.ProductVersion
	args[deploy.OperatorAccountKey] = reqDTO.Operator.Account
	args[deploy.AppKey] = plan.AppId
	// 需要交互
	if stage.Confirm != nil && stage.Confirm.NeedInteract {
		b, extra := stage.Confirm.CheckForm(reqDTO.Args)
		if !b {
			return util.InvalidArgsError()
		}
		for k, v := range extra {
			args[k] = v
		}
	}
	err = runner.Execute(func() {
		for _, ns := range notSuccessfulStages {
			b, err := updateStageStatusWithOldStatusById(ns.Id, deploymd.PendingStageStatus, ns.StageStatus)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return
			}
			if !b {
				break
			}
			redoAgentStage(reqDTO.PlanId, plan.AppId, dp, reqDTO.StageIndex, ns.Agent, args, ns.TaskId)
		}
	})
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

func updateStageStatusWithOldStatusById(stageId int64, newStatus, oldStatus deploymd.StageStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusWithOldStatusById(ctx, stageId, newStatus, oldStatus)
	return b, err
}

func serviceConfig2DeployStageReq(planId int64, s *deploy.Service) ([]deploymd.InsertDeployStageReqDTO, []map[string]string) {
	ret := make([]deploymd.InsertDeployStageReqDTO, 0)
	taskIdMapList := make([]map[string]string, 0)
	for i, stage := range s.Deploy {
		agents := make([]string, 0)
		if len(stage.Agents) > 0 {
			for _, agent := range stage.Agents {
				agents = append(agents, agent)
			}
		} else {
			for agent := range s.Agents {
				agents = append(agents, agent)
			}
		}
		taskIdMap := make(map[string]string)
		for _, agent := range agents {
			taskId := idutil.RandomUuid()
			ret = append(ret, deploymd.InsertDeployStageReqDTO{
				PlanId:      planId,
				StageIndex:  i,
				Agent:       agent,
				TaskId:      taskId,
				StageStatus: deploymd.PendingStageStatus,
			})
			taskIdMap[agent] = taskId
		}
		taskIdMapList = append(taskIdMapList, taskIdMap)
	}
	return ret, taskIdMapList
}

// ClosePlan 关闭发布计划
func (*outerImpl) ClosePlan(ctx context.Context, reqDTO ClosePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, b, err := deploymd.GetPlanById(ctx, reqDTO.PlanId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || plan.PlanStatus.IsFinalStatus() {
		return util.InvalidArgsError()
	}
	service, err := checkAppDevelopPermByServiceId(ctx, reqDTO.Operator, plan.ServiceId)
	if err != nil {
		return err
	}
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		_, err = deploymd.ClosePlanAndUpdateConfig(ctx, reqDTO.PlanId, plan.PlanStatus, service.Config)
	} else {
		// 判断是否有执行中的任务
		b, err = deploymd.ExistRunningStatusByPlanId(ctx, reqDTO.PlanId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if b {
			// todo 替换返回说明
			return util.OperationFailedError()
		}
		_, err = deploymd.UpdatePlanStatusWithOldStatus(ctx, reqDTO.PlanId, deploymd.ClosedPlanStatus, plan.PlanStatus)
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListPlan 发布计划列表
func (*outerImpl) ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]PlanDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, 0, err
	}
	const pageSize = 10
	plans, total, err := deploymd.ListPlan(ctx, deploymd.ListPlanReqDTO{
		AppId:    reqDTO.AppId,
		PageNum:  reqDTO.PageNum,
		PageSize: pageSize,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	serviceMap := make(map[int64]string, len(plans))
	serviceIdList, _ := listutil.Map(plans, func(t deploymd.Plan) (int64, error) {
		return t.ServiceId, nil
	})
	serviceIdList = listutil.Distinct(serviceIdList...)
	if len(serviceIdList) > 0 {
		services, err := deploymd.BatchGetServiceById(ctx, serviceIdList, []string{"id", "name"})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		for _, srv := range services {
			serviceMap[srv.Id] = srv.Name
		}
	}
	data, _ := listutil.Map(plans, func(t deploymd.Plan) (PlanDTO, error) {
		return PlanDTO{
			Id:             t.Id,
			ServiceId:      t.ServiceId,
			ServiceName:    serviceMap[t.ServiceId],
			Name:           t.Name,
			ProductVersion: t.ProductVersion,
			PlanStatus:     t.PlanStatus,
			Env:            t.Env,
			Creator:        t.Creator,
			Created:        t.Created,
		}, nil
	})
	return data, total, nil
}

// GetPlanDetail 获取单个发布计划详情
func (*outerImpl) GetPlanDetail(ctx context.Context, reqDTO GetPlanDetailReqDTO) (PlanDetailDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return PlanDetailDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, service, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return PlanDetailDTO{}, err
	}
	ret := PlanDetailDTO{
		Id:             plan.Id,
		ServiceId:      plan.ServiceId,
		ServiceName:    service.Name,
		Name:           plan.Name,
		ProductVersion: plan.ProductVersion,
		PlanStatus:     plan.PlanStatus,
		Env:            plan.Env,
		Creator:        plan.Creator,
		Created:        plan.Created,
	}
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		ret.ServiceConfig = service.Config
	} else {
		ret.ServiceConfig = plan.ServiceConfig
	}
	return ret, nil
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
	err := deploymd.InsertService(ctx, deploymd.InsertServiceReqDTO{
		AppId:       reqDTO.AppId,
		Config:      reqDTO.Config,
		Env:         reqDTO.Env,
		Name:        reqDTO.Name,
		ServiceType: reqDTO.service.Type,
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
	_, err := deploymd.UpdateService(ctx, deploymd.UpdateServiceReqDTO{
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
	// 存在正在执行的发布计划
	b, err := deploymd.ExistPendingOrRunningPlanByServiceId(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.AlreadyExistsError()
	}
	_, err = deploymd.DeleteService(ctx, reqDTO.ServiceId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListService 服务列表
func (*outerImpl) ListService(ctx context.Context, reqDTO ListServiceReqDTO) ([]ServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkServicePerm(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	services, err := deploymd.ListService(ctx, deploymd.ListServiceReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(services, func(t deploymd.Service) (ServiceDTO, error) {
		return ServiceDTO{
			Id:          t.Id,
			AppId:       t.AppId,
			Config:      t.Config,
			Env:         t.Env,
			Name:        t.Name,
			ServiceType: t.ServiceType,
		}, nil
	})
	return data, nil
}

// ListServiceWhenCreatePlan 创建发布计划时展示的服务列表
func (*outerImpl) ListServiceWhenCreatePlan(ctx context.Context, reqDTO ListServiceWhenCreatePlanReqDTO) ([]SimpleServiceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkDeployPlanPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, err
	}
	services, err := deploymd.ListService(ctx, deploymd.ListServiceReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
		Cols:  []string{"id", "env", "name", "service_type"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(services, func(t deploymd.Service) (SimpleServiceDTO, error) {
		return SimpleServiceDTO{
			Id:          t.Id,
			Env:         t.Env,
			Name:        t.Name,
			ServiceType: t.ServiceType,
		}, nil
	})
	return data, nil
}

// ListStages 展示发布计划流水线详情
func (*outerImpl) ListStages(ctx context.Context, reqDTO ListStagesReqDTO) ([]StageDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, srv, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return nil, err
	}
	var serviceConfig string
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		serviceConfig = srv.Config
	} else {
		serviceConfig = plan.ServiceConfig
	}
	var service deploy.Service
	err = yaml.Unmarshal([]byte(serviceConfig), &service)
	if err != nil || !service.IsValid() {
		return nil, nil
	}
	stagesMap := make(map[int][]deploymd.Stage)
	if plan.PlanStatus != deploymd.PendingPlanStatus {
		// 只有开始执行发布计划后才有下面的数据
		stages, err := deploymd.GetStageListByPlanId(ctx, reqDTO.PlanId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		for _, stage := range stages {
			sl, has := stagesMap[stage.StageIndex]
			if has {
				stagesMap[stage.StageIndex] = append(sl, stage)
			} else {
				stagesMap[stage.StageIndex] = append(make([]deploymd.Stage, 0), stage)
			}
		}
	}
	ret := make([]StageDTO, 0, len(stagesMap))
	for index, stage := range service.Deploy {
		var (
			total, done, percent, pending, running float64

			hasError bool
		)
		total = float64(len(stage.Agents))
		if total == 0 {
			total = float64(len(service.Agents))
		}
		subStages := make([]SubStageDTO, 0, len(stagesMap[index]))
		if len(stagesMap) > 0 {
			for _, md := range stagesMap[index] {
				sa := SubStageDTO{
					Id:         md.Id,
					Agent:      md.Agent,
					ExecuteLog: md.ExecuteLog,
				}
				if plan.PlanStatus != deploymd.ClosedPlanStatus || md.StageStatus != deploymd.PendingStageStatus {
					sa.StageStatus = md.StageStatus
				}
				if len(service.Agents) > 0 {
					sa.AgentHost = service.Agents[md.Agent].Host
				}
				subStages = append(subStages, sa)
				switch md.StageStatus {
				case deploymd.SuccessfulStageStatus:
					done++
				case deploymd.FailedStageStatus:
					hasError = true
				case deploymd.PendingStageStatus:
					pending++
				case deploymd.RunningStageStatus:
					running++
				}
			}
		} else {
			if len(stage.Agents) > 0 {
				for _, agent := range stage.Agents {
					sa := SubStageDTO{
						Agent: agent,
					}
					if len(service.Agents) > 0 {
						sa.AgentHost = service.Agents[agent].Host
					}
					subStages = append(subStages, sa)
				}
			} else {
				for id, agent := range service.Agents {
					sa := SubStageDTO{
						Agent:     id,
						AgentHost: agent.Host,
					}
					subStages = append(subStages, sa)
				}
			}
		}
		if total > 0 {
			percent = math.Round(done / total * 100)
		}
		waitInteract := plan.PlanStatus == deploymd.RunningPlanStatus &&
			pending == total && stage.Confirm != nil && stage.Confirm.NeedInteract
		if index > 0 {
			waitInteract = ret[index-1].IsAllDone && waitInteract
		}
		// 是否可以强制重新执行未成功agent
		canForceRedoUnSuccessAgentStages := !waitInteract && plan.PlanStatus == deploymd.RunningPlanStatus &&
			done < total
		if index > 0 {
			canForceRedoUnSuccessAgentStages = ret[index-1].IsAllDone && canForceRedoUnSuccessAgentStages
		}
		isRunning := plan.PlanStatus == deploymd.RunningPlanStatus &&
			running > 0
		isAllDone := done == total
		dto := StageDTO{
			Name:                             stage.Name,
			Percent:                          percent,
			Total:                            int(total),
			Done:                             int(done),
			IsAutomatic:                      stage.Confirm == nil || !stage.Confirm.NeedInteract,
			HasError:                         hasError,
			IsRunning:                        isRunning,
			WaitInteract:                     waitInteract,
			SubStages:                        subStages,
			IsAllDone:                        isAllDone,
			CanForceRedoUnSuccessAgentStages: canForceRedoUnSuccessAgentStages,
		}
		if len(service.Actions) > 0 {
			dto.Script = service.Actions[stage.Action].Script
		}
		dto.Confirm = stage.Confirm
		ret = append(ret, dto)
	}
	return ret, nil
}

// KillStage 中止执行
func (*outerImpl) KillStage(ctx context.Context, reqDTO KillStageReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Service
	err = yaml.Unmarshal([]byte(plan.ServiceConfig), &dp)
	if err != nil {
		return nil
	}
	if reqDTO.StageIndex >= len(dp.Deploy) {
		return util.InvalidArgsError()
	}
	stages, err := deploymd.ListStageByPlanIdAndStageIndexAndStatus(
		ctx,
		reqDTO.PlanId,
		reqDTO.StageIndex,
		deploymd.RunningStageStatus,
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	err = deploymd.KillStages(ctx, reqDTO.PlanId, reqDTO.StageIndex, "kill by "+reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	err = runner.Execute(func() {
		for _, stage := range stages {
			agent := dp.Agents[stage.Agent]
			err = sshagent.NewServiceCommand(agent.Host, agent.Token, plan.AppId).
				Kill(stage.TaskId)
			logger.Logger.Errorf("kill taskId: %s with err: %v", stage.TaskId, err)
		}
	})
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

// ConfirmInteractStage 交互阶段确认
func (*outerImpl) ConfirmInteractStage(ctx context.Context, reqDTO ConfirmInteractStageReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Service
	err = yaml.Unmarshal([]byte(plan.ServiceConfig), &dp)
	if err != nil {
		return nil
	}
	if reqDTO.StageIndex >= len(dp.Deploy) {
		return util.InvalidArgsError()
	}
	stage := dp.Deploy[reqDTO.StageIndex]
	// 判断不是交互阶段
	if stage.Confirm == nil || !stage.Confirm.NeedInteract {
		return util.InvalidArgsError()
	}
	args := reqDTO.Args
	// 防止非空指针
	if args == nil {
		args = make(map[string]string, 0)
	}
	filteredArgs, err := deploymd.MergeInputArgsByPlanIdAndLTIndex(ctx, reqDTO.PlanId, reqDTO.StageIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	filteredArgs[deploy.CurrentProductVersionKey] = plan.ProductVersion
	filteredArgs[deploy.OperatorAccountKey] = reqDTO.Operator.Account
	filteredArgs[deploy.AppKey] = plan.AppId
	if len(stage.Confirm.Form) > 0 {
		b, extra := stage.Confirm.CheckForm(args)
		if !b {
			return util.InvalidArgsError()
		}
		for k, v := range extra {
			filteredArgs[k] = v
		}
	}
	stages, err := deploymd.ListStagesByPlanIdAndLETIndex(ctx, reqDTO.PlanId, reqDTO.StageIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 不太可能会发生
	if len(stages) == 0 {
		return nil
	}
	for _, s := range stages {
		// 前面阶段都必须成功
		if s.StageIndex < reqDTO.StageIndex && s.StageStatus != deploymd.SuccessfulStageStatus {
			return util.InvalidArgsError()
		}
		// 当前阶段都属于pending状态
		if s.StageIndex == reqDTO.StageIndex && s.StageStatus != deploymd.PendingStageStatus {
			return util.InvalidArgsError()
		}
	}
	err = executeDeployOnConfirmStage(reqDTO.PlanId, plan.AppId, dp, filteredArgs, reqDTO.StageIndex)
	if err != nil {
		return util.OperationFailedError()
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
	probe, b, err := deploymd.GetServiceById(ctx, probeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkServicePerm(ctx, probe.AppId, operator)
}
