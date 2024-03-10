package detectsrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/tcpdetect/modules/model/detectmd"
	"github.com/LeeZXin/zall/util"
	"regexp"
	"time"
)

var (
	validIpPattern = regexp.MustCompile(`^((\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])$`)
)

type InsertDetectReqDTO struct {
	Ip       string              `json:"ip"`
	Port     int                 `json:"port"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertDetectReqDTO) IsValid() error {
	if !validIpPattern.MatchString(r.Ip) {
		return util.InvalidArgsError()
	}
	if r.Port <= 0 {
		return util.InvalidArgsError()
	}
	if !detectmd.IsNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateDetectReqDTO struct {
	Id       int64               `json:"id"`
	Ip       string              `json:"ip"`
	Port     int                 `json:"port"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateDetectReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !validIpPattern.MatchString(r.Ip) {
		return util.InvalidArgsError()
	}
	if r.Port <= 0 {
		return util.InvalidArgsError()
	}
	if !detectmd.IsNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDetectReqDTO struct {
	Name     string              `json:"name"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDetectReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit < 0 {
		return util.InvalidArgsError()
	}
	if len(r.Name) > 0 && !detectmd.IsNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteDetectReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteDetectReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DetectDTO struct {
	Id            int64  `json:"id"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Name          string `json:"name"`
	HeartbeatTime int64  `json:"heartbeatTime"`
	Enabled       bool   `json:"enabled"`
}

type ListLogReqDTO struct {
	Id       int64               `json:"id"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListLogReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 {
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

type LogDTO struct {
	Ip      string    `json:"ip"`
	Port    int       `json:"port"`
	Valid   bool      `json:"valid"`
	Created time.Time `json:"created"`
}

type EnableDetectReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EnableDetectReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisableDetectReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisableDetectReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
