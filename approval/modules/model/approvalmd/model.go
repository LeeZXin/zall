package approvalmd

import (
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

type FlowOp int

const (
	PendingOp FlowOp = iota
	AgreeOp
	DisagreeOp
	CancelOp
	AutoAgreeOp
	AutoDisagreeOp
)

func (o FlowOp) Readable() string {
	switch o {
	case PendingOp:
		return i18n.GetByKey(i18n.FlowPendingOp)
	case AgreeOp:
		return i18n.GetByKey(i18n.FlowAgreeOp)
	case DisagreeOp:
		return i18n.GetByKey(i18n.FlowDisagreeOp)
	case CancelOp:
		return i18n.GetByKey(i18n.FlowCancelOp)
	case AutoAgreeOp:
		return i18n.GetByKey(i18n.FlowAutoAgreeOp)
	case AutoDisagreeOp:
		return i18n.GetByKey(i18n.FlowAutoDisagreeOp)
	default:
		return i18n.GetByKey(i18n.FlowUnknownOp)
	}
}

type FlowStatus int

const (
	PendingFlowStatus FlowStatus = iota
	ErrFlowStatus
	CanceledFlowStatus
	AgreeFlowStatus
	DisagreeFlowStatus
)

func (s FlowStatus) Readable() string {
	switch s {
	case PendingFlowStatus:
		return i18n.GetByKey(i18n.FlowPendingStatus)
	case AgreeFlowStatus:
		return i18n.GetByKey(i18n.FlowAgreeStatus)
	case DisagreeFlowStatus:
		return i18n.GetByKey(i18n.FlowDisagreeStatus)
	case ErrFlowStatus:
		return i18n.GetByKey(i18n.FlowErrStatus)
	case CanceledFlowStatus:
		return i18n.GetByKey(i18n.FlowCanceledStatus)
	default:
		return i18n.GetByKey(i18n.FlowUnknownStatus)
	}
}

func (s FlowStatus) IsValid() bool {
	switch s {
	case PendingFlowStatus, AgreeFlowStatus, DisagreeFlowStatus, ErrFlowStatus, CanceledFlowStatus:
		return true
	default:
		return false
	}
}

type SourceType int

const (
	CustomSourceType SourceType = iota
	SystemSourceType
)

const (
	ProcessTableName = "zapproval_process"
	FlowTableName    = "zapproval_flow"
	NotifyTableName  = "zapproval_notify"
	GroupTableName   = "zapproval_group"
)

type Process struct {
	Id         int64             `json:"id" xorm:"pk autoincr"`
	Pid        string            `json:"pid"`
	GroupId    int64             `json:"groupId"`
	Name       string            `json:"name"`
	Content    *approval.Process `json:"content"`
	IconUrl    string            `json:"iconUrl"`
	SourceType SourceType        `json:"sourceType"`
	Created    time.Time         `json:"created" xorm:"created"`
	Updated    time.Time         `json:"updated" xorm:"updated"`
}

func (*Process) TableName() string {
	return ProcessTableName
}

type SimpleProcess struct {
	Id      int64  `json:"id" xorm:"pk autoincr"`
	Pid     string `json:"pid"`
	GroupId int64  `json:"groupId"`
	Name    string `json:"name"`
	IconUrl string `json:"iconUrl"`
}

func (*SimpleProcess) TableName() string {
	return ProcessTableName
}

type Flow struct {
	Id             int64             `json:"id" xorm:"pk autoincr"`
	ProcessId      int64             `json:"processId"`
	ProcessName    string            `json:"processName"`
	ProcessContent *approval.Process `json:"processContent"`
	CurrIndex      int               `json:"currIndex"`
	FlowStatus     FlowStatus        `json:"flowStatus"`
	ErrMsg         string            `json:"errMsg"`
	Creator        string            `json:"creator"`
	BizId          string            `json:"bizId"`
	Kvs            approval.Kvs      `json:"kvs"`
	Created        time.Time         `json:"created" xorm:"created"`
	Updated        time.Time         `json:"updated" xorm:"updated"`
}

func (*Flow) TableName() string {
	return FlowTableName
}

type Notify struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	FlowId    int64     `json:"flowId"`
	Account   string    `json:"account"`
	FlowIndex int       `json:"flowIndex"`
	Done      bool      `json:"done"`
	Op        FlowOp    `json:"op"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (*Notify) TableName() string {
	return NotifyTableName
}

type Group struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Name    string    `json:"name"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Group) TableName() string {
	return GroupTableName
}
