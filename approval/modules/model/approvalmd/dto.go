package approvalmd

import (
	"github.com/LeeZXin/zall/approval/approval"
	"time"
)

type InsertProcessReqDTO struct {
	Pid      string
	Name     string
	Approval approval.Approval
}

type UpdateProcessReqDTO struct {
	Id       int64
	Name     string
	Approval approval.Approval
}

type UpdateProcessByPidReqDTO struct {
	Pid      string
	Name     string
	Approval approval.Approval
}

type ProcessDTO struct {
	Id       int64
	Pid      string
	Approval approval.Approval
	Created  time.Time
}

type InsertFlowReqDTO struct {
	ProcessId  int64
	Approval   approval.Approval
	CurrIndex  int
	FlowStatus FlowStatus
	Creator    string
}

type InsertNotifyReqDTO struct {
	FlowId    int64
	Accounts  []string
	Done      bool
	FlowIndex int
}

type InsertDetailReqDTO struct {
	FlowId    int64
	FlowIndex int
	FlowOp    FlowOp
	Account   string
}
