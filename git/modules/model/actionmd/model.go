package actionmd

import (
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

type TriggerType int

const (
	WebhookTriggerType TriggerType = iota
	ManualTriggerType
)

func (t TriggerType) Readable() string {
	switch t {
	case WebhookTriggerType:
		return i18n.GetByKey(i18n.ActionWebhookTriggerType)
	case ManualTriggerType:
		return i18n.GetByKey(i18n.ActionManualTriggerType)
	default:
		return i18n.GetByKey(i18n.ActionUnknownTriggerType)
	}
}

type TaskStatus int

const (
	TaskRunningStatus TaskStatus = iota + 1
	TaskSuccessStatus
	TaskFailStatus
)

func (t TaskStatus) Readable() string {
	switch t {
	case TaskRunningStatus:
		return i18n.GetByKey(i18n.ActionTaskRunningStatus)
	case TaskSuccessStatus:
		return i18n.GetByKey(i18n.ActionTaskSuccessStatus)
	case TaskFailStatus:
		return i18n.GetByKey(i18n.ActionTaskFailStatus)
	default:
		return i18n.GetByKey(i18n.ActionTaskUnknownStatus)
	}
}

const (
	TaskTableName   = "zgit_action_task"
	StepTableName   = "zgit_action_step"
	NodeTableName   = "zgit_action_node"
	ActionTableName = "zgit_action"
)

type Task struct {
	Id            int64       `json:"id" xorm:"pk autoincr"`
	ActionId      int64       `json:"actionId"`
	TaskStatus    TaskStatus  `json:"taskStatus"`
	TriggerType   TriggerType `json:"triggerType"`
	ActionContent string      `json:"actionContent"`
	Operator      string      `json:"operator"`
	Created       time.Time   `json:"created" xorm:"created"`
}

func (*Task) TableName() string {
	return TaskTableName
}

type StepStatus int

const (
	StepWaitingStatus StepStatus = iota
	StepRunningStatus
	StepSuccessStatus
	StepFailStatus
)

func (t StepStatus) Readable() string {
	switch t {
	case StepWaitingStatus:
		return i18n.GetByKey(i18n.ActionStepWaitingStatus)
	case StepRunningStatus:
		return i18n.GetByKey(i18n.ActionStepRunningStatus)
	case StepSuccessStatus:
		return i18n.GetByKey(i18n.ActionStepSuccessStatus)
	case StepFailStatus:
		return i18n.GetByKey(i18n.ActionStepFailStatus)
	default:
		return i18n.GetByKey(i18n.ActionStepUnknownStatus)
	}
}

type Step struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	TaskId     int64      `json:"taskId"`
	JobName    string     `json:"jobName"`
	StepName   string     `json:"stepName"`
	StepIndex  int        `json:"stepIndex"`
	LogContent string     `json:"logContent"`
	StepStatus StepStatus `json:"stepStatus"`
	Created    time.Time  `json:"created" xorm:"created"`
	Updated    time.Time  `json:"updated" xorm:"updated"`
}

func (*Step) TableName() string {
	return StepTableName
}

type Node struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	Name     string    `json:"name"`
	HttpHost string    `json:"httpHost"`
	Token    string    `json:"token"`
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
}

func (*Node) TableName() string {
	return NodeTableName
}

type Action struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Aid     string    `json:"aid"`
	Name    string    `json:"name"`
	TeamId  int64     `json:"teamId"`
	Content string    `json:"content"`
	NodeId  int64     `json:"nodeId"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Action) TableName() string {
	return ActionTableName
}
