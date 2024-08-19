package notifysrv

import (
	"github.com/LeeZXin/zall/notify/modules/model/notifymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/notify/notify"
	"github.com/LeeZXin/zall/util"
)

type CreateTplReqDTO struct {
	Name     string              `json:"name"`
	TeamId   int64               `json:"teamId"`
	Cfg      notify.Cfg          `json:"cfg"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateTplReqDTO) IsValid() error {
	if !notifymd.IsTplNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Cfg.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTplReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	Cfg      notify.Cfg          `json:"cfg"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTplReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !notifymd.IsTplNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Cfg.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTplReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTplReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTplReqDTO struct {
	Name     string              `json:"name"`
	TeamId   int64               `json:"teamId"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTplReqDTO) IsValid() error {
	if r.Name != "" && !notifymd.IsTplNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type TplDTO struct {
	Id        int64
	Name      string
	ApiKey    string
	TeamId    int64
	NotifyCfg notify.Cfg
}

type SimpleTplDTO struct {
	Id   int64
	Name string
}

type ChangeTplApiKeyReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ChangeTplApiKeyReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SendNotifyByApiKeyReqDTO struct {
	ApiKey string            `json:"apiKey"`
	Params map[string]string `json:"params"`
}

func (r *SendNotifyByApiKeyReqDTO) IsValid() error {
	if len(r.ApiKey) != 32 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAllTplReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllTplReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
