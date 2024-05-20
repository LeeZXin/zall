package workflowsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/pkg/workflow"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	// TaskCallback 工作流回调
	TaskCallback(string, workflow.TaskStatus)
	// FindAndExecute 匹配仓库id 寻找并执行工作流
	FindAndExecute(int64, string, workflowmd.TriggerType, string, workflowmd.SourceType, int64)
	// Execute 执行工作流
	Execute(*workflowmd.Workflow, string, workflowmd.TriggerType, string, int64) error
}
type OuterService interface {
	// CreateWorkflow 创建工作流
	CreateWorkflow(context.Context, CreateWorkflowReqDTO) error
	// UpdateWorkflow 编辑工作流
	UpdateWorkflow(context.Context, UpdateWorkflowReqDTO) error
	// DeleteWorkflow 删除工作流
	DeleteWorkflow(context.Context, DeleteWorkflowReqDTO) error
	// ListWorkflowWithLastTask 工作流列表 + 最近执行任务
	ListWorkflowWithLastTask(context.Context, ListWorkflowWithLastTaskReqDTO) ([]WorkflowWithLastTaskDTO, error)
	// TriggerWorkflow 手动触发工作流
	TriggerWorkflow(context.Context, TriggerWorkflowReqDTO) error
	// ListTask 工作流任务列表
	ListTask(context.Context, ListTaskReqDTO) ([]TaskDTO, int64, error)
	ListStep(context.Context, ListStepReqDTO) ([]StepDTO, error)
	// GetWorkflowDetail 获取工作流详情
	GetWorkflowDetail(context.Context, GetWorkflowDetailReqDTO) (WorkflowDTO, error)
	// KillWorkflowTask 停止工作流
	KillWorkflowTask(context.Context, KillWorkflowTaskReqDTO) error
	// GetTaskDetail 获取工作流任务详情
	GetTaskDetail(context.Context, GetTaskDetailReqDTO) (TaskWithStepsDTO, error)
}
