package webhooksrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertWebHook(context.Context, InsertWebhookReqDTO) error
	DeleteWebhook(context.Context, DeleteWebhookReqDTO) error
	ListWebhook(context.Context, ListWebhookReqDTO) ([]WebhookDTO, error)
}
