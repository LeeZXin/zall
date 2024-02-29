package webhookapi

import (
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InsertWebhookReqVO struct {
	RepoId   int64              `json:"repoId"`
	HookUrl  string             `json:"hookUrl"`
	Headers  map[string]string  `json:"headers"`
	HookType webhookmd.HookType `json:"hookType"`
}

type DeleteWebhookReqVO struct {
	Id int64 `json:"id"`
}

type ListWebhookReqVO struct {
	RepoId   int64              `json:"repoId"`
	HookType webhookmd.HookType `json:"hookType"`
}

type WebhookVO struct {
	Id          int64             `json:"id"`
	RepoId      int64             `json:"repoId"`
	HookUrl     string            `json:"hookUrl"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	HookType    string            `json:"hookType"`
}

type ListWebhookRespVO struct {
	ginutil.BaseResp
	Data []WebhookVO
}
