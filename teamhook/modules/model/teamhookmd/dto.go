package teamhookmd

import (
	"github.com/LeeZXin/zall/pkg/commonhook"
	"github.com/LeeZXin/zall/pkg/teamhook"
)

type InsertTeamHookReqDTO struct {
	Name     string
	TeamId   int64
	Events   teamhook.Events
	HookType commonhook.HookType
	HookCfg  commonhook.Cfg
}

type UpdateTeamHookReqDTO struct {
	Id       int64
	Name     string
	Events   teamhook.Events
	HookType commonhook.HookType
	HookCfg  commonhook.Cfg
}
