package approvalmd

import (
	"github.com/LeeZXin/zall/pkg/approval"
	"time"
)

type InsertProcessReqDTO struct {
	Pid     string
	Name    string
	GroupId int64
	Process approval.Process
}

type UpdateProcessByIdReqDTO struct {
	Id      int64
	Name    string
	GroupId int64
	Process approval.Process
}

type UpdateProcessByPidReqDTO struct {
	Pid     string
	Name    string
	GroupId int64
	Process approval.Process
}

type ProcessDTO struct {
	Id      int64
	Pid     string
	GroupId int64
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
	Kvs        []approval.Kv
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

type InsertGroupReqDTO struct {
	Name string
}

type UpdateGroupReqDTO struct {
	Id   int64
	Name string
}
