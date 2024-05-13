package workflowapi

import (
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
)

type CreateWorkflowReqVO struct {
	Name        string            `json:"name"`
	RepoId      int64             `json:"repoId"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	Agent       zssh.AgentCfg     `json:"agent"`
	Desc        string            `json:"desc"`
}

type UpdateWorkflowReqVO struct {
	WorkflowId  int64             `json:"workflowId"`
	Name        string            `json:"name"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	Agent       zssh.AgentCfg     `json:"agent"`
}

type WorkflowWithLastTaskVO struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	LastTask *TaskVO `json:"lastTask,omitempty"`
}

type ListTaskReqVO struct {
	WorkflowId int64 `json:"workflowId"`
	Cursor     int64 `json:"cursor"`
	Limit      int   `json:"limit"`
}

type TaskVO struct {
	Id          int64                  `json:"id"`
	TaskStatus  workflowmd.TaskStatus  `json:"taskStatus"`
	TriggerType workflowmd.TriggerType `json:"triggerType"`
	YamlContent string                 `json:"yamlContent"`
	Operator    string                 `json:"operator"`
	Created     string                 `json:"created"`
	Branch      string                 `json:"branch"`
}

type StepVO struct {
	JobName    string `json:"jobName"`
	StepName   string `json:"stepName"`
	StepIndex  int    `json:"stepIndex"`
	LogContent string `json:"logContent"`
	StepStatus string `json:"stepStatus"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}