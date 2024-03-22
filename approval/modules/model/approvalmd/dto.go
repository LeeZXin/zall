package approvalmd

import (
	"github.com/LeeZXin/zall/pkg/approval"
)

type InsertProcessReqDTO struct {
	Pid        string
	Name       string
	GroupId    int64
	IconUrl    string
	SourceType SourceType
	Process    approval.Process
}

type UpdateProcessByIdReqDTO struct {
	Id      int64
	Name    string
	GroupId int64
	IconUrl string
	Process approval.Process
}

type UpdateProcessByPidReqDTO struct {
	Pid     string
	Name    string
	GroupId int64
	Process approval.Process
}

type InsertFlowReqDTO struct {
	ProcessId   int64
	ProcessName string
	Process     approval.Process
	CurrIndex   int
	FlowStatus  FlowStatus
	Creator     string
	BizId       string
	Kvs         []approval.Kv
}

type InsertNotifyReqDTO struct {
	FlowId    int64
	Accounts  []string
	Done      bool
	FlowIndex int
}

type InsertGroupReqDTO struct {
	Name string
}

type UpdateGroupReqDTO struct {
	Id   int64
	Name string
}
