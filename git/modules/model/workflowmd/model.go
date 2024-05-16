package workflowmd

import (
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

type TriggerType int

const (
	HookTriggerType TriggerType = iota + 1
	ManualTriggerType
)

func (t TriggerType) Readable() string {
	switch t {
	case HookTriggerType:
		return i18n.GetByKey(i18n.WorkflowHookTriggerType)
	case ManualTriggerType:
		return i18n.GetByKey(i18n.WorkflowManualTriggerType)
	default:
		return i18n.GetByKey(i18n.WorkflowUnknownTriggerType)
	}
}

func (t TriggerType) IsValid() bool {
	switch t {
	case HookTriggerType, ManualTriggerType:
		return true
	default:
		return false
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
		return i18n.GetByKey(i18n.WorkflowTaskRunningStatus)
	case TaskSuccessStatus:
		return i18n.GetByKey(i18n.WorkflowTaskSuccessStatus)
	case TaskFailStatus:
		return i18n.GetByKey(i18n.WorkflowTaskFailStatus)
	default:
		return i18n.GetByKey(i18n.WorkflowTaskUnknownStatus)
	}
}

const (
	TaskTableName     = "zgit_workflow_task"
	StepTableName     = "zgit_workflow_step"
	WorkflowTableName = "zgit_workflow"
)

type Task struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	WorkflowId  int64       `json:"workflowId"`
	TaskStatus  TaskStatus  `json:"taskStatus"`
	TriggerType TriggerType `json:"triggerType"`
	YamlContent string      `json:"yamlContent"`
	Branch      string      `json:"branch"`
	Operator    string      `json:"operator"`
	PrId        int64       `json:"prId"`
	Created     time.Time   `json:"created" xorm:"created"`
	Updated     time.Time   `json:"updated" xorm:"updated"`
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
		return i18n.GetByKey(i18n.WorkflowStepWaitingStatus)
	case StepRunningStatus:
		return i18n.GetByKey(i18n.WorkflowStepRunningStatus)
	case StepSuccessStatus:
		return i18n.GetByKey(i18n.WorkflowStepSuccessStatus)
	case StepFailStatus:
		return i18n.GetByKey(i18n.WorkflowStepFailStatus)
	default:
		return i18n.GetByKey(i18n.WorkflowStepUnknownStatus)
	}
}

type Step struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	TaskId     int64      `json:"taskId"`
	WorkflowId int64      `json:"workflowId"`
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

type Workflow struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RepoId      int64     `json:"repoId"`
	YamlContent string    `json:"yamlContent"`
	Source      *Source   `json:"source"`
	AgentHost   string    `json:"agentHost"`
	AgentToken  string    `json:"agentToken"`
	LastTaskId  int64     `json:"lastTaskId"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (*Workflow) TableName() string {
	return WorkflowTableName
}
