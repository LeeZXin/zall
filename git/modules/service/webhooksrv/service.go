package webhooksrv

import (
	"context"
)

var (
	Outer OuterService
)

func Init() {
	Outer = newOuterService()
}

type OuterService interface {
	// CreateWebhook 新增webhook
	CreateWebhook(context.Context, CreateWebhookReqDTO) error
	// UpdateWebhook 编辑
	UpdateWebhook(context.Context, UpdateWebhookReqDTO) error
	// DeleteWebhook 删除
	DeleteWebhook(context.Context, DeleteWebhookReqDTO) error
	// ListWebhook 列表
	ListWebhook(context.Context, ListWebhookReqDTO) ([]WebhookDTO, error)
	// PingWebhook ping
	PingWebhook(context.Context, PingWebhookReqDTO) error
}
