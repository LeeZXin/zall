package approvalsrv

import "context"

var (
	Inner InnerService = new(innerImpl)
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	InsertFlow(context.Context, string) error
	InsertProcess(context.Context, InsertProcessReqDTO) error
	UpdateProcess(context.Context, UpdateProcessReqDTO) error
}

type OuterService interface {
	Agree(context.Context, AgreeFlowReqDTO) error
	Disagree(context.Context, DisagreeFlowReqDTO) error
	//InsertProcess(context.Context, InsertProcessReqDTO) error
	//UpdateProcess(context.Context, UpdateProcessReqDTO) error
	//DeleteProcess(context.Context, DeleteProcessReqDTO) error
	//GetProcess(context.Context, GetProcessReqDTO) (approval.Approval, bool, error)
}
