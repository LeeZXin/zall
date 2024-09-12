package timersrv

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zall/timer/modules/model/timermd"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateTimerReqDTO struct {
	Name     string              `json:"name"`
	CronExp  string              `json:"cronExp"`
	Task     timer.Task          `json:"task"`
	TeamId   int64               `json:"teamId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateTimerReqDTO) IsValid() error {
	if !timermd.IsTimerNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	schd, err := ParseCron(r.CronExp)
	if err != nil {
		return util.InvalidArgsError()
	}
	now := time.Now()
	if schd.Next(now).Before(now) {
		return util.InvalidArgsError()
	}
	if !r.Task.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTimerReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	PageNum  int                 `json:"pageNum"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTimerReqDTO) IsValid() error {
	if len(r.Name) > 0 && !timermd.IsTimerNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 || r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type EnableTimerReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EnableTimerReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisableTimerReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisableTimerReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTimerReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTimerReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TimerDTO struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	CronExp   string     `json:"cronExp"`
	Task      timer.Task `json:"task"`
	TeamId    int64      `json:"teamId"`
	IsEnabled bool       `json:"isEnabled"`
	Env       string     `json:"env"`
	Creator   util.User  `json:"creator"`
}

type ListLogReqDTO struct {
	Id        int64               `json:"id"`
	PageNum   int                 `json:"pageNum"`
	Month     string              `json:"dateStr"`
	Operator  apisession.UserInfo `json:"operator"`
	monthTime time.Time
}

func (r *ListLogReqDTO) IsValid() error {
	if r.Id <= 0 || r.PageNum <= 0 || r.Month == "" {
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

type LogDTO struct {
	Task        timer.Task
	ErrLog      string
	TriggerType timer.TriggerType
	TriggerBy   util.User
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

type UpdateTimerReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	CronExp  string              `json:"cronExp"`
	Task     timer.Task          `json:"task"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTimerReqDTO) IsValid() error {
	if !timermd.IsTimerNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	schd, err := ParseCron(r.CronExp)
	if err != nil {
		return util.InvalidArgsError()
	}
	now := time.Now()
	if schd.Next(now).Before(now) {
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
