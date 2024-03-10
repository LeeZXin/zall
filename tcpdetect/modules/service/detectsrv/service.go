package detectsrv

import "context"

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertDetect(context.Context, InsertDetectReqDTO) error
	UpdateDetect(context.Context, UpdateDetectReqDTO) error
	ListDetect(context.Context, ListDetectReqDTO) ([]DetectDTO, int64, error)
	DeleteDetect(context.Context, DeleteDetectReqDTO) error
	ListLog(context.Context, ListLogReqDTO) ([]LogDTO, int64, error)
	EnabledDetect(context.Context, EnableDetectReqDTO) error
	DisableDetect(context.Context, DisableDetectReqDTO) error
}
