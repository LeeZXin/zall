package deploysrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/fileserv/modules/model/artifactmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/zalletmd"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/pkg/status"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppSourceTopic, func(data any) {
			req, ok := data.(event.AppSourceEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppManageServiceSourceAction:
							return cfg.AppSource.ManageServiceSource
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppDeployServiceTopic, func(data any) {
			req, ok := data.(event.AppDeployServiceEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Source.Env]
					if ok {
						switch req.Action {
						case event.AppDeployServiceTriggerActionAction:
							return cfg.AppDeployService.TriggerAction
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppDeployPipelineTopic, func(data any) {
			req, ok := data.(event.AppDeployPipelineEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppDeployPipelineCreatePipelineAction:
							return cfg.AppDeployPipeline.Create
						case event.AppDeployPipelineUpdatePipelineAction:
							return cfg.AppDeployPipeline.Update
						case event.AppDeployPipelineDeletePipelineAction:
							return cfg.AppDeployPipeline.Delete
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppDeployPipelineVarsTopic, func(data any) {
			req, ok := data.(event.AppDeployPipelineVarsEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppDeployPipelineVarsCreateAction:
							return cfg.AppDeployPipelineVars.Create
						case event.AppDeployPipelineVarsUpdateAction:
							return cfg.AppDeployPipelineVars.Update
						case event.AppDeployPipelineVarsDeleteAction:
							return cfg.AppDeployPipelineVars.Delete
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.AppDeployPlanTopic, func(data any) {
			req, ok := data.(event.AppDeployPlanEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppDeployPlanCreateAction:
							return cfg.AppDeployPlan.Create
						case event.AppDeployPlanCloseAction:
							return cfg.AppDeployPlan.Close
						case event.AppDeployPlanStartAction:
							return cfg.AppDeployPlan.Start
						default:
							return false
						}
					}
					return false
				})
			}
		})
	})
}

// CreatePlan 创建发布计划
func CreatePlan(ctx context.Context, reqDTO CreatePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	pipeline, app, team, err := checkDeployPlanPermByPipelineId(ctx, reqDTO.Operator, reqDTO.PipelineId)
	if err != nil {
		return err
	}
	// 检查制品号
	_, b, err := artifactmd.GetArtifactByAppIdAndNameAndEnv(ctx, artifactmd.GetArtifactReqDTO{
		AppId: app.AppId,
		Name:  reqDTO.ArtifactVersion,
		Env:   pipeline.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 检查是否有其他发布计划在
	b, err = deploymd.ExistPendingOrRunningPlanByPipelineId(ctx, reqDTO.PipelineId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	plan, err := deploymd.InsertPlan(ctx, deploymd.InsertPlanReqDTO{
		Name:            reqDTO.Name,
		PlanStatus:      deploymd.PendingPlanStatus,
		AppId:           pipeline.AppId,
		PipelineId:      reqDTO.PipelineId,
		PipelineName:    pipeline.Name,
		ArtifactVersion: reqDTO.ArtifactVersion,
		Creator:         reqDTO.Operator.Account,
		Env:             pipeline.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPlanEvent(
		reqDTO.Operator,
		team,
		app,
		plan,
		pipeline,
		event.AppDeployPlanCreateAction,
	)
	return nil
}

// StartPlan 开始发布计划
func StartPlan(ctx context.Context, reqDTO StartPlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, pipeline, app, team, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if plan.PlanStatus != deploymd.PendingPlanStatus {
		return util.InvalidArgsError()
	}
	if err != nil {
		return err
	}
	var plCfg deploy.Pipeline
	err = yaml.Unmarshal([]byte(pipeline.Config), &plCfg)
	if err != nil || !plCfg.IsValid() {
		return util.ThereHasBugErr()
	}
	nodes, err := zalletmd.ListAllZalletNode(ctx, []string{"node_id", "agent_host", "agent_token"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	agentMap := make(map[string]zalletmd.ZalletNode)
	for _, node := range nodes {
		agentMap[node.NodeId] = node
	}
	// 转化为insertReq
	stageList, err := pipelineCfg2DeployStageReq(&plan, &plCfg, agentMap)
	if err != nil {
		return util.OperationFailedError()
	}
	// 获取变量
	varsMap, err := deploymd.ListPipelineVarsMap(ctx, pipeline.AppId, pipeline.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	defaultEnv := util.MergeMap(varsMap, map[string]string{
		deploy.CurrentArtifactVersionKey: plan.ArtifactVersion,
		deploy.OperatorAccountKey:        reqDTO.Operator.Account,
		deploy.AppKey:                    plan.AppId,
	})
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b2, err2 := deploymd.StartPlan(ctx, reqDTO.PlanId, pipeline.Config)
		if err2 != nil {
			return err2
		}
		if !b2 {
			return fmt.Errorf("start plan failed: %v", reqDTO.PlanId)
		}
		// 新增部署步骤
		stages, err2 := deploymd.BatchInsertDeployStage(ctx, stageList...)
		if err2 != nil {
			return err2
		}
		return executeDeployOnStartPlan(
			plan,
			plCfg,
			defaultEnv,
			getStageMap(stages),
		)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.OperationFailedError()
	}
	notifyPlanEvent(
		reqDTO.Operator,
		team,
		app,
		plan,
		pipeline,
		event.AppDeployPlanStartAction,
	)
	return nil
}

func getStageMap(stages []deploymd.Stage) map[int][]deploymd.Stage {
	ret := make(map[int][]deploymd.Stage)
	for _, stage := range stages {
		dss, b := ret[stage.StageIndex]
		if !b {
			dss = make([]deploymd.Stage, 0)
		}
		ret[stage.StageIndex] = append(dss, stage)
	}
	return ret
}

// RedoAgentStage 重新执行agent
func RedoAgentStage(ctx context.Context, reqDTO RedoAgentStageReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	stage, plan, pipeline, _, _, err := checkAppDevelopPermByStageId(ctx, reqDTO.Operator, reqDTO.StageId)
	if err != nil {
		return err
	}
	// 校验状态
	if plan.PlanStatus != deploymd.RunningPlanStatus ||
		stage.StageStatus != deploymd.FailedStageStatus {
		return util.InvalidArgsError()
	}
	var dp deploy.Pipeline
	err = yaml.Unmarshal([]byte(plan.PipelineConfig), &dp)
	if err != nil {
		return nil
	}
	// 获取变量
	varsMap, err := deploymd.ListPipelineVarsMap(ctx, pipeline.AppId, pipeline.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 重置状态
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
		args[deploy.OperatorAccountKey] = reqDTO.Operator.Account
		err = redoAgentStageInRunner(dp, stage, util.MergeMap(varsMap, args))
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
	plan, pipeline, _, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
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
	for _, stage := range stages {
		// 前面阶段都必须成功
		if stage.StageIndex < reqDTO.StageIndex && stage.StageStatus != deploymd.SuccessfulStageStatus {
			return util.InvalidArgsError()
		}
		// 计算未成功数量
		if stage.StageIndex == reqDTO.StageIndex && stage.StageStatus != deploymd.SuccessfulStageStatus {
			notSuccessfulStages = append(notSuccessfulStages, stage)
		}
	}
	if len(notSuccessfulStages) == 0 {
		return util.InvalidArgsError()
	}
	args, err := deploymd.MergeInputArgsByPlanIdAndLETIndex(ctx, reqDTO.PlanId, reqDTO.StageIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	args[deploy.CurrentArtifactVersionKey] = plan.ArtifactVersion
	args[deploy.OperatorAccountKey] = reqDTO.Operator.Account
	args[deploy.AppKey] = plan.AppId
	err = redoAgentStagesInRunner(dp, notSuccessfulStages, util.MergeMap(varsMap, args))
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

func updateStageStatusWithOldStatusById(stageId int64, newStatus, oldStatus deploymd.StageStatus) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := deploymd.UpdateStageStatusWithOldStatusById(ctx, stageId, newStatus, oldStatus)
	if err != nil {
		logger.Logger.Error(err)
	}
	return b
}

func pipelineCfg2DeployStageReq(plan *deploymd.Plan, dp *deploy.Pipeline, agentMap map[string]zalletmd.ZalletNode) ([]deploymd.InsertDeployStageReqDTO, error) {
	ret := make([]deploymd.InsertDeployStageReqDTO, 0)
	for index, stage := range dp.Deploy {
		agents := make([]string, 0)
		// 如果配置没有指定agent 则认为全部agent
		if len(stage.Agents) > 0 {
			for _, agent := range stage.Agents {
				agents = append(agents, agent)
			}
		} else {
			for agent := range dp.Agents {
				agents = append(agents, agent)
			}
		}
		for _, agent := range agents {
			node, b := agentMap[dp.Agents[agent].NodeId]
			if !b {
				return nil, fmt.Errorf("%v not found", dp.Agents[agent].NodeId)
			}
			ret = append(ret, deploymd.InsertDeployStageReqDTO{
				PlanId:      plan.Id,
				AppId:       plan.AppId,
				StageIndex:  index,
				Agent:       agent,
				TaskId:      idutil.RandomUuid(),
				AgentHost:   node.AgentHost,
				AgentToken:  node.AgentToken,
				Script:      dp.Actions[stage.Action].Script,
				StageStatus: deploymd.PendingStageStatus,
			})
		}
	}
	return ret, nil
}

// ClosePlan 关闭发布计划
func ClosePlan(ctx context.Context, reqDTO ClosePlanReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	plan, pipeline, app, team, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return err
	}
	if plan.PlanStatus.IsFinalStatus() {
		return util.InvalidArgsError()
	}
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		_, err = deploymd.ClosePlanAndUpdateConfig(ctx, reqDTO.PlanId, plan.PlanStatus, pipeline.Config)
	} else {
		// 判断是否有执行中的任务
		b, err := deploymd.ExistRunningStatusByPlanId(ctx, reqDTO.PlanId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if b {
			return util.OperationFailedError()
		}
		_, err = deploymd.UpdatePlanStatusWithOldStatus(ctx, reqDTO.PlanId, deploymd.ClosedPlanStatus, plan.PlanStatus)
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPlanEvent(
		reqDTO.Operator,
		team,
		app,
		plan,
		pipeline,
		event.AppDeployPlanCloseAction,
	)
	return nil
}

// ListPlan 发布计划列表
func ListPlan(ctx context.Context, reqDTO ListPlanReqDTO) ([]PlanDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
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
	accounts := listutil.MapNe(plans, func(t deploymd.Plan) string {
		return t.Creator
	})
	userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, accounts...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data := listutil.MapNe(plans, func(t deploymd.Plan) PlanDTO {
		return PlanDTO{
			Id:              t.Id,
			PipelineId:      t.PipelineId,
			PipelineName:    t.PipelineName,
			Name:            t.Name,
			ArtifactVersion: t.ArtifactVersion,
			PlanStatus:      t.PlanStatus,
			Env:             t.Env,
			Creator:         userMap[t.Creator],
			Created:         t.Created,
		}
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
	plan, pipeline, _, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return PlanDetailDTO{}, err
	}
	userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, plan.Creator)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return PlanDetailDTO{}, util.InternalError(err)
	}
	ret := PlanDetailDTO{
		Id:              plan.Id,
		PipelineId:      plan.PipelineId,
		PipelineName:    pipeline.Name,
		Name:            plan.Name,
		ArtifactVersion: plan.ArtifactVersion,
		PlanStatus:      plan.PlanStatus,
		Env:             plan.Env,
		Creator:         userMap[plan.Creator],
		Created:         plan.Created,
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
	app, team, err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return err
	}
	pipeline, err := deploymd.InsertPipeline(ctx, deploymd.InsertPipelineReqDTO{
		AppId:  reqDTO.AppId,
		Config: reqDTO.Config,
		Env:    reqDTO.Env,
		Name:   reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPipelineEvent(
		reqDTO.Operator,
		team,
		app,
		pipeline,
		event.AppDeployPipelineCreatePipelineAction,
	)
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
	pipeline, app, team, err := checkPipelinePermByPipelineId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = deploymd.UpdatePipeline(ctx, deploymd.UpdatePipelineReqDTO{
		Id:     reqDTO.Id,
		Name:   reqDTO.Name,
		Config: reqDTO.Config,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPipelineEvent(
		reqDTO.Operator,
		team,
		app,
		pipeline,
		event.AppDeployPipelineUpdatePipelineAction,
	)
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
	pipeline, app, team, err := checkPipelinePermByPipelineId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 存在正在执行的发布计划
	b, err := deploymd.ExistPendingOrRunningPlanByPipelineId(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	_, err = deploymd.DeletePipelineById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPipelineEvent(
		reqDTO.Operator,
		team,
		app,
		pipeline,
		event.AppDeployPipelineDeletePipelineAction,
	)
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
	_, _, err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
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
	data := listutil.MapNe(pipelines, func(t deploymd.Pipeline) PipelineDTO {
		return PipelineDTO{
			Id:     t.Id,
			AppId:  t.AppId,
			Config: t.Config,
			Env:    t.Env,
			Name:   t.Name,
		}
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
	_, _, err := checkDeployPlanPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
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
	data := listutil.MapNe(pipelines, func(t deploymd.Pipeline) SimplePipelineDTO {
		return SimplePipelineDTO{
			Id:   t.Id,
			Env:  t.Env,
			Name: t.Name,
		}
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
	plan, srv, _, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
	if err != nil {
		return nil, err
	}
	var pipelineConfig string
	if plan.PlanStatus == deploymd.PendingPlanStatus {
		pipelineConfig = srv.Config
	} else {
		pipelineConfig = plan.PipelineConfig
	}
	var pl deploy.Pipeline
	err = yaml.Unmarshal([]byte(pipelineConfig), &pl)
	if err != nil || !pl.IsValid() {
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
	var agentMap map[string]zalletmd.ZalletNode
	for index, stage := range pl.Deploy {
		var (
			total, done, percent, pending, running float64

			hasError bool
		)
		total = float64(len(stage.Agents))
		if total == 0 {
			total = float64(len(pl.Agents))
		}
		subStages := make([]SubStageDTO, 0, len(stagesMap[index]))
		if len(stagesMap) > 0 {
			for _, md := range stagesMap[index] {
				sa := SubStageDTO{
					Id:         md.Id,
					Agent:      md.Agent,
					AgentHost:  md.AgentHost,
					ExecuteLog: md.ExecuteLog,
				}
				if plan.PlanStatus != deploymd.ClosedPlanStatus || md.StageStatus != deploymd.PendingStageStatus {
					sa.StageStatus = md.StageStatus
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
			if agentMap == nil {
				agentMap, err = getZalletMap()
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return nil, util.InternalError(err)
				}
			}
			if len(stage.Agents) > 0 {
				for _, agent := range stage.Agents {
					sa := SubStageDTO{
						Agent:     agent,
						AgentHost: agentMap[pl.Agents[agent].NodeId].AgentHost,
					}
					subStages = append(subStages, sa)
				}
			} else {
				for id, agent := range pl.Agents {
					sa := SubStageDTO{
						Agent:     id,
						AgentHost: agentMap[agent.NodeId].AgentHost,
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
		if len(pl.Actions) > 0 {
			dto.Script = pl.Actions[stage.Action].Script
		}
		dto.Confirm = stage.Confirm
		ret = append(ret, dto)
	}
	return ret, nil
}

func getZalletMap() (map[string]zalletmd.ZalletNode, error) {
	ret := make(map[string]zalletmd.ZalletNode)
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	nodes, err := zalletmd.ListAllZalletNode(ctx, []string{"node_id", "agent_host", "agent_token"})
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		ret[node.NodeId] = node
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
	plan, _, _, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
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
			err = sshagent.NewServiceCommand(stage.AgentHost, stage.AgentToken, plan.AppId).
				Kill(stage.TaskId)
			if err != nil {
				logger.Logger.Errorf("kill taskId: %s with err: %v", stage.TaskId, err)
			}
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
	plan, pipeline, _, _, err := checkAppDevelopPermByPlanId(ctx, reqDTO.Operator, reqDTO.PlanId)
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
	filteredArgs[deploy.CurrentArtifactVersionKey] = plan.ArtifactVersion
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
	err = executeDeployOnConfirmStage(
		plan,
		dp,
		util.MergeMap(filteredArgs, varsMap),
		reqDTO.StageIndex)
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
	data := listutil.MapNe(sources, func(t deploymd.ServiceSource) ServiceSourceDTO {
		return ServiceSourceDTO{
			Id:      t.Id,
			Name:    t.Name,
			Env:     t.Env,
			Host:    t.Host,
			ApiKey:  t.ApiKey,
			Created: t.Created,
		}
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
	data := listutil.MapNe(sources, func(t deploymd.ServiceSource) SimpleServiceSourceDTO {
		return SimpleServiceSourceDTO{
			Id:   t.Id,
			Name: t.Name,
		}
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
	_, _, err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
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
	app, team, err := checkManagePipelinePermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
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
	vars, err := deploymd.InsertPipelineVars(ctx, deploymd.InsertPipelineVarsReqDTO{
		AppId:   reqDTO.AppId,
		Env:     reqDTO.Env,
		Name:    reqDTO.Name,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPipelineVarsEvent(
		reqDTO.Operator,
		team,
		app,
		vars,
		event.AppDeployPipelineVarsCreateAction,
	)
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
	vars, app, team, err := checkManagePipelinePermByVarsId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = deploymd.UpdatePipelineVars(ctx, deploymd.UpdatePipelineVarsReqDTO{
		Id:      reqDTO.Id,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPipelineVarsEvent(
		reqDTO.Operator,
		team,
		app,
		vars,
		event.AppDeployPipelineVarsUpdateAction,
	)
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
	vars, app, team, err := checkManagePipelinePermByVarsId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = deploymd.DeletePipelineVarsById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyPipelineVarsEvent(
		reqDTO.Operator,
		team,
		app,
		vars,
		event.AppDeployPipelineVarsDeleteAction,
	)
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
	vars, _, _, err := checkManagePipelinePermByVarsId(ctx, reqDTO.Id, reqDTO.Operator)
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
	app, _, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
	if err != nil {
		return nil, err
	}
	url := strings.TrimSuffix(source.Host, "/") +
		fmt.Sprintf("/api/service/v1/status/list?app=%s&env=%s", app.AppId, source.Env)
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
	app, team, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
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
			notifyDeployServiceEvent(
				reqDTO.Operator,
				team,
				app,
				source,
				event.AppDeployServiceTriggerActionAction,
				reqDTO.Action,
			)
			return nil
		}
	}
	return nil
}

// ListStatusActions 获取服务操作列表
func ListStatusActions(ctx context.Context, reqDTO ListStatusActionReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, source, err := checkAppDevelopPermByBindId(ctx, reqDTO.Operator, reqDTO.BindId)
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
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
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
	sourceIdList := listutil.MapNe(binds, func(t deploymd.AppServiceSourceBind) int64 {
		bindMap[t.SourceId] = t
		return t.SourceId
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
	app, team, err := checkManageServiceSourcePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return err
	}
	if len(reqDTO.SourceIdList) == 0 {
		err = deploymd.DeleteAppServiceSourceBindByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
		if err != nil {
			return util.InternalError(err)
		}
		notifyServiceSourceEvent(
			reqDTO.Operator,
			team,
			app,
			nil,
			reqDTO.Env,
		)
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
	insertReqs := listutil.MapNe(reqDTO.SourceIdList, func(t int64) deploymd.InsertAppServiceSourceBindReqDTO {
		return deploymd.InsertAppServiceSourceBindReqDTO{
			SourceId: t,
			AppId:    reqDTO.AppId,
			Env:      reqDTO.Env,
		}
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
	notifyServiceSourceEvent(
		reqDTO.Operator,
		team,
		app,
		sources,
		reqDTO.Env,
	)
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

func notifyDeployServiceEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, source deploymd.ServiceSource, action event.AppDeployServiceEventAction, triggerAction string) {
	initPsub()
	psub.Publish(event.AppDeployServiceTopic, event.AppDeployServiceEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action:        action,
		TriggerAction: triggerAction,
		Source: event.AppSource{
			Id:   source.Id,
			Name: source.Name,
			Env:  source.Env,
		},
	})
}

func notifyServiceSourceEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, sources []deploymd.ServiceSource, env string) {
	initPsub()
	srcs := listutil.MapNe(sources, func(t deploymd.ServiceSource) event.AppSource {
		return event.AppSource{
			Id:   t.Id,
			Name: t.Name,
			Env:  t.Env,
		}
	})
	psub.Publish(event.AppSourceTopic, event.AppSourceEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, event.AppManageServiceSourceAction.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, event.AppManageServiceSourceAction.GetI18nValue()),
		},
		Env:     env,
		Action:  event.AppManageServiceSourceAction,
		Sources: srcs,
	})
}

func notifyPipelineEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, pipeline deploymd.Pipeline, action event.AppDeployPipelineEventAction) {
	initPsub()
	psub.Publish(event.AppDeployPipelineTopic, event.AppDeployPipelineEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		BasePipeline: event.BasePipeline{
			PipelineId:   pipeline.Id,
			PipelineName: pipeline.Name,
			Env:          pipeline.Env,
		},
		Action: action,
	})
}

func notifyPipelineVarsEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, vars deploymd.PipelineVars, action event.AppDeployPipelineVarsEventAction) {
	initPsub()
	psub.Publish(event.AppDeployPipelineVarsTopic, event.AppDeployPipelineVarsEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		VarsId:   vars.Id,
		VarsName: vars.Name,
		Env:      vars.Env,
		Action:   action,
	})
}

func notifyPlanEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, plan deploymd.Plan, pipeline deploymd.Pipeline, action event.AppDeployPlanEventAction) {
	initPsub()
	psub.Publish(event.AppDeployPlanTopic, event.AppDeployPlanEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		BasePipeline: event.BasePipeline{
			PipelineId:   pipeline.Id,
			PipelineName: pipeline.Name,
			Env:          pipeline.Env,
		},
		Action:   action,
		PlanId:   plan.Id,
		PlanName: plan.Name,
		Env:      plan.Env,
	})
}
