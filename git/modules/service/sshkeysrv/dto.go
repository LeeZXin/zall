package sshkeysrv

import (
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateSshKeyReqDTO struct {
	Name     string              `json:"name"`
	Content  string              `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateSshKeyReqDTO) IsValid() error {
	if !sshkeymd.IsSshKeyNameValid(r.Name) {
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

type DeleteSshKeyReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteSshKeyReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
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

type SshKeyDTO struct {
	Id           int64
	Name         string
	Fingerprint  string
	Created      time.Time
	LastOperated time.Time
}
