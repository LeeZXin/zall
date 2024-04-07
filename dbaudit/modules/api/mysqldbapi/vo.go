package mysqldbapi

import (
	"github.com/LeeZXin/zall/dbaudit/modules/model/mysqldbmd"
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv/command"
)

type InsertDbReqVO struct {
	Name     string `json:"name"`
	DbHost   string `json:"dbHost"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateDbReqVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	DbHost   string `json:"dbHost"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteDbReqVO struct {
	Id int64 `json:"id"`
}

type DbVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	DbHost   string `json:"dbHost"`
	Username string `json:"username"`
	Password string `json:"password"`
	Created  string `json:"created"`
}

type SimpleDbVO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	DbHost string `json:"dbHost"`
}

type ApplyDbPermReqVO struct {
	DbId         int64                  `json:"dbId"`
	AccessBase   string                 `json:"accessBase"`
	AccessTables mysqldbmd.AccessTables `json:"accessTables"`
	Reason       string                 `json:"reason"`
	ExpireDay    int                    `json:"expireDay"`
	PermType     mysqldbmd.PermType     `json:"permType"`
}

type AgreeDbPermReqVO struct {
	OrderId int64 `json:"orderId"`
}

type DisagreeDbPermReqVO struct {
	OrderId        int64  `json:"orderId"`
	DisagreeReason string `json:"disagreeReason"`
}

type CancelDbPermReqVO struct {
	OrderId int64 `json:"orderId"`
}

type ListDbPermReqVO struct {
	Cursor int64 `json:"cursor"`
	Limit  int   `json:"limit"`
}

type ListDbPermByAccountReqVO struct {
	Cursor  int64  `json:"cursor"`
	Limit   int    `json:"limit"`
	Account string `json:"account"`
}

type DeleteDbPermReqVO struct {
	PermId int64 `json:"permId"`
}

type ListPermApprovalOrderReqVO struct {
	Cursor      int64 `json:"cursor"`
	Limit       int   `json:"limit"`
	OrderStatus int   `json:"orderStatus"`
}

type ListAppliedPermApprovalOrderReqVO struct {
	Cursor      int64 `json:"cursor"`
	Limit       int   `json:"limit"`
	OrderStatus int   `json:"orderStatus"`
}

type PermApprovalOrderVO struct {
	Id           int64                  `json:"id"`
	Account      string                 `json:"account"`
	DbId         int64                  `json:"dbId"`
	DbHost       string                 `json:"dbHost"`
	DbName       string                 `json:"dbName"`
	AccessBase   string                 `json:"accessBase"`
	AccessTables mysqldbmd.AccessTables `json:"accessTables"`
	PermType     string                 `json:"permType"`
	OrderStatus  string                 `json:"orderStatus"`
	Auditor      string                 `json:"auditor"`
	ExpireDay    int                    `json:"expireDay"`
	Reason       string                 `json:"reason"`
	Created      string                 `json:"created"`
}

type ListUpdateApprovalOrderReqVO struct {
	Cursor      int64 `json:"cursor"`
	Limit       int   `json:"limit"`
	OrderStatus int   `json:"orderStatus"`
}

type ListAppliedUpdateApprovalOrderReqVO struct {
	Cursor      int64 `json:"cursor"`
	Limit       int   `json:"limit"`
	OrderStatus int   `json:"orderStatus"`
}

type UpdateApprovalOrderVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Account     string `json:"account"`
	DbId        int64  `json:"dbId"`
	DbHost      string `json:"dbHost"`
	DbName      string `json:"dbName"`
	AccessBase  string `json:"accessBase"`
	UpdateCmd   string `json:"updateCmd"`
	OrderStatus string `json:"orderStatus"`
	Auditor     string `json:"auditor"`
	ExecuteLog  string `json:"executeLog"`
	Created     string `json:"created"`
}

type PermVO struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	DbId        int64  `json:"dbId"`
	DbHost      string `json:"dbHost"`
	DbName      string `json:"dbName"`
	AccessBase  string `json:"accessBase"`
	AccessTable string `json:"accessTable"`
	PermType    string `json:"permType"`
	Created     string `json:"created"`
	Expired     string `json:"expired"`
}

type AllTablesReqVO struct {
	DbId       int64  `json:"dbId"`
	AccessBase string `json:"accessBase"`
}

type AllBasesReqVO struct {
	DbId int64 `json:"dbId"`
}

type SearchDbReqVO struct {
	DbId       int64  `json:"dbId"`
	AccessBase string `json:"accessBase"`
	Cmd        string `json:"cmd"`
	Limit      int    `json:"limit"`
}

type SearchDbResultVO struct {
	Columns []string   `json:"columns"`
	Result  [][]string `json:"result"`
}

type ApplyDbUpdateReqVO struct {
	Name       string `json:"name"`
	DbId       int64  `json:"dbId"`
	AccessBase string `json:"accessBase"`
	Cmd        string `json:"cmd"`
}

type ApplyDbUpdateResultVO struct {
	Results []command.ValidateUpdateResult `json:"results"`
	AllPass bool                           `json:"allPass"`
}

type AgreeDbUpdateReqVO struct {
	OrderId int64 `json:"orderId"`
}

type DisagreeDbUpdateReqVO struct {
	OrderId        int64  `json:"orderId"`
	DisagreeReason string `json:"disagreeReason"`
}

type CancelDbUpdateReqVO struct {
	OrderId int64 `json:"orderId"`
}

type ExecuteDbUpdateReqVO struct {
	OrderId int64 `json:"orderId"`
}

type AskToExecuteDbUpdateReqVO struct {
	OrderId int64 `json:"orderId"`
}
