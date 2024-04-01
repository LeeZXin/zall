package deploysrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	// DeployServiceWithoutPlan 不通过发布计划部署服务
	DeployServiceWithoutPlan(context.Context, DeployServiceWithoutPlanReqDTO) error
}

type OuterService interface {
	// ListConfig 获取部署配置
	ListConfig(context.Context, ListConfigReqDTO) ([]ConfigDTO, error)
	// UpdateConfig 编辑部署配置
	UpdateConfig(context.Context, UpdateConfigReqDTO) error
	// InsertConfig 新增部署配置
	InsertConfig(context.Context, InsertConfigReqDTO) error
	// InsertPlan 创建发布计划
	InsertPlan(context.Context, InsertPlanReqDTO) error
	// ClosePlan 关闭发布计划
	ClosePlan(context.Context, ClosePlanReqDTO) error
	// ListPlan 发布计划列表
	ListPlan(context.Context, ListPlanReqDTO) ([]PlanDTO, int64, error)
	// InsertPlanItem 添加发布计划部署服务
	InsertPlanItem(context.Context, InsertPlanItemReqDTO) error
	// ClosePlanItem 关闭发布计划单项服务
	ClosePlanItem(context.Context, ClosePlanItemReqDTO) error
	// ListPlanItem 展示发布计划的服务
	ListPlanItem(context.Context, ListPlanItemReqDTO) ([]PlanItemDTO, error)
	// DeployService 重建服务
	DeployService(context.Context, DeployServiceReqDTO) error
	// StopService 下线服务
	StopService(context.Context, StopServiceReqDTO) error
	// RestartService 重启服务
	RestartService(context.Context, RestartServiceReqDTO) error
	// ListService 服务列表
	ListService(context.Context, ListServiceReqDTO) ([]ServiceDTO, error)
	// ListDeployLog 查看部署日志
	ListDeployLog(context.Context, ListDeployLogReqDTO) ([]DeployLogDTO, int64, error)
	// ListOpLog 查看操作日志
	ListOpLog(context.Context, ListOpLogReqDTO) ([]OpLogDTO, int64, error)
	// DeployServiceWithPlan 通过发布计划部署服务
	DeployServiceWithPlan(context.Context, DeployServiceWithPlanReqDTO) error
	// RollbackServiceWithPlan 通过发布计划回滚服务
	RollbackServiceWithPlan(context.Context, RollbackServiceWithPlanReqDTO) error
}
