package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/pkg/status"
)

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = NewOuterService()
	}
}

type OuterService interface {
	// CreatePlan 创建发布计划
	CreatePlan(context.Context, CreatePlanReqDTO) error
	// StartPlan 开始发布计划
	StartPlan(context.Context, StartPlanReqDTO) error
	// ClosePlan 关闭发布计划
	ClosePlan(context.Context, ClosePlanReqDTO) error
	// ListPlan 发布计划列表
	ListPlan(context.Context, ListPlanReqDTO) ([]PlanDTO, int64, error)
	// GetPlanDetail 获取单个发布计划详情
	GetPlanDetail(context.Context, GetPlanDetailReqDTO) (PlanDetailDTO, error)
	// CreatePipeline 创建流水线
	CreatePipeline(context.Context, CreatePipelineReqDTO) error
	// UpdatePipeline 编辑流水线
	UpdatePipeline(context.Context, UpdatePipelineReqDTO) error
	// DeletePipeline 删除流水线
	DeletePipeline(context.Context, DeletePipelineReqDTO) error
	// ListPipeline 流水线列表
	ListPipeline(context.Context, ListPipelineReqDTO) ([]PipelineDTO, error)
	// ListPipelineWhenCreatePlan 创建发布计划时展示的流水线列表
	ListPipelineWhenCreatePlan(context.Context, ListPipelineWhenCreatePlanReqDTO) ([]SimplePipelineDTO, error)
	// ListPipelineVars 流水线变量
	ListPipelineVars(context.Context, ListPipelineVarsReqDTO) ([]PipelineVarsWithoutContentDTO, error)
	// CreatePipelineVars 创建流水线变量
	CreatePipelineVars(context.Context, CreatePipelineVarsReqDTO) error
	// UpdatePipelineVars 编辑流水线变量
	UpdatePipelineVars(context.Context, UpdatePipelineVarsReqDTO) error
	// DeletePipelineVars 删除流水线变量
	DeletePipelineVars(context.Context, DeletePipelineVarsReqDTO) error
	// GetPipelineVarsContent 获取流水线变量内容
	GetPipelineVarsContent(context.Context, GetPipelineVarsContentReqDTO) (PipelineVarsDTO, error)
	// ListStages 展示发布计划流水线详情
	ListStages(context.Context, ListStagesReqDTO) ([]StageDTO, error)
	// RedoAgentStage 重新执行agent
	RedoAgentStage(context.Context, RedoAgentStageReqDTO) error
	// ForceRedoNotSuccessfulAgentStages 强制重新执行未完成的任务
	ForceRedoNotSuccessfulAgentStages(context.Context, ForceRedoNotSuccessfulAgentStagesReqDTO) error
	// KillStage 中止执行
	KillStage(context.Context, KillStageReqDTO) error
	// ConfirmInteractStage 交互阶段确认
	ConfirmInteractStage(context.Context, ConfirmInteractStageReqDTO) error
	// ListServiceSource 查看服务状态来源
	ListServiceSource(context.Context, ListServiceSourceReqDTO) ([]ServiceSourceDTO, error)
	// ListAllServiceSource 所有服务状态来源
	ListAllServiceSource(context.Context, ListAllServiceSourceReqDTO) ([]SimpleServiceSourceDTO, error)
	// CreateServiceSource 插入服务状态来源
	CreateServiceSource(context.Context, CreateServiceSourceReqDTO) error
	// UpdateServiceSource 更新服务状态来源
	UpdateServiceSource(context.Context, UpdateServiceSourceReqDTO) error
	// DeleteServiceSource 删除服务状态来源
	DeleteServiceSource(context.Context, DeleteServiceSourceReqDTO) error
	// ListServiceStatus 展示服务状态列表
	ListServiceStatus(context.Context, ListServiceStatusReqDTO) ([]status.Service, error)
	// ListStatusActions 获取服务操作列表
	ListStatusActions(context.Context, ListStatusActionReqDTO) ([]string, error)
	// DoStatusAction 操作服务
	DoStatusAction(context.Context, DoStatusActionReqDTO) error
	// ListBindServiceSource 获取绑定服务来源
	ListBindServiceSource(context.Context, ListBindServiceSourceReqDTO) ([]SimpleBindServiceSourceDTO, error)
	// BindAppAndServiceSource 绑定应用服务和服务来源
	BindAppAndServiceSource(context.Context, BindAppAndServiceSourceReqDTO) error
}
