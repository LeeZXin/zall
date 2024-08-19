package taskapi

import (
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
)

type CreateTaskReqVO struct {
	Name    string     `json:"name"`
	CronExp string     `json:"cronExp"`
	Task    timer.Task `json:"task"`
	TeamId  int64      `json:"teamId"`
	Env     string     `json:"env"`
}

type ListTaskReqVO struct {
	TeamId  int64  `json:"teamId"`
	Name    string `json:"name"`
	PageNum int    `json:"pageNum"`
	Env     string `json:"env"`
}

type TaskVO struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	CronExp   string     `json:"cronExp"`
	Task      timer.Task `json:"task"`
	TeamId    int64      `json:"teamId"`
	IsEnabled bool       `json:"isEnabled"`
	Env       string     `json:"env"`
	Creator   string     `json:"creator"`
}

type ListLogReqVO struct {
	TaskId  int64  `json:"taskId"`
	PageNum int    `json:"pageNum"`
	Month   string `json:"month"`
}

type TaskLogVO struct {
	Task        timer.Task         `json:"task"`
	ErrLog      string             `json:"errLog"`
	TriggerType taskmd.TriggerType `json:"triggerType"`
	TriggerBy   string             `json:"triggerBy"`
	IsSuccess   bool               `json:"isSuccess"`
	Created     string             `json:"created"`
}

type UpdateTaskReqVO struct {
	TaskId  int64      `json:"taskId"`
	Name    string     `json:"name"`
	CronExp string     `json:"cronExp"`
	Task    timer.Task `json:"task"`
}

type BindFailedTaskNotifyTplReqVO struct {
	TeamId int64  `json:"teamId"`
	TplId  int64  `json:"tplId"`
	Env    string `json:"env"`
}
