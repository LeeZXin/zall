package mysqldbmd

import "time"

type InsertDbReqDTO struct {
	Name   string
	Config Config
}

type UpdateDbReqDTO struct {
	Id     int64
	Name   string
	Config Config
}

type InsertReadPermApplyReqDTO struct {
	Account      string
	DbId         int64
	DbName       string
	AccessBase   string
	AccessTables string
	OrderStatus  ReadPermApplyStatus
	ExpireDay    int
	ApplyReason  string
}

type ListReadPermApplyReqDTO struct {
	PageNum     int
	PageSize    int
	Account     string
	DbId        int64
	ApplyStatus ReadPermApplyStatus
}

type ListDataUpdateApplyReqDTO struct {
	PageNum     int
	PageSize    int
	Account     string
	OrderStatus DataUpdateApplyStatus
}

type InsertReadPermReqDTO struct {
	Account     string
	DbId        int64
	AccessBase  string
	AccessTable string
	ApplyId     int64
	Expired     time.Time
}

type UpdateReadPermApplyStatusReqDTO struct {
	Id             int64
	NewStatus      ReadPermApplyStatus
	OldStatus      ReadPermApplyStatus
	Auditor        string
	DisagreeReason string
}

type UpdateDataUpdateApplyStatusReqDTO struct {
	Id             int64
	NewStatus      DataUpdateApplyStatus
	OldStatus      DataUpdateApplyStatus
	Auditor        string
	DisagreeReason string
}

type ExistReadPermReqDTO struct {
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
	OrderStatus DataUpdateApplyStatus
}

type PageDbReqDTO struct {
	PageNum  int
	PageSize int
	Name     string
}

type PageReadPermReqDTO struct {
	Account  string
	PageNum  int
	PageSize int
	DbId     int64
}

type ListReadPermByAccountReqDTO struct {
	DbId       int64
	AccessBase string
	Account    string
	Cols       []string
}
