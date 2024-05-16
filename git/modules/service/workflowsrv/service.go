package workflowsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	// FindAndExecute 匹配仓库id 寻找并执行工作流
	FindAndExecute(int64, string, workflowmd.TriggerType, string, workflowmd.SourceType, int64)
	// Execute 执行工作流
	Execute(*workflowmd.Workflow, string, workflowmd.TriggerType, string, int64) error
}
type OuterService interface {
	// CreateWorkflow 创建工作流
	CreateWorkflow(context.Context, CreateWorkflowReqDTO) error
	UpdateWorkflow(context.Context, UpdateWorkflowReqDTO) error
	// DeleteWorkflow 删除工作流
	DeleteWorkflow(context.Context, DeleteWorkflowReqDTO) error
	// ListWorkflowWithLastTask 工作流列表 + 最近执行任务
	ListWorkflowWithLastTask(context.Context, ListWorkflowWithLastTaskReqDTO) ([]WorkflowWithLastTaskDTO, error)
	TriggerWorkflow(context.Context, TriggerWorkflowReqDTO) error
	// ListTask 工作流任务列表
	ListTask(context.Context, ListTaskReqDTO) ([]TaskDTO, int64, error)
	ListStep(context.Context, ListStepReqDTO) ([]StepDTO, error)
	// GetDetail 获取工作流详情

}
