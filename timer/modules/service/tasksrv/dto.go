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
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EnableTaskReqDTO) IsValid() error {
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
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisableTaskReqDTO) IsValid() error {
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
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTaskReqDTO) IsValid() error {
	if r.Id <= 0 {
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
	Creator   string     `json:"creator"`
}

type ListTaskLogReqDTO struct {
	TaskId    int64               `json:"taskId"`
	PageNum   int                 `json:"pageNum"`
	Month     string              `json:"dateStr"`
	Operator  apisession.UserInfo `json:"operator"`
	monthTime time.Time
}

func (r *ListTaskLogReqDTO) IsValid() error {
	if r.TaskId <= 0 || r.PageNum <= 0 || r.Month == "" {
		return util.InvalidArgsError()
	}
	var err error
	r.monthTime, err = time.Parse("2006-01", r.Month)
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
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TriggerTaskReqDTO) IsValid() error {
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

type GetFailedTaskNotifyTplIdReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetFailedTaskNotifyTplIdReqDTO) IsValid() error {
	if r.TeamId <= 0 {
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

type BindFailedTaskNotifyTplReqDTO struct {
	TeamId   int64               `json:"teamId"`
	TplId    int64               `json:"tplId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *BindFailedTaskNotifyTplReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	// 可以为0 0代表没有
	if r.TplId < 0 {
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
