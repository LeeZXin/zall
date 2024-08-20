package teamhooksrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/model/teamhookmd"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
)

type CreateTeamHookReqDTO struct {
	Name     string              `json:"name"`
	TeamId   int64               `json:"teamId"`
	Events   teamhook.Events     `json:"events"`
	HookType teamhook.HookType   `json:"hookType"`
	HookCfg  teamhook.Cfg        `json:"hookCfg"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateTeamHookReqDTO) IsValid() error {
	if !teamhookmd.IsTeamHookNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.HookType.IsValid() {
		return util.InvalidArgsError()
	}
	switch r.HookType {
	case teamhook.NotifyType:
		if r.HookCfg.NotifyTplId <= 0 {
			return util.InvalidArgsError()
		}
	case teamhook.WebhookType:
		parsedUrl, err := url.Parse(r.HookCfg.HookUrl)
		if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
			return util.InvalidArgsError()
		}
		if r.HookCfg.Secret == "" {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTeamHookReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	Events   teamhook.Events     `json:"events"`
	HookType teamhook.HookType   `json:"hookType"`
	HookCfg  teamhook.Cfg        `json:"hookCfg"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTeamHookReqDTO) IsValid() error {
	if !teamhookmd.IsTeamHookNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.HookType.IsValid() {
		return util.InvalidArgsError()
	}
	switch r.HookType {
	case teamhook.NotifyType:
		if r.HookCfg.NotifyTplId <= 0 {
			return util.InvalidArgsError()
		}
	case teamhook.WebhookType:
		parsedUrl, err := url.Parse(r.HookCfg.HookUrl)
		if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTeamHookReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTeamHookReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTeamHookReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTeamHookReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TeamHookDTO struct {
	Id       int64
	Name     string
	TeamId   int64
	Events   teamhook.Events
	HookType teamhook.HookType
	HookCfg  teamhook.Cfg
}
