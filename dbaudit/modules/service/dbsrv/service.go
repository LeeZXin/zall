package dbsrv

import "context"

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
	ListPermApprovalOrder(context.Context, ListPermApprovalOrderReqDTO) ([]ApprovalOrderDTO, int64, error)
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
}
