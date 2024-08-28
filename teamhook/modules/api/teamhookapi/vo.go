package teamhookapi

import (
	"github.com/LeeZXin/zall/pkg/commonhook"
	"github.com/LeeZXin/zall/pkg/teamhook"
)

type CreateTeamHookReqVO struct {
	Name   string          `json:"name"`
	TeamId int64           `json:"teamId"`
	Events teamhook.Events `json:"events"`
	commonhook.TypeAndCfg
}

type UpdateTeamHookReqVO struct {
	Id     int64           `json:"id"`
	Name   string          `json:"name"`
	Events teamhook.Events `json:"events"`
	commonhook.TypeAndCfg
}

type TeamHookVO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	TeamId   int64               `json:"teamId"`
	Events   teamhook.Events     `json:"events"`
	HookType commonhook.HookType `json:"hookType"`
	HookCfg  commonhook.Cfg      `json:"hookCfg"`
}
