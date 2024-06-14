package lfssrv

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
	Lock(context.Context, LockReqDTO) (LockRespDTO, error)
	ListLock(context.Context, ListLockReqDTO) (ListLockRespDTO, error)
	Unlock(context.Context, UnlockReqDTO) (LfsLockDTO, error)
	Verify(context.Context, VerifyReqDTO) (bool, bool, error)
	Download(context.Context, DownloadReqDTO) error
	Upload(context.Context, UploadReqDTO) error
	Batch(context.Context, BatchReqDTO) (BatchRespDTO, error)
}
