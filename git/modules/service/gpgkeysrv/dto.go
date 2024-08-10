package gpgkeysrv

import (
	"github.com/LeeZXin/zall/git/modules/model/gpgkeymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateGpgKeyReqDTO struct {
	Name     string              `json:"name"`
	Content  string              `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateGpgKeyReqDTO) IsValid() error {
	if !gpgkeymd.IsKeyNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.Content == "" {
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
	Id      int64
	Name    string
	KeyId   string
	Email   string
	Created time.Time
	Expired time.Time
	SubKeys string
}
