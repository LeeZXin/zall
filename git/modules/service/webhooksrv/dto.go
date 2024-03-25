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
	WildBranch  string              `json:"wildBranch"`
	WildTag     string              `json:"wildTag"`
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
	switch r.HookType {
	case webhookmd.PushHook:
		if len(r.WildBranch) == 0 || len(r.WildBranch) > 32 {
			return util.InvalidArgsError()
		}
	case webhookmd.TagHook:
		if len(r.WildTag) == 0 || len(r.WildTag) > 32 {
			return util.InvalidArgsError()
		}
	case webhookmd.PullRequestHook:
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
	HttpHeaders webhookmd.HttpHeaders
	HookType    webhookmd.HookType
	WildBranch  string
	WildTag     string
}

type UpdateWebhookReqDTO struct {
	Id          int64               `json:"id"`
	HookUrl     string              `json:"hookUrl"`
	HttpHeaders map[string]string   `json:"httpHeaders"`
	HookType    webhookmd.HookType  `json:"hookType"`
	WildBranch  string              `json:"wildBranch"`
	WildTag     string              `json:"wildTag"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *UpdateWebhookReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
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
	switch r.HookType {
	case webhookmd.PushHook:
		if len(r.WildBranch) == 0 || len(r.WildBranch) > 32 {
			return util.InvalidArgsError()
		}
	case webhookmd.TagHook:
		if len(r.WildTag) == 0 || len(r.WildTag) > 32 {
			return util.InvalidArgsError()
		}
	case webhookmd.PullRequestHook:
	}
	return nil
}
