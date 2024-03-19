package webhooksrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertWebhook(context.Context, InsertWebhookReqDTO) error
	UpdateWebhook(context.Context, UpdateWebhookReqDTO) error
	DeleteWebhook(context.Context, DeleteWebhookReqDTO) error
	ListWebhook(context.Context, ListWebhookReqDTO) ([]WebhookDTO, error)
}
