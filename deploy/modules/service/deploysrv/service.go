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
	// CreateService 创建服务
	CreateService(context.Context, CreateServiceReqDTO) error
	// UpdateService 编辑服务
	UpdateService(context.Context, UpdateServiceReqDTO) error
	// DeleteService 删除服务
	DeleteService(context.Context, DeleteServiceReqDTO) error
	// ListService 服务列表
	ListService(context.Context, ListServiceReqDTO) ([]ServiceDTO, error)
	// ListServiceWhenCreatePlan 创建发布计划时展示的服务列表
	ListServiceWhenCreatePlan(context.Context, ListServiceWhenCreatePlanReqDTO) ([]SimpleServiceDTO, error)
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
}
