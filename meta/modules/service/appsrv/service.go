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
	// GetAppWithPerm 获取服务信息
	GetAppWithPerm(context.Context, GetAppWithPermReqDTO) (AppWithPermDTO, error)
	// ListApp 应用服务列表
	ListApp(context.Context, ListAppReqDTO) ([]AppDTO, error)
	// ListAllAppByAdmin 所有应用服务列表 管理员权限
	ListAllAppByAdmin(context.Context, ListAllAppByAdminReqDTO) ([]AppDTO, error)
	// ListAllAppBySa 所有应用服务列表 超级管理员权限
	ListAllAppBySa(context.Context, ListAllAppBySaReqDTO) ([]AppDTO, error)
	// TransferTeam 迁移团队
	TransferTeam(context.Context, TransferTeamReqDTO) error
}
