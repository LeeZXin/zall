package actiontaskmd

import "time"

type TaskType int

const (
	GitTaskType TaskType = iota
)

type TaskStatus int

const (
	TaskSuspendStatus TaskStatus = iota
	TaskWaitingStatus
	TaskRunningStatus
	TaskSuccessStatus
	TaskFailStatus
	TaskErrStatus
)

type TriggerType int

func (t TriggerType) IsValid() bool {
	switch t {
	case SysTriggerType, ManualTriggerType:
		return true
	}
	return false
}

func (t TriggerType) Int() int {
	return int(t)
}

const (
	SysTriggerType TriggerType = iota
	ManualTriggerType
)

const (
	TaskTableName     = "action_task"
	StepTableName     = "action_step"
	InstanceTableName = "action_instance"
)

type Task struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	TaskName    string      `json:"taskName"`
	RepoId      int64       `json:"repoId"`
	TaskType    TaskType    `json:"taskType"`
	InstanceId  string      `json:"instanceId"`
	TaskStatus  TaskStatus  `json:"taskStatus"`
	TriggerType TriggerType `json:"triggerType"`
	HookContent string      `json:"hookContent"`
	Created     time.Time   `json:"created" xorm:"created"`
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

type Instance struct {
	Id           int64     `json:"id" xorm:"pk autoincr"`
	InstanceId   string    `json:"instanceId"`
	Name         string    `json:"name"`
	InstanceHost string    `json:"instanceHost"`
	JobCount     int64     `json:"jobCount"`
	Heartbeat    int64     `json:"heartbeat"`
	Created      time.Time `json:"created" xorm:"created"`
}

func (*Instance) TableName() string {
	return InstanceTableName
}
