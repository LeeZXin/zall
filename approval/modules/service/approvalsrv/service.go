package approvalsrv

import "context"

var (
	Inner InnerService = new(innerImpl)
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	InsertFlow(context.Context, InsertFlowReqDTO) error
	InsertProcess(context.Context, InsertProcessReqDTO) error
	UpdateProcess(context.Context, UpdateProcessReqDTO) error
}

type OuterService interface {
	Agree(context.Context, AgreeFlowReqDTO) error
	Disagree(context.Context, DisagreeFlowReqDTO) error
	//InsertProcess(context.Context, InsertProcessReqDTO) error
	//UpdateProcess(context.Context, UpdateProcessReqDTO) error
	//DeleteProcess(context.Context, DeleteProcessReqDTO) error
	//GetProcess(context.Context, GetProcessReqDTO) (approval.Process, bool, error)
}
