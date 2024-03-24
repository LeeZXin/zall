package actionsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/actionmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
	"time"
)

type InsertActionReqDTO struct {
	Name          string              `json:"name"`
	TeamId        int64               `json:"teamId"`
	ActionContent string              `json:"actionContent"`
	NodeId        int64               `json:"nodeId"`
	Operator      apisession.UserInfo `json:"operator"`
}

func (r *InsertActionReqDTO) IsValid() error {
	if !actionmd.IsActionNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.ActionContent == "" {
		return util.InvalidArgsError()
	}
	if r.NodeId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
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
	TeamId   int64               `json:"teamId"`
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
	return nil
}

type TriggerActionReqDTO struct {
	Id       int64               `json:"id"`
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
	Token    string              `json:"token"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !actionmd.IsNodeNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !actionmd.IsNodeTokenValid(r.Token) {
		return util.InvalidArgsError()
	}
	parsedUrl, err := url.Parse(r.HttpHost)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateNodeReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	HttpHost string              `json:"httpHost"`
	Token    string              `json:"token"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !actionmd.IsNodeNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !actionmd.IsNodeTokenValid(r.Token) {
		return util.InvalidArgsError()
	}
	parsedUrl, err := url.Parse(r.HttpHost)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
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

type AllNodeReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AllNodeReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type NodeDTO struct {
	Id       int64
	Name     string
	HttpHost string
	Token    string
}

type SimpleNodeDTO struct {
	Id   int64
	Name string
}

type ListTaskReqDTO struct {
	ActionId int64               `json:"actionId"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTaskReqDTO) IsValid() error {
	if r.ActionId <= 0 || r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListStepReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListStepReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TaskDTO struct {
	TaskStatus    actionmd.TaskStatus
	TriggerType   actionmd.TriggerType
	ActionContent string
	Operator      string
	Created       time.Time
}

type StepDTO struct {
	JobName    string
	StepName   string
	StepIndex  int
	LogContent string
	StepStatus actionmd.StepStatus
	Created    time.Time
	Updated    time.Time
}
