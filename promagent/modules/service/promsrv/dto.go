package promsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
	"time"
)

type InsertScrapeReqDTO struct {
	ServerUrl  string              `json:"serverUrl"`
	AppId      string              `json:"appId"`
	Target     string              `json:"target"`
	TargetType prommd.TargetType   `json:"targetType"`
	Env        string              `json:"env"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *InsertScrapeReqDTO) IsValid() error {
	parsedUrl, err := url.Parse(r.ServerUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !prommd.IsTargetValid(r.Target) {
		return util.InvalidArgsError()
	}
	if !r.TargetType.IsValid() {
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

type UpdateScrapeReqDTO struct {
	Id         int64               `json:"id"`
	ServerUrl  string              `json:"serverUrl"`
	Target     string              `json:"target"`
	TargetType prommd.TargetType   `json:"targetType"`
	Env        string              `json:"env"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *UpdateScrapeReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	parsedUrl, err := url.Parse(r.ServerUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	if !prommd.IsTargetValid(r.Target) {
		return util.InvalidArgsError()
	}
	if !r.TargetType.IsValid() {
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

type ListScrapeReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListScrapeReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
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

type ScrapeDTO struct {
	Id         int64
	ServerUrl  string
	AppId      string
	Target     string
	TargetType prommd.TargetType
	Created    time.Time
}

type DeleteScrapeReqDTO struct {
	Id       int64               `json:"id"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteScrapeReqDTO) IsValid() error {
	if r.Id <= 0 {
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
