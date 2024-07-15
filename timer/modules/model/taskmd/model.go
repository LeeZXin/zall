package taskmd

import (
	"github.com/LeeZXin/zall/pkg/timer"
	"time"
)

const (
	TaskTableName    = "ztimer_task"
	ExecuteTableName = "ztimer_execute"
	LogTableName     = "ztimer_log"

	DefaultTrigger = "system"
)

type TriggerType int

func (t TriggerType) IsValid() bool {
	switch t {
	case AutoTriggerType, ManualTriggerType:
		return true
	default:
		return false
	}
}

const (
	AutoTriggerType = iota + 1
	ManualTriggerType
)

type Task struct {
	Id        int64       `json:"id" xorm:"pk autoincr"`
	Name      string      `json:"name"`
	CronExp   string      `json:"cronExp"`
	Content   *timer.Task `json:"content"`
	TeamId    int64       `json:"teamId"`
	Env       string      `json:"env"`
	IsEnabled bool        `json:"isEnabled"`
	Creator   string      `json:"creator"`
	Created   time.Time   `json:"created" xorm:"created"`
}

func (*Task) TableName() string {
	return TaskTableName
}

func (t *Task) GetContent() timer.Task {
	if t.Content != nil {
		return *t.Content
	}
	return timer.Task{}
}

type Execute struct {
	Id         int64  `json:"id" xorm:"pk autoincr"`
	TaskId     int64  `json:"taskId"`
	IsEnabled  bool   `json:"isEnabled"`
	NextTime   int64  `json:"nextTime"`
	Env        string `json:"env"`
	RunVersion int64  `json:"runVersion"`
}

func (*Execute) TableName() string {
	return ExecuteTableName
}

type TaskLog struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	TaskId      int64       `json:"taskId"`
	TaskContent *timer.Task `json:"taskContent"`
	ErrLog      string      `json:"errLog"`
	TriggerType TriggerType `json:"triggerType"`
	TriggerBy   string      `json:"triggerBy"`
	IsSuccess   bool        `json:"isSuccess"`
	Created     time.Time   `json:"created" xorm:"created"`
}

func (l *TaskLog) GetTaskContent() timer.Task {
	if l.TaskContent != nil {
		return *l.TaskContent
	}
	return timer.Task{}
}

func (*TaskLog) TableName() string {
	return LogTableName
}
