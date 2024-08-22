package timermd

import (
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	TimerTableName   = "ztimer_timer"
	ExecuteTableName = "ztimer_execute"
	LogTableName     = "ztimer_log"
)

type Timer struct {
	Id        int64                            `json:"id" xorm:"pk autoincr"`
	Name      string                           `json:"name"`
	CronExp   string                           `json:"cronExp"`
	Content   *xormutil.Conversion[timer.Task] `json:"content"`
	TeamId    int64                            `json:"teamId"`
	Env       string                           `json:"env"`
	IsEnabled bool                             `json:"isEnabled"`
	Creator   string                           `json:"creator"`
	Created   time.Time                        `json:"created" xorm:"created"`
	Updated   time.Time                        `json:"updated" xorm:"updated"`
}

func (*Timer) TableName() string {
	return TimerTableName
}

func (t *Timer) GetContent() timer.Task {
	if t.Content != nil {
		return t.Content.Data
	}
	return timer.Task{}
}

type Execute struct {
	Id         int64  `json:"id" xorm:"pk autoincr"`
	TimerId    int64  `json:"timerId"`
	IsEnabled  bool   `json:"isEnabled"`
	NextTime   int64  `json:"nextTime"`
	Env        string `json:"env"`
	RunVersion int64  `json:"runVersion"`
}

func (*Execute) TableName() string {
	return ExecuteTableName
}

type Log struct {
	Id          int64                            `json:"id" xorm:"pk autoincr"`
	TimerId     int64                            `json:"timerId"`
	TaskContent *xormutil.Conversion[timer.Task] `json:"taskContent"`
	ErrLog      string                           `json:"errLog"`
	TriggerType timer.TriggerType                `json:"triggerType"`
	TriggerBy   string                           `json:"triggerBy"`
	IsSuccess   bool                             `json:"isSuccess"`
	Created     time.Time                        `json:"created" xorm:"created"`
}

func (l *Log) GetTaskContent() timer.Task {
	if l.TaskContent != nil {
		return l.TaskContent.Data
	}
	return timer.Task{}
}

func (*Log) TableName() string {
	return LogTableName
}
