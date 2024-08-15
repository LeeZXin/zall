package appsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
)

type AppDTO struct {
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type AppWithPermDTO struct {
	AppDTO
	Perm perm.AppPerm `json:"perm"`
}

type ListAppReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAppReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAllAppByAdminReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllAppByAdminReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAllAppBySaReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllAppBySaReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreateAppReqDTO struct {
	AppId    string              `json:"appId"`
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateAppReqDTO) IsValid() error {
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

type GetAppWithPermReqDTO struct {
	AppId    string              `json:"appId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetAppWithPermReqDTO) IsValid() error {
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
