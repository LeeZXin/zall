package webhooksrv

import (
	"context"
)

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// CreateWebhook 新增webhook
	CreateWebhook(context.Context, CreateWebhookReqDTO) error
	// UpdateWebhook 编辑webhook
	UpdateWebhook(context.Context, UpdateWebhookReqDTO) error
	// DeleteWebhook 删除webhook
	DeleteWebhook(context.Context, DeleteWebhookReqDTO) error
	// ListWebhook 列表webhook
	ListWebhook(context.Context, ListWebhookReqDTO) ([]WebhookDTO, error)
	// PingWebhook ping webhook
	PingWebhook(context.Context, PingWebhookReqDTO) error
}
