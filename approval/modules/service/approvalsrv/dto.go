package approvalsrv

import (
	"github.com/LeeZXin/zall/approval/modules/model/approvalmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zall/util"
)

type InsertAttachedProcessReqDTO struct {
	Pid       string              `json:"pid"`
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Process   approval.ProcessCfg `json:"process"`
}

func (r *InsertAttachedProcessReqDTO) IsValid() error {
	if r.Namespace == "" {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsProcessNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Process.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateAttachedProcessReqDTO struct {
	Pid     string              `json:"pid"`
	Name    string              `json:"name"`
	Process approval.ProcessCfg `json:"process"`
}

func (r *UpdateAttachedProcessReqDTO) IsValid() error {
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsProcessNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Process.IsValid() {
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

type InsertCustomProcessReqDTO struct {
	Pid      string              `json:"pid"`
	Name     string              `json:"name"`
	GroupId  int64               `json:"groupId"`
	Process  approval.ProcessCfg `json:"process"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertCustomProcessReqDTO) IsValid() error {
	// 不能等于1 1是default
	if r.GroupId <= 1 {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsProcessNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Process.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateCustomProcessReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	GroupId  int64               `json:"groupId"`
	Process  approval.ProcessCfg `json:"process"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateCustomProcessReqDTO) IsValid() error {
	// 不能等于1 1是default
	if r.GroupId <= 1 {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsProcessNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Process.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertFlowReqDTO struct {
	Pid     string
	Account string
	BizId   string
	Kvs     []approval.Kv
}

func (r *InsertFlowReqDTO) IsValid() error {
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertCustomFlowReqDTO struct {
	Pid      string              `json:"pid"`
	BizId    string              `json:"bizId"`
	Kvs      []approval.Kv       `json:"kvs"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertCustomFlowReqDTO) IsValid() error {
	if !approvalmd.IsPidValid(r.Pid) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CancelCustomFlowReqDTO struct {
	FlowId   int64               `json:"flowId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CancelCustomFlowReqDTO) IsValid() error {
	if r.FlowId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListUnDoneNotifyReqDTO struct {
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListUnDoneNotifyReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UnDoneNotifyDTO struct {
}

type ListDetailReqDTO struct {
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDetailReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DetailDTO struct {
}

type ListCustomFlowReqDTO struct {
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListCustomFlowReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type FlowDTO struct {
}

type ListCustomProcessReqDTO struct {
	Pid      string              `json:"pid"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	GroupId  int64               `json:"groupId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListCustomProcessReqDTO) IsValid() error {
	// 不能展示系统默认group审批流
	if r.GroupId == 1 {
		return util.InvalidArgsError()
	}
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ProcessDTO struct {
}

type DeleteCustomProcessReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteCustomProcessReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
