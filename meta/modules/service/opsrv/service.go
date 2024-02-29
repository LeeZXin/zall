package opsrv

import (
	"context"
)

var (
	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	InsertOpLog(context.Context, InsertOpLogReqDTO)
}
