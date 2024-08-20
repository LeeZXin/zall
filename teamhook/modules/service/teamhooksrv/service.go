package teamhooksrv

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
	// CreateTeamHook 创建team hook
	CreateTeamHook(context.Context, CreateTeamHookReqDTO) error
	// UpdateTeamHook 编辑team hook
	UpdateTeamHook(context.Context, UpdateTeamHookReqDTO) error
	// DeleteTeamHook 删除team hook
	DeleteTeamHook(context.Context, DeleteTeamHookReqDTO) error
	// ListTeamHook team hook 列表
	ListTeamHook(context.Context, ListTeamHookReqDTO) ([]TeamHookDTO, error)
}
