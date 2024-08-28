package alertsrv

import (
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"math"
)

type CreateConfigReqDTO struct {
	Name        string              `json:"name"`
	Alert       alert.Alert         `json:"alert"`
	AppId       string              `json:"appId"`
	IntervalSec int                 `json:"intervalSec"`
	Env         string              `json:"env"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *CreateConfigReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !alertmd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Alert.IsValid() {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if r.IntervalSec < 10 || r.IntervalSec > 3600 {
		return util.InvalidArgsError()
	}
	r.IntervalSec = int(math.Floor(float64(r.IntervalSec)/10) * 10)
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateConfigReqDTO struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name"`
	Alert       alert.Alert         `json:"alert"`
	IntervalSec int                 `json:"intervalSec"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *UpdateConfigReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !alertmd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Alert.IsValid() {
		return util.InvalidArgsError()
	}
	if r.IntervalSec < 10 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteConfigReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteConfigReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type EnableConfigReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EnableConfigReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisableConfigReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisableConfigReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListConfigReqDTO struct {
	PageNum  int                 `json:"pageNum"`
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListConfigReqDTO) IsValid() error {
	if r.PageNum < 0 {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
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

type ConfigDTO struct {
	Id          int64
	Name        string
	AppId       string
	Content     alert.Alert
	IntervalSec int
	IsEnabled   bool
	Creator     string
	Env         string
}
