package gitnodesrv

import (
	"github.com/LeeZXin/zall/git/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type InsertNodeReqDTO struct {
	Name     string              `json:"name"`
	HttpHost string              `json:"httpHost"`
	SshHost  string              `json:"sshHost"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !gitnodemd.IsNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.HttpHost) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.SshHost) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateNodeReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	HttpHost string              `json:"httpHost"`
	SshHost  string              `json:"sshHost"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !gitnodemd.IsNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.HttpHost) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.SshHost) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteNodeReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListNodeReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type NodeDTO struct {
	Id       int64
	Name     string
	HttpHost string
	SshHost  string
}
