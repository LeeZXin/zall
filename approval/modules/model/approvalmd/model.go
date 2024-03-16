package approvalmd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/approval/approval"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

type FlowOp int

const (
	AgreeOp FlowOp = iota
	DisagreeOp
	CancelOp
)

func (o FlowOp) Readable() string {
	switch o {
	case AgreeOp:
		return i18n.GetByKey(i18n.FlowAgreeOp)
	case DisagreeOp:
		return i18n.GetByKey(i18n.FlowDisagreeOp)
	case CancelOp:
		return i18n.GetByKey(i18n.FlowCancelOp)
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

const (
	ProcessTableName = "zapproval_process"
	FlowTableName    = "zapproval_flow"
	DetailTableName  = "zapproval_detail"
	NotifyTableName  = "zapproval_notify"
)

type Process struct {
	Id      int64     `xorm:"pk autoincr"`
	Pid     string    `json:"pid"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Process) TableName() string {
	return ProcessTableName
}

func (p *Process) GetUnmarshalProcess() (approval.Process, error) {
	var ret approval.Process
	err := json.Unmarshal([]byte(p.Content), &ret)
	return ret, err
}

type Flow struct {
	Id             int64      `xorm:"pk autoincr"`
	ProcessId      int64      `json:"processId"`
	ProcessContent string     `json:"processContent"`
	CurrIndex      int        `json:"currIndex"`
	FlowStatus     FlowStatus `json:"flowStatus"`
	ErrMsg         string     `json:"errMsg"`
	Creator        string     `json:"creator"`
	BizId          string     `json:"bizId"`
	Created        time.Time  `json:"created" xorm:"created"`
	Updated        time.Time  `json:"updated" xorm:"updated"`
}

func (f *Flow) GetProcess() (approval.Process, error) {
	var ret approval.Process
	err := json.Unmarshal([]byte(f.ProcessContent), &ret)
	return ret, err
}

func (*Flow) TableName() string {
	return FlowTableName
}

type Detail struct {
	Id        int64     `xorm:"pk autoincr"`
	FlowId    int64     `json:"flowId"`
	Account   string    `json:"account"`
	FlowIndex int       `json:"flowIndex"`
	Op        FlowOp    `json:"op"`
	Created   time.Time `json:"created" xorm:"created"`
}

func (*Detail) TableName() string {
	return DetailTableName
}

type Notify struct {
	Id        int64     `xorm:"pk autoincr"`
	FlowId    int64     `json:"flowId"`
	Account   string    `json:"account"`
	FlowIndex int       `json:"flowIndex"`
	Done      bool      `json:"done"`
	Created   time.Time `json:"created" xorm:"created"`
}

func (*Notify) TableName() string {
	return NotifyTableName
}
