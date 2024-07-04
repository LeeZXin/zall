package statussrv

import (
	"context"
	"github.com/LeeZXin/zall/pkg/status"
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
	// GetActions 获取基本信息
	GetActions(context.Context) []status.Action
	// ListService 服务列表
	ListService(context.Context, status.ListServiceReq) ([]status.Service, error)
	// KillService 关闭服务
	KillService(context.Context, string) error
	// RestartService 重启服务
	RestartService(context.Context, string) error
}
