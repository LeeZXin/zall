package webhooksrv

import (
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"net/url"
)

type InsertWebhookReqDTO struct {
	RepoId      int64               `json:"repoId"`
	HookUrl     string              `json:"hookUrl"`
	HttpHeaders map[string]string   `json:"httpHeaders"`
	HookType    webhookmd.HookType  `json:"hookType"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *InsertWebhookReqDTO) IsValid() error {
	_, err := url.Parse(r.HookUrl)
	if err != nil {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.HookType.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteWebhookReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteWebhookReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListWebhookReqDTO struct {
	RepoId   int64               `json:"repoId"`
	HookType webhookmd.HookType  `json:"hookType"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListWebhookReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.HookType.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type WebhookDTO struct {
	Id          int64
	RepoId      int64
	HookUrl     string
	HttpHeaders map[string]string
	HookType    webhookmd.HookType
}
