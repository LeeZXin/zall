package promsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
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
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateScrapeReqDTO struct {
	Id         int64               `json:"id"`
	Endpoint   string              `json:"endpoint"`
	Target     string              `json:"target"`
	TargetType prommd.TargetType   `json:"targetType"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *UpdateScrapeReqDTO) IsValid() error {
	if r.Id <= 0 {
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
	Endpoint string              `json:"endpoint"`
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListScrapeReqDTO) IsValidBySa() error {
	if len(r.Endpoint) > 0 && !prommd.IsEndpointValid(r.Endpoint) {
		return util.InvalidArgsError()
	}
	if len(r.AppId) > 0 && !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

func (r *ListScrapeReqDTO) IsValidByTeam() error {
	if len(r.Endpoint) > 0 && !prommd.IsEndpointValid(r.Endpoint) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}

type ScrapeByTeamDTO struct {
	Id         int64
	Endpoint   string
	Target     string
	TargetType prommd.TargetType
	Env        string
}

type ScrapeBySaDTO struct {
	Id         int64
	Endpoint   string
	AppId      string
	AppName    string
	TeamId     int64
	TeamName   string
	Target     string
	TargetType prommd.TargetType
	Env        string
}

type DeleteScrapeReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteScrapeReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
