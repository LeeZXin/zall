package gitactionmd

import "time"

type TaskStatus int

const (
	TaskRunningStatus TaskStatus = iota + 1
	TaskSuccessStatus
	TaskFailStatus
)

const (
	TaskTableName   = "zgit_action_task"
	StepTableName   = "zgit_action_step"
	NodeTableName   = "zgit_action_node"
	ActionTableName = "zgit_action"
)

type Task struct {
	Id            int64      `json:"id" xorm:"pk autoincr"`
	TaskName      string     `json:"taskName"`
	ActionId      int64      `json:"actionId"`
	TaskStatus    TaskStatus `json:"taskStatus"`
	TriggerType   int        `json:"triggerType"`
	ActionContent string     `json:"actionContent"`
	Created       time.Time  `json:"created" xorm:"created"`
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
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
}

func (*Node) TableName() string {
	return NodeTableName
}

type Action struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	Name       string    `json:"name"`
	RepoId     int64     `json:"repoId"`
	Content    string    `json:"content"`
	NodeId     int64     `json:"nodeId"`
	WildBranch string    `json:"wildBranch"`
	Created    time.Time `json:"created" xorm:"created"`
	Updated    time.Time `json:"updated" xorm:"updated"`
}

func (*Action) TableName() string {
	return ActionTableName
}
