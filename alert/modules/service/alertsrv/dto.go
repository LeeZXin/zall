package alertsrv

import (
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type InsertConfigReqDTO struct {
	Name        string              `json:"name"`
	Alert       alert.Alert         `json:"alert"`
	AppId       string              `json:"appId"`
	IntervalSec int                 `json:"intervalSec"`
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
	if r.IntervalSec <= 0 || r.IntervalSec%5 != 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
