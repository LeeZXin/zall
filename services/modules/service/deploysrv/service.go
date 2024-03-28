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
}
