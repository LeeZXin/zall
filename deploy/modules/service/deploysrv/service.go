package deploysrv

import (
	"context"
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
	// CreatePipeline 创建服务
	CreatePipeline(context.Context, CreatePipelineReqDTO) error
	// UpdatePipeline 编辑服务
	UpdatePipeline(context.Context, UpdatePipelineReqDTO) error
	// DeletePipeline 删除服务
	DeletePipeline(context.Context, DeletePipelineReqDTO) error
	// ListPipeline 服务列表
	ListPipeline(context.Context, ListPipelineReqDTO) ([]PipelineDTO, error)
	// ListPipelineWhenCreatePlan 创建发布计划时展示的服务列表
	ListPipelineWhenCreatePlan(context.Context, ListPipelineWhenCreatePlanReqDTO) ([]SimplePipelineDTO, error)
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
	// ListServiceSource 查看服务数据来源
	ListServiceSource(context.Context, ListServiceSourceReqDTO) ([]ServiceSourceDTO, error)
	// CreateServiceSource 插入服务数据来源
	CreateServiceSource(context.Context, CreateServiceSourceReqDTO) error
	// UpdateServiceSource 更新数据服务来源
	UpdateServiceSource(context.Context, UpdateServiceSourceReqDTO) error
	// DeleteServiceSource 删除数据服务来源
	DeleteServiceSource(context.Context, DeleteServiceSourceReqDTO) error
}
