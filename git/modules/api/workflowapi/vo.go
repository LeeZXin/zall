package workflowapi

import (
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
)

type CreateWorkflowReqVO struct {
	Name        string            `json:"name"`
	RepoId      int64             `json:"repoId"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	AgentHost   string            `json:"agentHost"`
	AgentToken  string            `json:"agentToken"`
	Desc        string            `json:"desc"`
}

type UpdateWorkflowReqVO struct {
	WorkflowId  int64             `json:"workflowId"`
	Name        string            `json:"name"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	AgentHost   string            `json:"agentHost"`
	AgentToken  string            `json:"agentToken"`
	Desc        string            `json:"desc"`
}

type WorkflowWithLastTaskVO struct {
	Id       int64                     `json:"id"`
	Name     string                    `json:"name"`
	Desc     string                    `json:"desc"`
	LastTask *TaskWithoutYamlContentVO `json:"lastTask,omitempty"`
}

type TaskWithoutYamlContentVO struct {
	Id          int64                  `json:"id"`
	TaskStatus  workflowmd.TaskStatus  `json:"taskStatus"`
	TriggerType workflowmd.TriggerType `json:"triggerType"`
	Operator    string                 `json:"operator"`
	Created     string                 `json:"created"`
	Branch      string                 `json:"branch"`
	PrId        int64                  `json:"prId"`
	Duration    int64                  `json:"duration"`
	WorkflowId  int64                  `json:"workflowId"`
}

type TaskVO struct {
	TaskWithoutYamlContentVO
	YamlContent string `json:"yamlContent,omitempty"`
}

type WorkflowTaskVO struct {
	Name string `json:"name"`
	TaskWithoutYamlContentVO
}

type WorkflowVO struct {
	Id          int64             `json:"id"`
	Name        string            `json:"name"`
	Desc        string            `json:"desc"`
	RepoId      int64             `json:"repoId"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	AgentHost   string            `json:"agentHost"`
	AgentToken  string            `json:"agentToken"`
}
