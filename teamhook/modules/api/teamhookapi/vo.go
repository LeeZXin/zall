package teamhookapi

import "github.com/LeeZXin/zall/pkg/teamhook"

type CreateTeamHookReqVO struct {
	Name     string            `json:"name"`
	TeamId   int64             `json:"teamId"`
	Events   teamhook.Events   `json:"events"`
	HookType teamhook.HookType `json:"hookType"`
	HookCfg  teamhook.Cfg      `json:"hookCfg"`
}

type UpdateTeamHookReqVO struct {
	Id       int64             `json:"id"`
	Name     string            `json:"name"`
	Events   teamhook.Events   `json:"events"`
	HookType teamhook.HookType `json:"hookType"`
	HookCfg  teamhook.Cfg      `json:"hookCfg"`
}

type TeamHookVO struct {
	Id       int64             `json:"id"`
	Name     string            `json:"name"`
	TeamId   int64             `json:"teamId"`
	Events   teamhook.Events   `json:"events"`
	HookType teamhook.HookType `json:"hookType"`
	HookCfg  teamhook.Cfg      `json:"hookCfg"`
}
