package alertsrv

import "context"

type OuterService interface {
	InsertConfig(context.Context, InsertConfigReqDTO) error
}
