package deploysrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/pkg/status"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"io"
	"math"
	"net/http"
	"strings"
)

// CreatePlan 创建发布计划
func CreatePlan(ctx context.Context, reqDTO CreatePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	pipeline, err := checkDeployPlanPermByPipelineId(ctx, reqDTO.Operator, reqDTO.PipelineId)
	if err != nil {
		return err
	}
	// 检查是否有其他发布计划在
	b, err := deploymd.ExistPendingOrRunningPlanByPipelineId(ctx, reqDTO.PipelineId)
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
		AppId:          pipeline.AppId,
		PipelineId:     reqDTO.PipelineId,
		PipelineName:   pipeline.Name,
		ProductVersion: reqDTO.ProductVersion,
		Creator:        reqDTO.Operator.Account,
		Env:            pipeline.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// StartPlan 开始发布计划
func StartPlan(ctx context.Context, reqDTO StartPlanReqDTO) error {
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
	pipeline, err := checkAppDevelopPermByPipelineId(ctx, reqDTO.Operator, plan.PipelineId)
	if err != nil {
		return err
	}
	var pipelineConfig deploy.Pipeline
	err = yaml.Unmarshal([]byte(pipeline.Config), &pipelineConfig)
	if err != nil || !pipelineConfig.IsValid() {
		return util.ThereHasBugErr()
	}
	insertStageList, taskIdMapList := pipelineConfig2DeployStageReq(reqDTO.PlanId, &pipelineConfig)
	varsMap, err := deploymd.ListPipelineVarsMap(ctx, pipeline.AppId, pipeline.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b2, err2 := deploymd.StartPlan(ctx, reqDTO.PlanId, pipeline.Config)
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
			pipeline.AppId,
			pipelineConfig,
			map[string]string{
				deploy.CurrentProductVersionKey: plan.ProductVersion,
				deploy.OperatorAccountKey:       reqDTO.Operator.Account,
				deploy.AppKey:                   plan.AppId,
			},
			taskIdMapList,
			varsMap,
		)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	return nil
}

// RedoAgentStage 重新执行agent
func RedoAgentStage(ctx context.Context, reqDTO RedoAgentStageReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	stage, plan, pipeline, err := checkAppDevelopPermByStageId(ctx, reqDTO.Operator, reqDTO.StageId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus ||
		stage.StageStatus != deploymd.FailedStageStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Pipeline
	err = yaml.Unmarshal([]byte(plan.PipelineConfig), &dp)
	if err != nil {
		return nil
	}
	varsMap, err := deploymd.ListPipelineVarsMap(ctx, pipeline.AppId, pipeline.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
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
		initRunner()
		err = runner.Execute(func() {
			redoAgentStage(stage.PlanId, plan.AppId, dp, stage.StageIndex, stage.Agent, args, varsMap, stage.TaskId)
		})
		if err != nil {
			return util.OperationFailedError()
		}
	}
	return nil
}

// ForceRedoNotSuccessfulAgentStages 强制重新执行未完成的任务
func ForceRedoNotSuccessfulAgentStages(ctx context.Context, reqDTO ForceRedoNotSuccessfulAgentStagesReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, pipeline, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Pipeline
	err = yaml.Unmarshal([]byte(plan.PipelineConfig), &dp)
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
	varsMap, err := deploymd.ListPipelineVarsMap(ctx, pipeline.AppId, pipeline.Env)
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
	initRunner()
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
			redoAgentStage(reqDTO.PlanId, plan.AppId, dp, reqDTO.StageIndex, ns.Agent, args, varsMap, ns.TaskId)
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

func pipelineConfig2DeployStageReq(planId int64, s *deploy.Pipeline) ([]deploymd.InsertDeployStageReqDTO, []map[string]string) {
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
func ClosePlan(ctx context.Context, reqDTO ClosePlanReqDTO) error {
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
	pipeline, err := checkAppDevelopPermByPipelineId(ctx, reqDTO.Operator, plan.PipelineId)
	if err != nil {
		return err
	}
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		_, err = deploymd.ClosePlanAndUpdateConfig(ctx, reqDTO.PlanId, plan.PlanStatus, pipeline.Config)
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
func ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]PlanDTO, int64, error) {
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
	data, _ := listutil.Map(plans, func(t deploymd.Plan) (PlanDTO, error) {
		return PlanDTO{
			Id:             t.Id,
			PipelineId:     t.PipelineId,
			PipelineName:   t.PipelineName,
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
func GetPlanDetail(ctx context.Context, reqDTO GetPlanDetailReqDTO) (PlanDetailDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return PlanDetailDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, pipeline, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return PlanDetailDTO{}, err
	}
	ret := PlanDetailDTO{
		Id:             plan.Id,
		PipelineId:     plan.PipelineId,
		PipelineName:   pipeline.Name,
		Name:           plan.Name,
		ProductVersion: plan.ProductVersion,
		PlanStatus:     plan.PlanStatus,
		Env:            plan.Env,
		Creator:        plan.Creator,
		Created:        plan.Created,
	}
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		ret.PipelineConfig = pipeline.Config
	} else {
		ret.PipelineConfig = plan.PipelineConfig
	}
	return ret, nil
}

// CreatePipeline 创建流水线
func CreatePipeline(ctx context.Context, reqDTO CreatePipelineReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return err
	}
	err := deploymd.InsertPipeline(ctx, deploymd.InsertPipelineReqDTO{
		AppId:  reqDTO.AppId,
		Config: reqDTO.Config,
		Env:    reqDTO.Env,
		Name:   reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdatePipeline 编辑流水线
func UpdatePipeline(ctx context.Context, reqDTO UpdatePipelineReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkPipelinePermByPipelineId(ctx, reqDTO.PipelineId, reqDTO.Operator); err != nil {
		return err
	}
	_, err := deploymd.UpdatePipeline(ctx, deploymd.UpdatePipelineReqDTO{
		PipelineId: reqDTO.PipelineId,
		Name:       reqDTO.Name,
		Config:     reqDTO.Config,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeletePipeline 删除流水线
func DeletePipeline(ctx context.Context, reqDTO DeletePipelineReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkPipelinePermByPipelineId(ctx, reqDTO.PipelineId, reqDTO.Operator); err != nil {
		return err
	}
	// 存在正在执行的发布计划
	b, err := deploymd.ExistPendingOrRunningPlanByPipelineId(ctx, reqDTO.PipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	_, err = deploymd.DeletePipelineById(ctx, reqDTO.PipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListPipeline 流水线列表
func ListPipeline(ctx context.Context, reqDTO ListPipelineReqDTO) ([]PipelineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	pipelines, err := deploymd.ListPipeline(ctx, deploymd.ListPipelineReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(pipelines, func(t deploymd.Pipeline) (PipelineDTO, error) {
		return PipelineDTO{
			Id:     t.Id,
			AppId:  t.AppId,
			Config: t.Config,
			Env:    t.Env,
			Name:   t.Name,
		}, nil
	})
	return data, nil
}

// ListPipelineWhenCreatePlan 创建发布计划时展示的流水线列表
func ListPipelineWhenCreatePlan(ctx context.Context, reqDTO ListPipelineWhenCreatePlanReqDTO) ([]SimplePipelineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err := checkDeployPlanPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId); err != nil {
		return nil, err
	}
	pipelines, err := deploymd.ListPipeline(ctx, deploymd.ListPipelineReqDTO{
		AppId: reqDTO.AppId,
		Env:   reqDTO.Env,
		Cols:  []string{"id", "env", "name"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(pipelines, func(t deploymd.Pipeline) (SimplePipelineDTO, error) {
		return SimplePipelineDTO{
			Id:   t.Id,
			Env:  t.Env,
			Name: t.Name,
		}, nil
	})
	return data, nil
}

// ListStages 展示发布计划流水线详情
func ListStages(ctx context.Context, reqDTO ListStagesReqDTO) ([]StageDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, srv, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return nil, err
	}
	var pipelineConfig string
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		pipelineConfig = srv.Config
	} else {
		pipelineConfig = plan.PipelineConfig
	}
	var pipeline deploy.Pipeline
	err = yaml.Unmarshal([]byte(pipelineConfig), &pipeline)
	if err != nil || !pipeline.IsValid() {
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
	for index, stage := range pipeline.Deploy {
		var (
			total, done, percent, pending, running float64

			hasError bool
		)
		total = float64(len(stage.Agents))
		if total == 0 {
			total = float64(len(pipeline.Agents))
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
				if len(pipeline.Agents) > 0 {
					sa.AgentHost = pipeline.Agents[md.Agent].Host
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
					if len(pipeline.Agents) > 0 {
						sa.AgentHost = pipeline.Agents[agent].Host
					}
					subStages = append(subStages, sa)
				}
			} else {
				for id, agent := range pipeline.Agents {
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
		if len(pipeline.Actions) > 0 {
			dto.Script = pipeline.Actions[stage.Action].Script
		}
		dto.Confirm = stage.Confirm
		ret = append(ret, dto)
	}
	return ret, nil
}

// KillStage 中止执行
func KillStage(ctx context.Context, reqDTO KillStageReqDTO) error {
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
	var dp deploy.Pipeline
	err = yaml.Unmarshal([]byte(plan.PipelineConfig), &dp)
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
	initRunner()
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
func ConfirmInteractStage(ctx context.Context, reqDTO ConfirmInteractStageReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, pipeline, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return err
	}
	if plan.PlanStatus != deploymd.RunningPlanStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Pipeline
	err = yaml.Unmarshal([]byte(plan.PipelineConfig), &dp)
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
	varsMap, err := deploymd.ListPipelineVarsMap(ctx, pipeline.AppId, pipeline.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	err = executeDeployOnConfirmStage(reqDTO.PlanId, plan.AppId, dp, filteredArgs, varsMap, reqDTO.StageIndex)
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

// ListServiceSource 查看服务数据来源
func ListServiceSource(ctx context.Context, reqDTO ListServiceSourceReqDTO) ([]ServiceSourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	sources, err := deploymd.ListServiceSource(ctx, deploymd.ListServiceSourceReqDTO{
		Env: reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(sources, func(t deploymd.ServiceSource) (ServiceSourceDTO, error) {
		return ServiceSourceDTO{
			Id:      t.Id,
			Name:    t.Name,
			Env:     t.Env,
			Host:    t.Host,
			ApiKey:  t.ApiKey,
			Created: t.Created,
		}, nil
	})
	return data, nil
}

// ListAllServiceSource 所有服务状态来源
func ListAllServiceSource(ctx context.Context, reqDTO ListAllServiceSourceReqDTO) ([]SimpleServiceSourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	sources, err := deploymd.ListServiceSource(ctx, deploymd.ListServiceSourceReqDTO{
		Env:  reqDTO.Env,
		Cols: []string{"id", "name"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(sources, func(t deploymd.ServiceSource) (SimpleServiceSourceDTO, error) {
		return SimpleServiceSourceDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
	return data, nil
}

// CreateServiceSource 插入服务数据来源
func CreateServiceSource(ctx context.Context, reqDTO CreateServiceSourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	err := deploymd.InsertServiceSource(ctx, deploymd.InsertServiceSourceReqDTO{
		Name:   reqDTO.Name,
		Env:    reqDTO.Env,
		Host:   reqDTO.Host,
		ApiKey: reqDTO.ApiKey,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateServiceSource 更新数据服务来源
func UpdateServiceSource(ctx context.Context, reqDTO UpdateServiceSourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	_, err := deploymd.UpdateServiceSource(ctx, deploymd.UpdateServiceSourceReqDTO{
		Id:     reqDTO.SourceId,
		Name:   reqDTO.Name,
		Host:   reqDTO.Host,
		ApiKey: reqDTO.ApiKey,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteServiceSource 删除数据服务来源
func DeleteServiceSource(ctx context.Context, reqDTO DeleteServiceSourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := deploymd.DeleteServiceSourceById(ctx, reqDTO.SourceId)
		if err2 != nil {
			return err2
		}
		return deploymd.DeleteAppServiceSourceBindBySourceId(ctx, reqDTO.SourceId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListPipelineVars 流水线变量
func ListPipelineVars(ctx context.Context, reqDTO ListPipelineVarsReqDTO) ([]PipelineVarsWithoutContentDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	vars, err := deploymd.ListPipelineVars(ctx, reqDTO.AppId, reqDTO.Env, []string{"id", "app_id", "env", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(vars, func(t deploymd.PipelineVars) (PipelineVarsWithoutContentDTO, error) {
		return PipelineVarsWithoutContentDTO{
			Id:    t.Id,
			Name:  t.Name,
			AppId: t.AppId,
			Env:   t.Env,
		}, nil
	})
}

// CreatePipelineVars 创建流水线变量
func CreatePipelineVars(ctx context.Context, reqDTO CreatePipelineVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return err
	}
	b, err := deploymd.ExistPipelineVarsByAppIdAndEnvAndName(ctx, reqDTO.AppId, reqDTO.Env, reqDTO.Name)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = deploymd.InsertPipelineVars(ctx, deploymd.InsertPipelineVarsReqDTO{
		AppId:   reqDTO.AppId,
		Env:     reqDTO.Env,
		Name:    reqDTO.Name,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdatePipelineVars 编辑流水线变量
func UpdatePipelineVars(ctx context.Context, reqDTO UpdatePipelineVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if _, err := checkManagePipelinePermByVarsId(ctx, reqDTO.Id, reqDTO.Operator); err != nil {
		return err
	}
	_, err := deploymd.UpdatePipelineVars(ctx, deploymd.UpdatePipelineVarsReqDTO{
		Id:      reqDTO.Id,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeletePipelineVars 删除流水线变量
func DeletePipelineVars(ctx context.Context, reqDTO DeletePipelineVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	if _, err := checkManagePipelinePermByVarsId(ctx, reqDTO.Id, reqDTO.Operator); err != nil {
		return err
	}
	_, err := deploymd.DeletePipelineVarsById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetPipelineVarsContent 获取流水线变量内容
func GetPipelineVarsContent(ctx context.Context, reqDTO GetPipelineVarsContentReqDTO) (PipelineVarsDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return PipelineVarsDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	vars, err := checkManagePipelinePermByVarsId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return PipelineVarsDTO{}, err
	}
	return PipelineVarsDTO{
		Id:      vars.Id,
		Name:    vars.Name,
		AppId:   vars.AppId,
		Env:     vars.Env,
		Content: vars.Content,
	}, nil
}

// ListServiceStatus 展示服务状态列表
func ListServiceStatus(ctx context.Context, reqDTO ListServiceStatusReqDTO) ([]status.Service, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	appId, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return nil, err
	}
	url := strings.TrimSuffix(source.Host, "/") +
		fmt.Sprintf("/api/service/v1/status/list?app=%s&env=%s", appId, source.Env)
	resp := make([]status.Service, 0)
	err = httputil.Get(
		ctx,
		http.DefaultClient,
		url,
		map[string]string{
			"Authorization": source.ApiKey,
		},
		&resp,
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.OperationFailedError()
	}
	return resp, nil
}

// DoStatusAction 操作服务
func DoStatusAction(ctx context.Context, reqDTO DoStatusActionReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return err
	}
	actions, err := listActions(ctx, source)
	if err != nil {
		return util.OperationFailedError()
	}
	for _, action := range actions {
		if action.Label == reqDTO.Action {
			url := strings.TrimSuffix(source.Host, "/") +
				fmt.Sprintf("/%s?serviceId=%s", strings.TrimPrefix(action.Api.Url, "/"), reqDTO.ServiceId)
			request, err := http.NewRequestWithContext(ctx, action.Api.Method, url, nil)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return util.OperationFailedError()
			}
			for k, v := range action.Api.Headers {
				request.Header.Set(k, v)
			}
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return util.OperationFailedError()
			}
			all, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				logger.Logger.WithContext(ctx).Infof("url: %s message: %s", url, string(all))
				return util.OperationFailedError()
			}
			return nil
		}
	}
	return util.OperationFailedError()
}

// ListStatusActions 获取服务操作列表
func ListStatusActions(ctx context.Context, reqDTO ListStatusActionReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return nil, err
	}
	actions, err := listActions(ctx, source)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.OperationFailedError()
	}
	return listutil.Map(actions, func(t status.Action) (string, error) {
		return t.Label, nil
	})
}

// ListBindServiceSource 获取绑定服务来源
func ListBindServiceSource(ctx context.Context, reqDTO ListBindServiceSourceReqDTO) ([]SimpleBindServiceSourceDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return nil, err
	}
	binds, err := deploymd.ListAppServiceSourceBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if len(binds) == 0 {
		return []SimpleBindServiceSourceDTO{}, nil
	}
	bindMap := make(map[int64]deploymd.AppServiceSourceBind, len(binds))
	sourceIdList, _ := listutil.Map(binds, func(t deploymd.AppServiceSourceBind) (int64, error) {
		bindMap[t.SourceId] = t
		return t.SourceId, nil
	})
	sources, err := deploymd.BatchGetServiceSourceByIdList(ctx, sourceIdList, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(sources, func(t deploymd.ServiceSource) (SimpleBindServiceSourceDTO, error) {
		bind := bindMap[t.Id]
		return SimpleBindServiceSourceDTO{
			Id:     t.Id,
			Name:   t.Name,
			BindId: bind.Id,
			Env:    bind.Env,
		}, nil
	})
}

// BindAppAndServiceSource 绑定应用服务和服务来源
func BindAppAndServiceSource(ctx context.Context, reqDTO BindAppAndServiceSourceReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkManageServiceSourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	if len(reqDTO.SourceIdList) == 0 {
		err = deploymd.DeleteAppServiceSourceBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err != nil {
			return util.InternalError(err)
		}
		return nil
	}
	// 校验sourceIdList
	sources, err := deploymd.BatchGetServiceSourceByIdList(ctx, reqDTO.SourceIdList, nil)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if len(sources) == 0 {
		return util.InvalidArgsError()
	}
	for _, source := range sources {
		if source.Env != reqDTO.Env {
			return util.InvalidArgsError()
		}
	}
	insertReqs, _ := listutil.Map(reqDTO.SourceIdList, func(t int64) (deploymd.InsertAppServiceSourceBindReqDTO, error) {
		return deploymd.InsertAppServiceSourceBindReqDTO{
			SourceId: t,
			AppId:    reqDTO.AppId,
			Env:      reqDTO.Env,
		}, nil
	})
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 先删除
		err2 := deploymd.DeleteAppServiceSourceBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err2 != nil {
			return err2
		}
		// 批量插入
		return deploymd.BatchInsertAppServiceSourceBind(ctx, insertReqs)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func listActions(ctx context.Context, source deploymd.ServiceSource) ([]status.Action, error) {
	url := strings.TrimSuffix(source.Host, "/") + "/api/service/v1/status/actions"
	resp := make([]status.Action, 0)
	err := httputil.Get(
		ctx,
		http.DefaultClient,
		url,
		map[string]string{
			"Authorization": source.ApiKey,
		},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
