package mysqldbsrv

import (
	"context"
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv/command"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertDb(context.Context, InsertDbReqDTO) error
	UpdateDb(context.Context, UpdateDbReqDTO) error
	DeleteDb(context.Context, DeleteDbReqDTO) error
	ListDb(context.Context, ListDbReqDTO) ([]DbDTO, error)
	ListSimpleDb(context.Context, ListSimpleDbReqDTO) ([]SimpleDbDTO, error)
	// ApplyDbPerm 申请库表权限
	ApplyDbPerm(context.Context, ApplyDbPermReqDTO) error
	// ListPermApprovalOrder 展示审批列表
	ListPermApprovalOrder(context.Context, ListPermApprovalOrderReqDTO) ([]PermApprovalOrderDTO, int64, error)
	// ListAppliedPermApprovalOrder 展示申请的审批列表
	ListAppliedPermApprovalOrder(context.Context, ListAppliedPermApprovalOrderReqDTO) ([]PermApprovalOrderDTO, int64, error)
	// AgreeDbPerm 同意审批
	AgreeDbPerm(context.Context, AgreeDbPermReqDTO) error
	// DisagreeDbPerm 不同意审批
	DisagreeDbPerm(context.Context, DisagreeDbPermReqDTO) error
	// CancelDbPerm 取消申请
	CancelDbPerm(context.Context, CancelDbPermReqDTO) error
	// ListDbPerm 权限列表
	ListDbPerm(context.Context, ListDbPermReqDTO) ([]PermDTO, int64, error)
	// DeleteDbPerm 删除权限
	DeleteDbPerm(context.Context, DeleteDbPermReqDTO) error
	// ListDbPermByAccount 权限列表
	ListDbPermByAccount(context.Context, ListDbPermByAccountReqDTO) ([]PermDTO, int64, error)
	// AllBases 所有库
	AllBases(context.Context, AllBasesReqDTO) ([]string, error)
	// AllTables 展示单个数据库所有表
	AllTables(context.Context, AllTablesReqDTO) ([]string, error)
	// SearchDb 搜索
	SearchDb(context.Context, SearchDbReqDTO) ([]string, [][]string, error)
	// ApplyDbUpdate 提数据库修改单
	ApplyDbUpdate(context.Context, ApplyDbUpdateReqDTO) ([]command.ValidateUpdateResult, bool, error)
	// ListUpdateApprovalOrder 数据库修改审批单
	ListUpdateApprovalOrder(context.Context, ListUpdateApprovalOrderReqDTO) ([]UpdateApprovalOrderDTO, int64, error)
	// ListAppliedUpdateApprovalOrder 申请的数据库修改审批单
	ListAppliedUpdateApprovalOrder(context.Context, ListAppliedUpdateApprovalOrderReqDTO) ([]UpdateApprovalOrderDTO, int64, error)
	// AgreeDbUpdate 同意修改单
	AgreeDbUpdate(context.Context, AgreeDbUpdateReqDTO) error
	// DisagreeDbUpdate 不同意修改单
	DisagreeDbUpdate(context.Context, DisagreeDbUpdateReqDTO) error
	// CancelDbUpdate 取消修改单
	CancelDbUpdate(context.Context, CancelDbUpdateReqDTO) error
	// ExecuteDbUpdate 执行修改单
	ExecuteDbUpdate(context.Context, ExecuteDbUpdateReqDTO) error
	// AskToExecuteDbUpdate 请求执行修改单
	AskToExecuteDbUpdate(context.Context, AskToExecuteDbUpdateReqDTO) error
}
