package gitnodesrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type InsertNodeReqDTO struct {
	NodeId    string              `json:"nodeId"`
	HttpHosts []string            `json:"httpHosts"`
	SshHosts  []string            `json:"sshHosts"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *InsertNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !gitnodemd.IsNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if len(r.HttpHosts) == 0 || len(r.SshHosts) == 0 {
		return util.InvalidArgsError()
	}
	for _, host := range r.HttpHosts {
		if !util.IpPortPattern.MatchString(host) {
			return util.InvalidArgsError()
		}
	}
	for _, host := range r.SshHosts {
		if !util.IpPortPattern.MatchString(host) {
			return util.InvalidArgsError()
		}
	}
	return nil
}

type UpdateNodeReqDTO struct {
	NodeId    string              `json:"nodeId"`
	HttpHosts []string            `json:"httpHosts"`
	SshHosts  []string            `json:"sshHosts"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !gitnodemd.IsNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if len(r.HttpHosts) == 0 || len(r.SshHosts) == 0 {
		return util.InvalidArgsError()
	}
	for _, host := range r.HttpHosts {
		if !util.IpPortPattern.MatchString(host) {
			return util.InvalidArgsError()
		}
	}
	for _, host := range r.SshHosts {
		if !util.IpPortPattern.MatchString(host) {
			return util.InvalidArgsError()
		}
	}
	return nil
}

type DeleteNodeReqDTO struct {
	NodeId   string              `json:"nodeId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !gitnodemd.IsNodeIdValid(r.NodeId) {
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

type GitNodeDTO struct {
	NodeId    string
	HttpHosts []string
	SshHosts  []string
}
