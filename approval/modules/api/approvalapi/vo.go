package approvalapi

import (
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type AgreeApprovalReqVO struct {
	NotifyId int64 `json:"notifyId"`
}

type DisagreeApprovalReqVO struct {
	NotifyId int64 `json:"notifyId"`
}

type InsertCustomProcessReqVO struct {
	Pid     string              `json:"pid"`
	Name    string              `json:"name"`
	GroupId int64               `json:"groupId"`
	IconUrl string              `json:"iconUrl"`
	Process approval.ProcessCfg `json:"process"`
}

type DeleteCustomProcessReqVO struct {
	Id int64 `json:"id"`
}

type ListCustomProcessReqVO struct {
	GroupId int64 `json:"groupId"`
}

type UpdateCustomProcessReqVO struct {
	Id      int64               `json:"id"`
	Name    string              `json:"name"`
	GroupId int64               `json:"groupId"`
	IconUrl string              `json:"iconUrl"`
	Process approval.ProcessCfg `json:"process"`
}

type InsertCustomFlowReqVO struct {
	Pid string        `json:"pid"`
	Kvs []approval.Kv `json:"kvs"`
}

type InsertCustomFlowRespVO struct {
	ginutil.BaseResp
	ErrKeys []string `json:"errKeys"`
}

type CancelCustomFlowReqVO struct {
	FlowId int64 `json:"flowId"`
}

type GroupProcessVO struct {
	Id        int64             `json:"id"`
	Name      string            `json:"name"`
	Processes []SimpleProcessVO `json:"processes"`
}

type SimpleProcessVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	IconUrl string `json:"iconUrl"`
}

type ProcessVO struct {
	Id      int64            `json:"id"`
	Pid     string           `json:"pid"`
	GroupId int64            `json:"groupId"`
	Name    string           `json:"name"`
	Content approval.Process `json:"content"`
	IconUrl string           `json:"iconUrl"`
	Created string           `json:"created"`
}

type GetFlowReqVO struct {
	FlowId int64 `json:"flowId"`
}

type FlowDetailVO struct {
	Id          int64            `json:"flowId"`
	ProcessName string           `json:"processName"`
	FlowStatus  string           `json:"flowStatus"`
	Creator     string           `json:"creator"`
	Created     string           `json:"created"`
	Kvs         []approval.Kv    `json:"kvs"`
	Process     approval.Process `json:"process"`
	NotifyList  []NotifyVO       `json:"notifyList"`
}

type FlowVO struct {
	Id          int64  `json:"flowId"`
	ProcessName string `json:"processName"`
	FlowStatus  string `json:"flowStatus"`
	Creator     string `json:"creator"`
	Created     string `json:"created"`
}

type NotifyVO struct {
	Account   string `json:"account"`
	FlowIndex int    `json:"flowIndex"`
	Done      bool   `json:"done"`
	Op        string `json:"op"`
	Updated   string `json:"updated" xorm:"updated"`
}

type GetFlowRespVO struct {
	ginutil.BaseResp
	Data FlowDetailVO `json:"data"`
}

type InsertGroupReqVO struct {
	Name string `json:"name"`
}

type DeleteGroupReqVO struct {
	Id int64 `json:"id"`
}

type UpdateGroupReqVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GroupVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ListCustomFlowReqVO struct {
	DayTime string `json:"dayTime"`
}

type ListOperateFlowReqVO struct {
	DayTime string `json:"dayTime"`
	Done    bool   `json:"done"`
}
