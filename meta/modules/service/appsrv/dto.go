package appsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type AppDTO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type ListAppReqDTO struct {
	AppId    string              `json:"appId"`
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAppReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertAppReqDTO struct {
	AppId    string              `json:"appId"`
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertAppReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteAppReqDTO struct {
	AppId    string              `json:"appId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteAppReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateAppReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateAppReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TransferTeamReqDTO struct {
	AppId    string              `json:"appId"`
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TransferTeamReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
