package mysqldbmd

import "time"

type InsertDbReqDTO struct {
	Name     string
	DbHost   string
	Username string
	Password string
}

type UpdateDbReqDTO struct {
	Id       int64
	Name     string
	DbHost   string
	Username string
	Password string
}

type InsertPermApprovalOrderReqDTO struct {
	Account      string
	DbId         int64
	AccessBase   string
	AccessTables AccessTables
	PermType     PermType
	OrderStatus  PermOrderStatus
	ExpireDay    int
	Reason       string
}

type ListPermApprovalOrderReqDTO struct {
	Cursor      int64
	Limit       int
	Account     string
	OrderStatus int
}

type ListUpdateApprovalOrderReqDTO struct {
	Cursor      int64
	Limit       int
	Account     string
	OrderStatus int
}

type InsertPermReqDTO struct {
	Account     string
	DbId        int64
	AccessBase  string
	AccessTable string
	PermType    PermType
	Expired     time.Time
}

type UpdatePermApprovalOrderStatusReqDTO struct {
	Id             int64
	NewStatus      PermOrderStatus
	OldStatus      PermOrderStatus
	Auditor        string
	DisagreeReason string
}

type UpdateUpdateApprovalOrderStatusReqDTO struct {
	Id             int64
	NewStatus      UpdateOrderStatus
	OldStatus      UpdateOrderStatus
	Auditor        string
	DisagreeReason string
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

type InsertUpdateApprovalOrderReqDTO struct {
	Name        string
	Account     string
	DbId        int64
	AccessBase  string
	UpdateCmd   string
	OrderStatus UpdateOrderStatus
}
