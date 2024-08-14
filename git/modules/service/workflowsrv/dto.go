package workflowsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"time"
)

type CreateWorkflowReqDTO struct {
	Name        string              `json:"name"`
	RepoId      int64               `json:"repoId"`
	YamlContent string              `json:"yamlContent"`
	AgentId     int64               `json:"agentId"`
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
	if r.AgentId <= 0 {
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
	AgentId     int64               `json:"agentId"`
	Source      workflowmd.Source   `json:"source"`
	Desc        string              `json:"desc"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *UpdateWorkflowReqDTO) IsValid() error {
	if r.WorkflowId <= 0 {
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
	if r.AgentId <= 0 {
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

type TriggerWorkflowReqDTO struct {
	WorkflowId int64               `json:"workflowId"`
	Branch     string              `json:"branch"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *TriggerWorkflowReqDTO) IsValid() error {
	if r.WorkflowId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.Branch) == 0 || len(r.Branch) > 1024 {
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

type ListTaskByPrIdReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTaskByPrIdReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TaskWithoutYamlContentDTO struct {
	Id          int64
	TaskStatus  sshagent.Status
	TriggerType workflowmd.TriggerType
	Branch      string
	PrId        int64
	Operator    string
	Created     time.Time
	Duration    int64
	WorkflowId  int64
}

type TaskDTO struct {
	TaskWithoutYamlContentDTO
	YamlContent string
}

type WorkflowTaskDTO struct {
	Name string
	TaskWithoutYamlContentDTO
}

type WorkflowWithLastTaskDTO struct {
	Id       int64
	Name     string
	Desc     string
	LastTask *TaskWithoutYamlContentDTO
}

type WorkflowDTO struct {
	Id          int64             `json:"id"`
	Name        string            `json:"name"`
	Desc        string            `json:"desc"`
	RepoId      int64             `json:"repoId"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	AgentId     int64             `json:"agentId"`
}

type GetWorkflowDetailReqDTO struct {
	WorkflowId int64               `json:"workflowId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *GetWorkflowDetailReqDTO) IsValid() error {
	if r.WorkflowId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type KillWorkflowTaskReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *KillWorkflowTaskReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetTaskDetailReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetTaskDetailReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetTaskStatusReqDTO struct {
	TaskId   int64               `json:"taskId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetTaskStatusReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetLogContentReqDTO struct {
	TaskId    int64               `json:"taskId"`
	JobName   string              `json:"jobName"`
	StepIndex int                 `json:"stepIndex"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *GetLogContentReqDTO) IsValid() error {
	if r.TaskId <= 0 {
		return util.InvalidArgsError()
	}
	if r.JobName == "" {
		return util.InvalidArgsError()
	}
	if r.StepIndex < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListVarsReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListVarsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreateVarsReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Name     string              `json:"name"`
	Content  string              `json:"content"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateVarsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !workflowmd.IsVarsNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !workflowmd.IsVarsContentValid(r.Content) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateVarsReqDTO struct {
	VarsId   int64               `json:"varsId"`
	Content  string              `json:"content"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateVarsReqDTO) IsValid() error {
	if r.VarsId <= 0 {
		return util.InvalidArgsError()
	}
	if !workflowmd.IsVarsContentValid(r.Content) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteVarsReqDTO struct {
	VarsId   int64               `json:"varsId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteVarsReqDTO) IsValid() error {
	if r.VarsId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetVarsContentReqDTO struct {
	VarsId   int64               `json:"varsId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetVarsContentReqDTO) IsValid() error {
	if r.VarsId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type VarsWithoutContentDTO struct {
	VarsId int64
	Name   string
}

type VarsDTO struct {
	VarsWithoutContentDTO
	Content string
}

type FindAndExecuteWorkflowReqDTO struct {
	RepoId      int64
	RepoPath    string
	Operator    string
	TriggerType workflowmd.TriggerType
	Branch      string
	Source      workflowmd.SourceType
	PrId        int64
}

type ExecuteWorkflowReqDTO struct {
	RepoPath    string
	Operator    string
	TriggerType workflowmd.TriggerType
	Branch      string
	PrId        int64
}
