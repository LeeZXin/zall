package gitactionsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type InsertActionReqDTO struct {
	Id            int64               `json:"id"`
	ActionContent string              `json:"actionContent"`
	NodeId        int64               `json:"nodeId"`
	PushBranch    string              `json:"pushBranch"`
	Operator      apisession.UserInfo `json:"operator"`
}

func (r *InsertActionReqDTO) IsValid() error {
	if r.ActionContent == "" {
		return util.InvalidArgsError()
	}
	if r.NodeId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.PushBranch == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteActionReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteActionReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListActionReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListActionReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateActionReqDTO struct {
	Id            int64               `json:"id"`
	ActionContent string              `json:"actionContent"`
	NodeId        int64               `json:"nodeId"`
	PushBranch    string              `json:"pushBranch"`
	Operator      apisession.UserInfo `json:"operator"`
}

func (r *UpdateActionReqDTO) IsValid() error {
	if r.ActionContent == "" {
		return util.InvalidArgsError()
	}
	if r.NodeId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.PushBranch == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type TriggerActionReqDTO struct {
	Id       int64               `json:"id"`
	Args     map[string]string   `json:"args"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TriggerActionReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertNodeReqDTO struct {
	Name     string              `json:"name"`
	HttpHost string              `json:"httpHost"`
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
	return nil
}

type UpdateNodeReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	HttpHost string              `json:"httpHost"`
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
}
