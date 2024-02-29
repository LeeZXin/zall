package gpgkeysrv

import (
	"github.com/LeeZXin/zall/git/modules/model/gpgkeymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type InsertGpgKeyReqDTO struct {
	Name      string              `json:"name"`
	Content   string              `json:"-"`
	Signature string              `json:"-"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *InsertGpgKeyReqDTO) IsValid() error {
	if !gpgkeymd.IsKeyNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.Signature == "" || r.Content == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteGpgKeyReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteGpgKeyReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetTokenReqDTO struct {
	Content  string              `json:"content"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetTokenReqDTO) IsValid() error {
	if r.Content == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListGpgKeyReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListGpgKeyReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GpgKeyDTO struct {
	Id         int64     `json:"keyId"`
	Name       string    `json:"name"`
	PubKeyId   string    `json:"pubKeyId"`
	EmailList  []string  `json:"emailList"`
	ExpireTime time.Time `json:"expireTime"`
}
