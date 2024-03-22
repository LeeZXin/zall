package approvalsrv

import (
	"github.com/LeeZXin/zall/approval/modules/model/approvalmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	dayTimePattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
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
	IconUrl  string              `json:"iconUrl"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertCustomProcessReqDTO) IsValid() error {
	parsedUrl, err := url.Parse(r.IconUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	if r.GroupId <= 0 {
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
	IconUrl  string              `json:"iconUrl"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateCustomProcessReqDTO) IsValid() error {
	if r.GroupId <= 0 {
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
	parsedUrl, err := url.Parse(r.IconUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
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

type ListOperateFlowReqDTO struct {
	DayTime  string              `json:"dayTime"`
	Done     bool                `json:"done"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListOperateFlowReqDTO) IsValid() error {
	if !dayTimePattern.MatchString(r.DayTime) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDetailReqDTO struct {
	DayTime  string              `json:"dayTime"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDetailReqDTO) IsValid() error {
	if !dayTimePattern.MatchString(r.DayTime) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListCustomFlowReqDTO struct {
	DayTime  string              `json:"dayTime"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListCustomFlowReqDTO) IsValid() error {
	if !dayTimePattern.MatchString(r.DayTime) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type FlowDTO struct {
	Id          int64                 `json:"flowId"`
	ProcessName string                `json:"processName"`
	FlowStatus  approvalmd.FlowStatus `json:"flowStatus"`
	Creator     string                `json:"creator"`
	Created     time.Time             `json:"created"`
}

type ListAllGroupProcessReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllGroupProcessReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GroupProcessDTO struct {
	Id        int64
	Name      string
	Processes []SimpleProcessDTO
}

type SimpleProcessDTO struct {
	Id      int64
	Name    string
	IconUrl string
}

type ProcessDTO struct {
	Id      int64
	Pid     string
	GroupId int64
	Name    string
	Content approval.Process
	IconUrl string
	Created time.Time
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

type ListCustomProcessReqDTO struct {
	GroupId  int64               `json:"groupId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListCustomProcessReqDTO) IsValid() error {
	if r.GroupId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetFlowDetailReqDTO struct {
	FlowId   int64               `json:"flowId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetFlowDetailReqDTO) IsValid() error {
	if r.FlowId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type NotifyDTO struct {
	Account   string            `json:"account"`
	FlowIndex int               `json:"flowIndex"`
	Done      bool              `json:"done"`
	Op        approvalmd.FlowOp `json:"op"`
	Updated   time.Time         `json:"updated" xorm:"updated"`
}

type FlowDetailDTO struct {
	Id          int64                 `json:"flowId"`
	ProcessName string                `json:"processName"`
	FlowStatus  approvalmd.FlowStatus `json:"flowStatus"`
	Creator     string                `json:"creator"`
	Created     time.Time             `json:"created"`
	Kvs         []approval.Kv         `json:"kvs"`
	Process     approval.Process      `json:"process"`
	NotifyList  []NotifyDTO           `json:"notifyList"`
}

type InsertGroupReqDTO struct {
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertGroupReqDTO) IsValid() error {
	if !approvalmd.IsGroupNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteGroupReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteGroupReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListGroupReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListGroupReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GroupDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateGroupReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateGroupReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !approvalmd.IsGroupNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
