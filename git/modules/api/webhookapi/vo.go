package webhookapi

import (
	"github.com/LeeZXin/zall/pkg/webhook"
)

type CreateWebhookReqVO struct {
	RepoId  int64          `json:"repoId"`
	HookUrl string         `json:"hookUrl"`
	Secret  string         `json:"secret"`
	Events  webhook.Events `json:"events"`
}

type WebhookVO struct {
	Id      int64          `json:"id"`
	RepoId  int64          `json:"repoId"`
	HookUrl string         `json:"hookUrl"`
	Secret  string         `json:"secret"`
	Events  webhook.Events `json:"events"`
}

type UpdateWebhookReqVO struct {
	WebhookId int64          `json:"webhookId"`
	HookUrl   string         `json:"hookUrl"`
	Secret    string         `json:"secret"`
	Events    webhook.Events `json:"events"`
}
