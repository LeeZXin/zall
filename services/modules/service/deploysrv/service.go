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
}
