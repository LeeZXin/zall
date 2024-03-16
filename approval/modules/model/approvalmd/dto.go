package approvalmd

import (
	"github.com/LeeZXin/zall/approval/approval"
	"time"
)

type InsertProcessReqDTO struct {
	Pid     string
	Name    string
	Process approval.Process
}

type UpdateProcessReqDTO struct {
	Id      int64
	Name    string
	Process approval.Process
}

type UpdateProcessByPidReqDTO struct {
	Pid     string
	Name    string
	Process approval.Process
}

type ProcessDTO struct {
	Id      int64
	Pid     string
	Process approval.Process
	Created time.Time
}

type InsertFlowReqDTO struct {
	ProcessId  int64
	Process    approval.Process
	CurrIndex  int
	FlowStatus FlowStatus
	Creator    string
	BizId      string
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
