package alertsrv

import (
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type InsertConfigReqDTO struct {
	Name        string              `json:"name"`
	Alert       alert.Alert         `json:"alert"`
	AppId       string              `json:"appId"`
	IntervalSec int                 `json:"intervalSec"`
	SilenceSec  int                 `json:"silenceSec"`
	Enabled     bool                `json:"enabled"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *InsertConfigReqDTO) IsValid() error {
	if !alertmd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Alert.IsValid() {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if r.IntervalSec <= 0 || r.IntervalSec%5 != 0 || r.SilenceSec < r.IntervalSec || r.SilenceSec%5 != 0 {
		return util.InvalidArgsError()
	}
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
	SilenceSec  int                 `json:"silenceSec"`
	Enabled     bool                `json:"enabled"`
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
	if r.IntervalSec <= 0 || r.IntervalSec%5 != 0 || r.SilenceSec < r.IntervalSec || r.SilenceSec%5 != 0 {
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

type ListConfigReqDTO struct {
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	AppId    string              `json:"appId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListConfigReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ConfigDTO struct {
	Id          int64
	Name        string
	AppId       string
	Content     *alert.Alert
	IntervalSec int
	SilenceSec  int
	Enabled     bool
	NextTime    int64
	Created     time.Time
}
