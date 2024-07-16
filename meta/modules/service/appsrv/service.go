package appsrv

import "context"

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// CreateApp 创建应用服务
	CreateApp(context.Context, CreateAppReqDTO) error
	// DeleteApp 删除应用服务
	DeleteApp(context.Context, DeleteAppReqDTO) error
	// UpdateApp 编辑应用服务
	UpdateApp(context.Context, UpdateAppReqDTO) error
	// GetApp 获取服务信息
	GetApp(context.Context, GetAppReqDTO) (AppDTO, error)
	// ListApp 应用服务列表
	ListApp(context.Context, ListAppReqDTO) ([]AppDTO, error)
	// ListAllAppByAdmin 所有应用服务列表 管理员权限
	ListAllAppByAdmin(context.Context, ListAppReqDTO) ([]AppDTO, error)
	TransferTeam(context.Context, TransferTeamReqDTO) error
}
