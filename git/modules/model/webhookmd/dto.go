package webhookmd

import "github.com/LeeZXin/zall/pkg/webhook"

type InsertWebhookReqDTO struct {
	RepoId  int64
	HookUrl string
	Secret  string
	Events  webhook.Events
}

type UpdateWebhookReqDTO struct {
	Id      int64
	HookUrl string
	Secret  string
	Events  webhook.Events
}
