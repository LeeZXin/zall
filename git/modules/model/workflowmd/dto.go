package workflowmd

import (
	zssh "github.com/LeeZXin/zall/pkg/ssh"
)

type InsertTaskReqDTO struct {
	WorkflowId  int64
	TaskStatus  TaskStatus
	TriggerType TriggerType
	Operator    string
	Branch      string
	Workflow    WorkflowCfg
}

type InsertStepReqDTO struct {
	WorkflowId int64
	TaskId     int64
	JobName    string
	StepName   string
	StepIndex  int
	StepStatus StepStatus
}

type InsertWorkflowReqDTO struct {
	RepoId      int64
	Name        string
	YamlContent string
	Agent       zssh.AgentCfg
	Source      Source
	Desc        string
}

type UpdateWorkflowReqDTO struct {
	Id      int64
	Name    string
	Content string
	Agent   zssh.AgentCfg
}

type ListTaskByWorkflowIdReqDTO struct {
	WorkflowId int64
	PageNum    int
	PageSize   int
}
