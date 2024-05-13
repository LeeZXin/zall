package webhooksrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
)

type CreateWebhookReqDTO struct {
	RepoId   int64               `json:"repoId"`
	HookUrl  string              `json:"hookUrl"`
	Secret   string              `json:"secret"`
	Events   webhook.Events      `json:"events"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateWebhookReqDTO) IsValid() error {
	parsedUrl, err := url.Parse(r.HookUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Secret) == 0 || len(r.Secret) > 1024 {
		return util.InvalidArgsError()
	}
	if !r.Events.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteWebhookReqDTO struct {
	WebhookId int64               `json:"webhookId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DeleteWebhookReqDTO) IsValid() error {
	if r.WebhookId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PingWebhookReqDTO struct {
	WebhookId int64               `json:"webhookId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *PingWebhookReqDTO) IsValid() error {
	if r.WebhookId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListWebhookReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListWebhookReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type WebhookDTO struct {
	Id      int64
	RepoId  int64
	HookUrl string
	Secret  string
	Events  webhook.Events
}

type UpdateWebhookReqDTO struct {
	WebhookId int64               `json:"webhookId"`
	HookUrl   string              `json:"hookUrl"`
	Secret    string              `json:"secret"`
	Events    webhook.Events      `json:"events"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateWebhookReqDTO) IsValid() error {
	if r.WebhookId <= 0 {
		return util.InvalidArgsError()
	}
	parsedUrl, err := url.Parse(r.HookUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Secret) > 1024 {
		return util.InvalidArgsError()
	}
	if !r.Events.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
