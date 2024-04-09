package alertsrv

import "context"

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	// InsertConfig 新增配置
	InsertConfig(context.Context, InsertConfigReqDTO) error
	// UpdateConfig 修改配置
	UpdateConfig(context.Context, UpdateConfigReqDTO) error
	// DeleteConfig 删除配置
	DeleteConfig(context.Context, DeleteConfigReqDTO) error
	// ListConfig 展示配置
	ListConfig(context.Context, ListConfigReqDTO) ([]ConfigDTO, int64, error)
}
