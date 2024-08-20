package taskmd

import (
	"github.com/LeeZXin/zall/pkg/timer"
	"time"
)

type InsertTaskReqDTO struct {
	Name      string
	CronExp   string
	Content   timer.Task
	TeamId    int64
	Env       string
	IsEnabled bool
	Creator   string
}

type UpdateTaskReqDTO struct {
	Id      int64
	Name    string
	CronExp string
	Content timer.Task
}

type PageTaskReqDTO struct {
	TeamId   int64
	Name     string
	PageNum  int
	PageSize int
	Env      string
}

type InsertTaskLogReqDTO struct {
	TaskId      int64
	TaskContent *timer.Task
	ErrLog      string
	TriggerType TriggerType
	TriggerBy   string
	IsSuccess   bool
}

type ListTaskLogReqDTO struct {
	TaskId    int64
	PageNum   int
	PageSize  int
	BeginTime time.Time
	EndTime   time.Time
}

type InsertExecuteReqDTO struct {
	TaskId    int64
	IsEnabled bool
	NextTime  int64
	Env       string
}
