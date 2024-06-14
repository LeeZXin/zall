package appsrv

import "context"

var (
	Outer OuterService
)

func Init() {
	Outer = new(outerImpl)
}

type OuterService interface {
	// CreateApp 创建应用服务
	CreateApp(context.Context, CreateAppReqDTO) error
	DeleteApp(context.Context, DeleteAppReqDTO) error
	UpdateApp(context.Context, UpdateAppReqDTO) error
	// ListApp 应用服务列表
	ListApp(context.Context, ListAppReqDTO) ([]AppDTO, error)
	// ListAllAppByAdmin 所有应用服务列表 管理员权限
	ListAllAppByAdmin(context.Context, ListAppReqDTO) ([]AppDTO, error)
	TransferTeam(context.Context, TransferTeamReqDTO) error
}
