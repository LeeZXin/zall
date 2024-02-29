package actiontaskmd

import "github.com/LeeZXin/zall/pkg/action"

type InsertTaskReqDTO struct {
	RepoId      int64
	TaskName    string
	InstanceId  string
	TaskType    TaskType
	TaskStatus  TaskStatus
	TriggerType TriggerType
	Hook        action.Webhook
}

type InsertStepReqDTO struct {
	TaskId     int64      `json:"taskId"`
	JobName    string     `json:"jobName"`
	StepName   string     `json:"stepName"`
	StepIndex  int        `json:"stepIndex"`
	StepStatus StepStatus `json:"stepStatus"`
}
