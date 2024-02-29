package sshkeysrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type InsertSshKeyReqDTO struct {
	Name          string              `json:"name"`
	PubKeyContent string              `json:"-"`
	Signature     string              `json:"-"`
	Operator      apisession.UserInfo `json:"operator"`
}

func (r *InsertSshKeyReqDTO) IsValid() error {
	if len(r.Name) == 0 || len(r.Name) > 128 {
		return util.InvalidArgsError()
	}
	if r.PubKeyContent == "" || r.Signature == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteSshKeyReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteSshKeyReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListSshKeyReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListSshKeyReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetTokenReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetTokenReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
