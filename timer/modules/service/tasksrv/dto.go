package tasksrv

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateTaskReqDTO struct {
	Name     string              `json:"name"`
	CronExp  string              `json:"cronExp"`
	Task     timer.Task          `json:"task"`
	TeamId   int64               `json:"teamId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateTaskReqDTO) IsValid() error {
	if !taskmd.IsTaskNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	cron, err := ParseCron(r.CronExp)
	if err != nil {
		return util.InvalidArgsError()
	}
	now := time.Now()
	if cron.Next(now).Before(now) {
		return util.InvalidArgsError()
	}
	if !r.Task.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
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
	PageNum  int                 `json:"pageNum"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTaskReqDTO) IsValid() error {
	if len(r.Name) > 0 && !taskmd.IsTaskNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 || r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type EnableTaskReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EnableTaskReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisableTaskReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisableTaskReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTaskReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTaskReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TaskDTO struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	CronExp   string     `json:"cronExp"`
	Task      timer.Task `json:"task"`
	TeamId    int64      `json:"teamId"`
	IsEnabled bool       `json:"isEnabled"`
	Env       string     `json:"env"`
}

type PageTaskLogReqDTO struct {
	TaskId   int64               `json:"taskId"`
	PageNum  int                 `json:"pageNum"`
	DateStr  string              `json:"dateStr"`
	Operator apisession.UserInfo `json:"operator"`

	dateTime time.Time
}

func (r *PageTaskLogReqDTO) IsValid() error {
	if r.TaskId <= 0 || r.PageNum <= 0 || r.DateStr == "" {
		return util.InvalidArgsError()
	}
	var err error
	r.dateTime, err = time.Parse(time.DateOnly, r.DateStr)
	if err != nil {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TaskLogDTO struct {
	Task        timer.Task
	ErrLog      string
	TriggerType taskmd.TriggerType
	TriggerBy   string
	IsSuccess   bool
	Created     time.Time
}

type TriggerTaskReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TriggerTaskReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTaskReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Name     string              `json:"name"`
	CronExp  string              `json:"cronExp"`
	Task     timer.Task          `json:"task"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTaskReqDTO) IsValid() error {
	if !taskmd.IsTaskNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	_, err := ParseCron(r.CronExp)
	if err != nil {
		return util.InvalidArgsError()
	}
	if !r.Task.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
