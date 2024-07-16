package promsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateScrapeReqDTO struct {
	Endpoint   string              `json:"endpoint"`
	AppId      string              `json:"appId"`
	Target     string              `json:"target"`
	TargetType prommd.TargetType   `json:"targetType"`
	Env        string              `json:"env"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *CreateScrapeReqDTO) IsValid() error {
	if !prommd.IsEndpointValid(r.Endpoint) {
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
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateScrapeReqDTO struct {
	ScrapeId   int64               `json:"scrapeId"`
	Endpoint   string              `json:"endpoint"`
	Target     string              `json:"target"`
	TargetType prommd.TargetType   `json:"targetType"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *UpdateScrapeReqDTO) IsValid() error {
	if r.ScrapeId <= 0 {
		return util.InvalidArgsError()
	}
	if !prommd.IsEndpointValid(r.Endpoint) {
		return util.InvalidArgsError()
	}
	if !prommd.IsTargetValid(r.Target) {
		return util.InvalidArgsError()
	}
	if !r.TargetType.IsValid() {
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
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type ScrapeDTO struct {
	Id         int64
	Endpoint   string
	AppId      string
	Target     string
	TargetType prommd.TargetType
	Created    time.Time
	Env        string
}

type DeleteScrapeReqDTO struct {
	ScrapeId int64               `json:"scrapeId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteScrapeReqDTO) IsValid() error {
	if r.ScrapeId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
