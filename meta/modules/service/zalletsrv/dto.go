package zalletsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/zalletmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type CreateZalletNodeReqDTO struct {
	NodeId     string              `json:"nodeId"`
	Name       string              `json:"name"`
	AgentHost  string              `json:"agentHost"`
	AgentToken string              `json:"agentToken"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *CreateZalletNodeReqDTO) IsValid() error {
	if !zalletmd.IsZalletNodeIdValid(r.NodeId) {
		return util.InvalidArgsError()
	}
	if !zalletmd.IsZalletNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !zalletmd.IsZalletAgentHostValid(r.AgentHost) {
		return util.InvalidArgsError()
	}
	if !zalletmd.IsZalletAgentTokenValid(r.AgentToken) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListZalletNodeReqDTO struct {
	Name     string              `json:"name"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListZalletNodeReqDTO) IsValid() error {
	if len(r.Name) > 0 && !zalletmd.IsZalletNameValid(r.Name) {
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

type DeleteZalletNodeReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteZalletNodeReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateZalletNodeReqDTO struct {
	Id         int64               `json:"id"`
	Name       string              `json:"name"`
	AgentHost  string              `json:"agentHost"`
	AgentToken string              `json:"agentToken"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *UpdateZalletNodeReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !zalletmd.IsZalletNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !zalletmd.IsZalletAgentHostValid(r.AgentHost) {
		return util.InvalidArgsError()
	}
	if !zalletmd.IsZalletAgentTokenValid(r.AgentToken) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ZalletNodeDTO struct {
	Id         int64  `json:"id"`
	NodeId     string `json:"nodeId"`
	Name       string `json:"name"`
	AgentHost  string `json:"agentHost"`
	AgentToken string `json:"agentToken"`
}

type ListAllZalletNodeReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllZalletNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SimpleZalletNodeDTO struct {
	Id     int64
	NodeId string
	Name   string
}
