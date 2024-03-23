package taskapi

import (
	"github.com/LeeZXin/zall/timer/modules/service/tasksrv"
)

type InsertTaskReqVO struct {
	Name     string           `json:"name"`
	CronExp  string           `json:"cronExp"`
	TaskType string           `json:"taskType"`
	HttpTask tasksrv.HttpTask `json:"httpTask"`
	TeamId   int64            `json:"teamId"`
	Env      string           `json:"env"`
}

type EnabledTaskReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type DisableTaskReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type DeleteTaskReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type TriggerTaskReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type ListTaskReqVO struct {
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
	Cursor int64  `json:"cursor"`
	Limit  int    `json:"limit"`
	Env    string `json:"env"`
}

type TaskVO struct {
	Id         int64            `json:"id"`
	Name       string           `json:"name"`
	CronExp    string           `json:"cronExp"`
	TaskType   string           `json:"taskType"`
	HttpTask   tasksrv.HttpTask `json:"httpTask"`
	TeamId     int64            `json:"teamId"`
	NextTime   string           `json:"nextTime"`
	TaskStatus string           `json:"taskStatus"`
}

type ListLogReqVO struct {
	Id     int64  `json:"id"`
	Cursor int64  `json:"cursor"`
	Limit  int    `json:"limit"`
	Env    string `json:"env"`
}

type TaskLogVO struct {
	TaskType    string           `json:"taskType"`
	HttpTask    tasksrv.HttpTask `json:"httpTask"`
	LogContent  string           `json:"logContent"`
	TriggerType string           `json:"triggerType"`
	TriggerBy   string           `json:"triggerBy"`
	TaskStatus  string           `json:"taskStatus"`
	Created     string           `json:"created"`
}

type UpdateTaskReqVO struct {
	Id       int64            `json:"id"`
	Name     string           `json:"name"`
	CronExp  string           `json:"cronExp"`
	TaskType string           `json:"taskType"`
	HttpTask tasksrv.HttpTask `json:"httpTask"`
	Env      string           `json:"env"`
}
