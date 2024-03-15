package taskmd

import (
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	TaskTableName     = "ztimer_task_content"
	LogTableName      = "ztimer_task_log"
	InstanceTableName = "ztimer_instance"

	DefaultTrigger = "system"
)

type TaskStatus int

func (t TaskStatus) IsValid() bool {
	switch t {
	case Pending, Running, Successful, Failed, Closed:
		return true
	default:
		return false
	}
}

func (t TaskStatus) Readable() string {
	switch t {
	case Pending:
		return i18n.GetByKey(i18n.TimerTaskPendingStatus)
	case Running:
		return i18n.GetByKey(i18n.TimerTaskRunningStatus)
	case Successful:
		return i18n.GetByKey(i18n.TimerTaskSuccessfulStatus)
	case Failed:
		return i18n.GetByKey(i18n.TimerTaskFailedStatus)
	case Closed:
		return i18n.GetByKey(i18n.TimerTaskClosedStatus)
	default:
		return i18n.GetByKey(i18n.TimerTaskUnknownStatus)
	}
}

const (
	Pending TaskStatus = iota + 1
	Running
	Successful
	Failed
	Closed
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

func (t TriggerType) Readable() string {
	switch t {
	case AutoTriggerType:
		return i18n.GetByKey(i18n.TimerTaskAutoTriggerType)
	case ManualTriggerType:
		return i18n.GetByKey(i18n.TimerTaskManualTriggerType)
	default:
		return i18n.GetByKey(i18n.TimerTaskUnknownTriggerType)
	}
}

const (
	AutoTriggerType = iota + 1
	ManualTriggerType
)

type Task struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	Name       string     `json:"name"`
	CronExp    string     `json:"cronExp"`
	Content    string     `json:"content"`
	NextTime   int64      `json:"nextTime"`
	TaskStatus TaskStatus `json:"taskStatus"`
	TeamId     int64      `json:"teamId"`
	Version    int64      `json:"version"`
	Created    time.Time  `json:"created" xorm:"created"`
}

func (*Task) TableName() string {
	return TaskTableName
}

type TaskLog struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	TaskId      int64       `json:"taskId"`
	TaskContent string      `json:"taskContent"`
	LogContent  string      `json:"logContent"`
	TriggerType TriggerType `json:"triggerType"`
	TriggerBy   string      `json:"triggerBy"`
	TaskStatus  TaskStatus  `json:"taskStatus"`
	Created     time.Time   `json:"created" xorm:"created"`
}

func (*TaskLog) TableName() string {
	return LogTableName
}

type Instance struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	InstanceId    string    `json:"instanceId"`
	HeartbeatTime int64     `json:"heartbeatTime"`
	Created       time.Time `json:"created" xorm:"created"`
}

func (*Instance) TableName() string {
	return InstanceTableName
}
