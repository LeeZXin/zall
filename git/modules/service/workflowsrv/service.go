package workflowsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/pkg/workflow"
)

var (
	Outer = newOuterService()

	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	// TaskCallback 工作流回调
	TaskCallback(string, workflow.TaskStatusCallbackReq)
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
	ListTask(context.Context, ListTaskReqDTO) ([]TaskWithoutYamlContentDTO, int64, error)
	// ListTaskByPrId 合并请求相关工作流任务列表
	ListTaskByPrId(context.Context, ListTaskByPrIdReqDTO) ([]WorkflowTaskDTO, error)
	// GetTaskDetail 获取任务详情
	GetTaskDetail(context.Context, GetTaskDetailReqDTO) (TaskDTO, error)
	// GetWorkflowDetail 获取工作流详情
	GetWorkflowDetail(context.Context, GetWorkflowDetailReqDTO) (WorkflowDTO, error)
	// KillWorkflowTask 停止工作流
	KillWorkflowTask(context.Context, KillWorkflowTaskReqDTO) error
	// GetTaskStatus 获取任务状态
	GetTaskStatus(context.Context, GetTaskStatusReqDTO) (workflow.TaskStatus, error)
	// GetLogContent 获取日志内容
	GetLogContent(context.Context, GetLogContentReqDTO) ([]string, error)
}
