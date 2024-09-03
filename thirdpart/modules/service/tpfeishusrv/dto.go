package tpfeishusrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpfeishumd"
	"github.com/LeeZXin/zall/util"
	"time"
)

type ListAccessTokenReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Key      string              `json:"key"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAccessTokenReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AccessTokenDTO struct {
	Id          int64
	TeamId      int64
	Name        string
	AppId       string
	Creator     string
	Secret      string
	Token       string
	TenantToken string
	ApiKey      string
	Expired     time.Time
}

type CreateAccessTokenReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	AppId    string              `json:"appId"`
	Secret   string              `json:"secret"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateAccessTokenReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !tpfeishumd.IsAccessTokenNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !tpfeishumd.IsAccessTokenAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !tpfeishumd.IsAccessTokenSecretValid(r.Secret) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateAccessTokenReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	AppId    string              `json:"appId"`
	Secret   string              `json:"secret"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateAccessTokenReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !tpfeishumd.IsAccessTokenNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !tpfeishumd.IsAccessTokenAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !tpfeishumd.IsAccessTokenSecretValid(r.Secret) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteAccessTokenReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteAccessTokenReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RefreshAccessTokenReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RefreshAccessTokenReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ChangeAccessTokenApiKeyReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ChangeAccessTokenApiKeyReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
