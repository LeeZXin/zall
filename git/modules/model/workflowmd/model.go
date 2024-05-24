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
	TaskQueueStatus TaskStatus = iota
	TaskRunningStatus
	TaskSuccessStatus
	TaskFailStatus
	TaskCancelStatus
	TaskTimeoutStatus
)

func (s TaskStatus) IsEndType() bool {
	switch s {
	case TaskQueueStatus, TaskRunningStatus:
		return false
	default:
		return true
	}
}

const (
	TaskTableName     = "zgit_workflow_task"
	SecretTableName   = "zgit_workflow_secret"
	WorkflowTableName = "zgit_workflow"
	TokenTableName    = "zgit_workflow_token"
)

type Task struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	WorkflowId  int64       `json:"workflowId"`
	TaskStatus  TaskStatus  `json:"taskStatus"`
	TriggerType TriggerType `json:"triggerType"`
	YamlContent string      `json:"yamlContent"`
	AgentHost   string      `json:"agentHost"`
	AgentToken  string      `json:"agentToken"`
	Branch      string      `json:"branch"`
	Operator    string      `json:"operator"`
	PrId        int64       `json:"prId"`
	Duration    int64       `json:"duration"`
	BizId       string      `json:"bizId"`
	Created     time.Time   `json:"created" xorm:"created"`
	Updated     time.Time   `json:"updated" xorm:"updated"`
}

func (*Task) TableName() string {
	return TaskTableName
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

type Secret struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	RepoId  int64     `json:"repoId"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Secret) TableName() string {
	return SecretTableName
}

type Token struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	TaskId   int64     `json:"taskId"`
	RepoId   int64     `json:"repoId"`
	Content  string    `json:"content"`
	Operator string    `json:"operator"`
	Expired  time.Time `json:"expired"`
	Created  time.Time `json:"created" xorm:"created"`
}

func (*Token) TableName() string {
	return TokenTableName
}

func (t *Token) IsExpired() bool {
	return t.Expired.Before(time.Now())
}
