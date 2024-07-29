package mysqldbsrv

import (
	"context"
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv/command"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
)

var (
	Outer        OuterService
	defaultBases *hashset.HashSet[string]
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
		defaultBases = hashset.NewHashSet("information_schema", "mysql", "sys", "performance_schema")
	}
}

type OuterService interface {
	// CreateDb 创建数据库
	CreateDb(context.Context, CreateDbReqDTO) error
	// UpdateDb 编辑数据库
	UpdateDb(context.Context, UpdateDbReqDTO) error
	// DeleteDb 删除数据库
	DeleteDb(context.Context, DeleteDbReqDTO) error
	// ListDb 数据库详细列表
	ListDb(context.Context, ListDbReqDTO) ([]DbDTO, int64, error)
	// ListSimpleDb 数据库简单列表
	ListSimpleDb(context.Context, ListSimpleDbReqDTO) ([]SimpleDbDTO, error)
	// ApplyReadPerm 申请库表读权限
	ApplyReadPerm(context.Context, ApplyReadPermReqDTO) error
	// ListReadPermApplyByDba 展示审批列表
	ListReadPermApplyByDba(context.Context, ListReadPermApplyByDbaReqDTO) ([]ReadPermApplyDTO, int64, error)
	// ListReadPermApplyByOperator 展示申请的审批列表
	ListReadPermApplyByOperator(context.Context, ListReadPermApplyByOperatorReqDTO) ([]ReadPermApplyDTO, int64, error)
	// AgreeReadPermApply 同意审批
	AgreeReadPermApply(context.Context, AgreeReadPermApplyReqDTO) error
	// DisagreeReadPermApply 不同意审批
	DisagreeReadPermApply(context.Context, DisagreeReadPermApplyReqDTO) error
	// CancelReadPermApply 取消申请
	CancelReadPermApply(context.Context, CancelReadPermApplyReqDTO) error
	// GetReadPermApply 查看读权限审批单
	GetReadPermApply(context.Context, GetReadPermApplyReqDTO) (ReadPermApplyDTO, error)
	// ListReadPermByOperator 权限列表
	ListReadPermByOperator(context.Context, ListReadPermByOperatorReqDTO) ([]ReadPermDTO, int64, error)
	// DeleteReadPermByDba 删除权限
	DeleteReadPermByDba(context.Context, DeleteReadPermByDbaReqDTO) error
	// ListReadPermByDba 权限列表
	ListReadPermByDba(context.Context, ListReadPermByDbaReqDTO) ([]ReadPermDTO, int64, error)
	// ListAuthorizedDb 展示授权的数据库
	ListAuthorizedDb(context.Context, ListAuthorizedDbReqDTO) ([]SimpleDbDTO, error)
	// ListAuthorizedBase 展示授权的库
	ListAuthorizedBase(context.Context, ListAuthorizedBaseReqDTO) ([]string, error)
	// ListAuthorizedTable 展示授权的表
	ListAuthorizedTable(context.Context, ListAuthorizedTableReqDTO) ([]string, error)
	// GetCreateTableSql 展示建表语句
	GetCreateTableSql(context.Context, GetCreateSqlReqDTO) (string, error)
	// ShowTableIndex 展示索引语句
	ShowTableIndex(context.Context, ShowTableIndexReqDTO) ([]string, [][]string, error)
	// ExecuteSelectSql 搜索
	ExecuteSelectSql(context.Context, ExecuteSelectSqlReqDTO) (util.MysqlQueryResult, error)
	// ApplyDataUpdate 提数据库修改单
	ApplyDataUpdate(context.Context, ApplyDataUpdateReqDTO) ([]command.ValidateUpdateResult, bool, error)
	// ExplainDataUpdate 数据库修改单的执行计划
	ExplainDataUpdate(context.Context, ExplainDataUpdateReqDTO) (string, error)
	// ListDataUpdateApplyByDba 数据库修改审批单
	ListDataUpdateApplyByDba(context.Context, ListDataUpdateApplyByDbaReqDTO) ([]DataUpdateApplyDTO, int64, error)
	// ListDataUpdateApplyByOperator 申请的数据库修改审批单
	ListDataUpdateApplyByOperator(context.Context, ListDataUpdateApplyByOperatorReqDTO) ([]DataUpdateApplyDTO, int64, error)
	// AgreeDataUpdateApply 同意修改单
	AgreeDataUpdateApply(context.Context, AgreeDbUpdateReqDTO) error
	// DisagreeDataUpdateApply 不同意修改单
	DisagreeDataUpdateApply(context.Context, DisagreeDataUpdateApplyReqDTO) error
	// CancelDataUpdateApply 取消修改单
	CancelDataUpdateApply(context.Context, CancelDataUpdateApplyReqDTO) error
	// ExecuteDataUpdateApply 执行修改单
	ExecuteDataUpdateApply(context.Context, ExecuteDataUpdateApplyReqDTO) error
	// AskToExecuteDataUpdateApply 请求执行修改单
	AskToExecuteDataUpdateApply(context.Context, AskToExecuteDataUpdateApplyReqDTO) error
}
