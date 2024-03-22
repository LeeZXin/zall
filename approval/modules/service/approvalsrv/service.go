package approvalsrv

import "context"

var (
	Inner InnerService = new(innerImpl)
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	InsertFlow(context.Context, InsertFlowReqDTO) error
	// InsertAttachedProcess 创建系统审批流
	InsertAttachedProcess(context.Context, InsertAttachedProcessReqDTO) error
	// UpdateAttachedProcess 编辑系统审批流
	UpdateAttachedProcess(context.Context, UpdateAttachedProcessReqDTO) error
	// DeleteAttachedProcess 删除系统审批流
	DeleteAttachedProcess(context.Context, string) error
}

type OuterService interface {
	// Agree 同意
	Agree(context.Context, AgreeFlowReqDTO) error
	// Disagree 不同意
	Disagree(context.Context, DisagreeFlowReqDTO) error
	// InsertCustomProcess 创建自定义审批流
	InsertCustomProcess(context.Context, InsertCustomProcessReqDTO) error
	// UpdateCustomProcess 编辑自定义审批流
	UpdateCustomProcess(context.Context, UpdateCustomProcessReqDTO) error
	// DeleteCustomProcess 删除自定义审批流
	DeleteCustomProcess(context.Context, DeleteCustomProcessReqDTO) error
	// ListCustomProcess 展示详细审批流列表
	ListCustomProcess(context.Context, ListCustomProcessReqDTO) ([]ProcessDTO, error)
	// ListAllGroupProcess 自定义审批流列表
	ListAllGroupProcess(context.Context, ListAllGroupProcessReqDTO) ([]GroupProcessDTO, error)
	// InsertCustomFlow 发起自定义审批流
	InsertCustomFlow(context.Context, InsertCustomFlowReqDTO) ([]string, error)
	// CancelCustomFlow 取消自定义审批流
	CancelCustomFlow(context.Context, CancelCustomFlowReqDTO) error
	// ListCustomFlow 获取申请的审批流
	ListCustomFlow(context.Context, ListCustomFlowReqDTO) ([]FlowDTO, error)
	// ListOperateFlow 获取操作的审批流
	ListOperateFlow(context.Context, ListOperateFlowReqDTO) ([]FlowDTO, error)
	// GetFlowDetail 获取审批流详情
	GetFlowDetail(context.Context, GetFlowDetailReqDTO) (FlowDetailDTO, error)
	// InsertGroup 插入审批流分组
	InsertGroup(context.Context, InsertGroupReqDTO) error
	// DeleteGroup 删除分组
	DeleteGroup(context.Context, DeleteGroupReqDTO) error
	// ListGroup 展示分组
	ListGroup(context.Context, ListGroupReqDTO) ([]GroupDTO, error)
	// UpdateGroup 修改分组
	UpdateGroup(context.Context, UpdateGroupReqDTO) error
}
