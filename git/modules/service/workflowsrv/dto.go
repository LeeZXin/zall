package workflowsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"time"
)

type CreateWorkflowReqDTO struct {
	Name        string              `json:"name"`
	RepoId      int64               `json:"repoId"`
	YamlContent string              `json:"yamlContent"`
	Agent       zssh.AgentCfg       `json:"agent"`
	Source      workflowmd.Source   `json:"source"`
	Desc        string              `json:"desc"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *CreateWorkflowReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !workflowmd.IsWorkflowNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !workflowmd.IsWorkflowDescValid(r.Desc) {
		return util.InvalidArgsError()
	}
	if r.YamlContent == "" {
		return util.InvalidArgsError()
	}
	if !r.Agent.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Source.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteWorkflowReqDTO struct {
	WorkflowId int64               `json:"workflowId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *DeleteWorkflowReqDTO) IsValid() error {
	if r.WorkflowId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListWorkflowWithLastTaskReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListWorkflowWithLastTaskReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateWorkflowReqDTO struct {
	WorkflowId  int64               `json:"workflowId"`
	Name        string              `json:"name"`
	YamlContent string              `json:"yamlContent"`
	Agent       zssh.AgentCfg       `json:"agent"`
	Source      workflowmd.Source   `json:"source"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *UpdateWorkflowReqDTO) IsValid() error {
	if r.WorkflowId <= 0 {
		return util.InvalidArgsError()
	}
	if !workflowmd.IsWorkflowNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.YamlContent == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Agent.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Source.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TriggerWorkflowReqDTO struct {
	WorkflowId int64               `json:"workflowId"`
	Branch     string              `json:"branch"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *TriggerWorkflowReqDTO) IsValid() error {
	if r.WorkflowId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Branch == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTaskReqDTO struct {
	WorkflowId int64 `json:"workflowId"`
	ginutil.Page2Req
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTaskReqDTO) IsValid() error {
	if r.WorkflowId <= 0 || r.PageNum <= 0 || r.PageSize <= 0 || r.PageSize > 1000 {
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
	Id          int64
	TaskStatus  workflowmd.TaskStatus
	TriggerType workflowmd.TriggerType
	YamlContent string
	Branch      string
	Operator    string
	Created     time.Time
}

type StepDTO struct {
	JobName    string
	StepName   string
	StepIndex  int
	LogContent string
	StepStatus workflowmd.StepStatus
	Created    time.Time
	Updated    time.Time
}

type WorkflowWithLastTaskDTO struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	LastTask *TaskDTO `json:"lastTask"`
}

type WorkflowDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}