package dbmd

import "time"

type InsertDbReqDTO struct {
	Name     string
	DbHost   string
	Username string
	Password string
	DbType   DbType
}

type UpdateDbReqDTO struct {
	Id       int64
	Name     string
	DbHost   string
	Username string
	Password string
}

type InsertApprovalOrderReqDTO struct {
	Account      string
	DbId         int64
	AccessBase   string
	AccessTables AccessTables
	PermType     PermType
	OrderStatus  OrderStatus
	ExpireDay    int
	Reason       string
}

type ListApprovalOrderReqDTO struct {
	Cursor      int64
	Limit       int
	OrderStatus OrderStatus
}

type InsertPermReqDTO struct {
	Account     string
	DbId        int64
	AccessBase  string
	AccessTable string
	PermType    PermType
	Expired     time.Time
}

type UpdateApprovalOrderStatusReqDTO struct {
	Id        int64
	NewStatus OrderStatus
	OldStatus OrderStatus
	Auditor   string
}

type ListPermReqDTO struct {
	Cursor  int64
	Limit   int
	Account string
}

type SearchPermReqDTO struct {
	Account      string
	DbId         int64
	AccessBase   string
	AccessTables []string
}
