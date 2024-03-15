package approvalsrv

import (
	"github.com/LeeZXin/zall/approval/approval"
	"github.com/LeeZXin/zall/approval/modules/model/approvalmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type InsertProcessReqDTO struct {
	Pid      string        `json:"pid"`
	Name     string        `json:"name"`
	Approval approval.Node `json:"approval"`
}

func (r *InsertProcessReqDTO) IsValid() error {
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsProcessNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Approval.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateProcessReqDTO struct {
	Pid      string        `json:"pid"`
	Name     string        `json:"name"`
	Approval approval.Node `json:"approval"`
}

func (r *UpdateProcessReqDTO) IsValid() error {
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsProcessNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Approval.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteProcessReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteProcessReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetProcessReqDTO struct {
	Pid      string              `json:"pid"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetProcessReqDTO) IsValid() error {
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AgreeFlowReqDTO struct {
	NotifyId int64               `json:"notifyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AgreeFlowReqDTO) IsValid() error {
	if r.NotifyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisagreeFlowReqDTO struct {
	NotifyId int64               `json:"notifyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisagreeFlowReqDTO) IsValid() error {
	if r.NotifyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
