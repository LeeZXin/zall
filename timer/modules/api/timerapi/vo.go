package timerapi

import (
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zall/util"
)

type CreateTimerReqVO struct {
	Name    string     `json:"name"`
	CronExp string     `json:"cronExp"`
	Task    timer.Task `json:"task"`
	TeamId  int64      `json:"teamId"`
	Env     string     `json:"env"`
}

type ListTimerReqVO struct {
	TeamId  int64  `json:"teamId"`
	Name    string `json:"name"`
	PageNum int    `json:"pageNum"`
	Env     string `json:"env"`
}

type TimerVO struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	CronExp   string     `json:"cronExp"`
	Task      timer.Task `json:"task"`
	TeamId    int64      `json:"teamId"`
	IsEnabled bool       `json:"isEnabled"`
	Env       string     `json:"env"`
	Creator   util.User  `json:"creator"`
}

type ListLogReqVO struct {
	Id      int64  `json:"id"`
	PageNum int    `json:"pageNum"`
	Month   string `json:"month"`
}

type TaskLogVO struct {
	Task        timer.Task        `json:"task"`
	ErrLog      string            `json:"errLog"`
	TriggerType timer.TriggerType `json:"triggerType"`
	TriggerBy   util.User         `json:"triggerBy"`
	IsSuccess   bool              `json:"isSuccess"`
	Created     string            `json:"created"`
}

type UpdateTimerReqVO struct {
	Id      int64      `json:"id"`
	Name    string     `json:"name"`
	CronExp string     `json:"cronExp"`
	Task    timer.Task `json:"task"`
}
