package tasksrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"time"
)

const (
	HttpTaskType = "http"
)

type InsertTaskReqDTO struct {
	Name     string              `json:"name"`
	CronExp  string              `json:"cronExp"`
	TaskType string              `json:"taskType"`
	HttpTask HttpTask            `json:"httpTask"`
	TeamId   int64               `json:"teamId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertTaskReqDTO) IsValid() error {
	if !taskmd.IsTaskNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	_, err := parseCron(r.CronExp)
	if err != nil {
		return util.InvalidArgsError()
	}
	switch r.TaskType {
	case HttpTaskType:
		if !r.HttpTask.IsValid() {
			return util.InvalidArgsError()
		}
	default:
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTaskReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTaskReqDTO) IsValid() error {
	if len(r.Name) > 0 && !taskmd.IsTaskNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 || r.Limit < 0 || r.Cursor < 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	return nil
}

type EnableTaskReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EnableTaskReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisableTaskReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisableTaskReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTaskReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTaskReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TaskDTO struct {
	Id         int64             `json:"id"`
	Name       string            `json:"name"`
	CronExp    string            `json:"cronExp"`
	TaskType   string            `json:"taskType"`
	HttpTask   HttpTask          `json:"httpTask"`
	TeamId     int64             `json:"teamId"`
	NextTime   int64             `json:"nextTime"`
	TaskStatus taskmd.TaskStatus `json:"taskStatus"`
}

type ListTaskLogReqDTO struct {
	Id       int64               `json:"id"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTaskLogReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 || r.Cursor < 0 || r.Limit < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TaskLogDTO struct {
	TaskType    string             `json:"taskType"`
	HttpTask    HttpTask           `json:"httpTask"`
	LogContent  string             `json:"logContent"`
	TriggerType taskmd.TriggerType `json:"triggerType"`
	TriggerBy   string             `json:"triggerBy"`
	TaskStatus  taskmd.TaskStatus  `json:"taskStatus"`
	Created     time.Time          `json:"created"`
}

type TriggerTaskReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TriggerTaskReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTaskReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	CronExp  string              `json:"cronExp"`
	TaskType string              `json:"taskType"`
	HttpTask HttpTask            `json:"httpTask"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTaskReqDTO) IsValid() error {
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if !taskmd.IsTaskNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	_, err := parseCron(r.CronExp)
	if err != nil {
		return util.InvalidArgsError()
	}
	switch r.TaskType {
	case HttpTaskType:
		if !r.HttpTask.IsValid() {
			return util.InvalidArgsError()
		}
	default:
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
