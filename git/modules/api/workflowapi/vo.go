package workflowapi

import (
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/util"
)

type CreateWorkflowReqVO struct {
	Name        string            `json:"name"`
	RepoId      int64             `json:"repoId"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	AgentId     string            `json:"agentId"`
	Desc        string            `json:"desc"`
}

type UpdateWorkflowReqVO struct {
	WorkflowId  int64             `json:"workflowId"`
	Name        string            `json:"name"`
	YamlContent string            `json:"yamlContent"`
	Source      workflowmd.Source `json:"source"`
	AgentId     string            `json:"agentId"`
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
	TaskStatus  sshagent.Status        `json:"taskStatus"`
	TriggerType workflowmd.TriggerType `json:"triggerType"`
	Operator    util.User              `json:"operator"`
	Created     string                 `json:"created"`
	Branch      string                 `json:"branch"`
	PrId        int64                  `json:"prId"`
	PrIndex     int                    `json:"prIndex"`
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
	AgentId     string            `json:"agentId"`
}

type CreateVarsReqVO struct {
	RepoId  int64  `json:"repoId"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type UpdateVarsReqVO struct {
	VarsId  int64  `json:"varsId"`
	Content string `json:"content"`
}

type VarsWithoutContentVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type VarsVO struct {
	VarsWithoutContentVO
	Content string `json:"content"`
}

type ListTaskReqVO struct {
	WorkflowId int64 `json:"workflowId"`
	PageNum    int   `json:"pageNum"`
}
