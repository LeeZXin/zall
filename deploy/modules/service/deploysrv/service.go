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
	// ListConfig 获取部署配置
	ListConfig(context.Context, ListConfigReqDTO) ([]ConfigDTO, error)
	// UpdateConfig 编辑部署配置
	UpdateConfig(context.Context, UpdateConfigReqDTO) error
	// CreateConfig 新增部署配置
	CreateConfig(context.Context, CreateConfigReqDTO) error
	// DeleteConfig 删除配置
	DeleteConfig(context.Context, DeleteConfigReqDTO) error
	// CreatePlan 创建发布计划
	CreatePlan(context.Context, CreatePlanReqDTO) error
	// ClosePlan 关闭发布计划
	ClosePlan(context.Context, ClosePlanReqDTO) error
	// ListPlan 发布计划列表
	ListPlan(context.Context, ListPlanReqDTO) ([]PlanDTO, int64, error)
	// AddPlanService 添加发布计划部署服务
	AddPlanService(context.Context, AddPlanServiceReqDTO) error
	// DeletePendingPlanService 删除未执行发布计划单项服务
	DeletePendingPlanService(context.Context, DeletePendingPlanServiceReqDTO) error
	// ListPlanService 展示发布计划的服务
	ListPlanService(context.Context, ListPlanServiceReqDTO) ([]PlanServiceDTO, error)
	// StartPlanService 启动部署服务流水线
	StartPlanService(context.Context, StartPlanServiceReqDTO) error
	// FinishPlanService 完成服务流水线
	FinishPlanService(context.Context, FinishPlanServiceReqDTO) error
	// ConfirmPlanServiceStep 执行流水线其中一个环节
	ConfirmPlanServiceStep(context.Context, ConfirmPlanServiceStepReqDTO) error
	// RollbackPlanServiceStep 回滚执行流水线其中一个环节
	RollbackPlanServiceStep(context.Context, RollbackPlanServiceStepReqDTO) error
}
