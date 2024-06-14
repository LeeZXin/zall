package smartsrv

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
	UploadPack(context.Context, UploadPackReqDTO) error
	ReceivePack(context.Context, ReceivePackReqDTO) error
	InfoRefs(context.Context, InfoRefsReqDTO) error
}
