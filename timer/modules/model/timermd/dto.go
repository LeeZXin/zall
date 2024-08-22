package timermd

import (
	"github.com/LeeZXin/zall/pkg/timer"
	"time"
)

type InsertTimerReqDTO struct {
	Name      string
	CronExp   string
	Content   timer.Task
	TeamId    int64
	Env       string
	IsEnabled bool
	Creator   string
}

type UpdateTimerReqDTO struct {
	Id      int64
	Name    string
	CronExp string
	Content timer.Task
}

type ListTimerReqDTO struct {
	TeamId   int64
	Name     string
	PageNum  int
	PageSize int
	Env      string
}

type InsertLogReqDTO struct {
	TimerId     int64
	TaskContent timer.Task
	ErrLog      string
	TriggerType timer.TriggerType
	TriggerBy   string
	IsSuccess   bool
}

type ListLogReqDTO struct {
	TimerId   int64
	PageNum   int
	PageSize  int
	BeginTime time.Time
	EndTime   time.Time
}

type InsertExecuteReqDTO struct {
	TimerId   int64
	IsEnabled bool
	NextTime  int64
	Env       string
}
