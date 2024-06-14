package servicesrv

import "context"

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		initTask()
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// CreateService 创建服务
	CreateService(context.Context, CreateServiceReqDTO) error
	// UpdateService 编辑服务
	UpdateService(context.Context, UpdateServiceReqDTO) error
	// DeleteService 删除服务
	DeleteService(context.Context, DeleteServiceReqDTO) error
	// ListService 服务列表
	ListService(context.Context, ListServiceReqDTO) ([]ServiceDTO, int64, error)
	// EnableService 启动服务
	EnableService(context.Context, EnableServiceReqDTO) error
	// DisableService 关闭服务
	DisableService(context.Context, DisableServiceReqDTO) error
}
