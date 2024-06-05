package workflowmd

import "time"

type InsertTaskReqDTO struct {
	WorkflowId  int64
	TaskStatus  TaskStatus
	TriggerType TriggerType
	YamlContent string
	Operator    string
	Branch      string
	PrId        int64
	AgentHost   string
	AgentToken  string
	BizId       string
}

type InsertWorkflowReqDTO struct {
	RepoId      int64
	Name        string
	YamlContent string
	AgentHost   string
	AgentToken  string
	Source      Source
	Desc        string
}

type UpdateWorkflowReqDTO struct {
	Id         int64
	Name       string
	Content    string
	AgentHost  string
	AgentToken string
	Desc       string
	Source     Source
}

type ListTaskByWorkflowIdReqDTO struct {
	WorkflowId int64
	PageNum    int
	PageSize   int
}

type InsertVarsReqDTO struct {
	RepoId  int64
	Name    string
	Content string
}

type UpdateVarsReqDTO struct {
	Id      int64
	Content string
}

type ExistsVarsReqDTO struct {
	RepoId int64
	Name   string
}

type InsertTokenReqDTO struct {
	RepoId   int64
	TaskId   int64
	Content  string
	Expired  time.Time
	Operator string
}
