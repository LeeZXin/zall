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
	// ConfirmPlanServiceStep 执行流水线其中一个环节
	ConfirmPlanServiceStep(context.Context, ConfirmPlanServiceStepReqDTO) error
	// RollbackPlanServiceStep 回滚执行流水线其中一个环节
	RollbackPlanServiceStep(context.Context, RollbackPlanServiceStepReqDTO) error
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
}
