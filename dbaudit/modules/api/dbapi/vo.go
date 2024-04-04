package dbapi

import (
	"github.com/LeeZXin/zall/dbaudit/modules/model/dbmd"
)

type InsertDbReqVO struct {
	Name     string      `json:"name"`
	DbHost   string      `json:"dbHost"`
	Username string      `json:"username"`
	Password string      `json:"password"`
	DbType   dbmd.DbType `json:"dbType"`
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
	DbType   string `json:"dbType"`
	Created  string `json:"created"`
}

type SimpleDbVO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	DbHost string `json:"dbHost"`
}

type ApplyDbPermReqVO struct {
	DbId        int64         `json:"dbId"`
	AccessTable string        `json:"accessTable"`
	Reason      string        `json:"reason"`
	ExpireDay   int           `json:"expireDay"`
	PermType    dbmd.PermType `json:"permType"`
}

type AgreeDbPermReqVO struct {
	OrderId int64 `json:"orderId"`
}

type DisagreeDbPermReqVO struct {
	OrderId int64 `json:"orderId"`
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

type ListApprovalOrderReqVO struct {
	Cursor      int64            `json:"cursor"`
	Limit       int              `json:"limit"`
	OrderStatus dbmd.OrderStatus `json:"orderStatus"`
}

type ApprovalOrderVO struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	DbId        int64  `json:"dbId"`
	DbHost      string `json:"dbHost"`
	DbName      string `json:"dbName"`
	AccessTable string `json:"accessTable"`
	PermType    string `json:"permType"`
	OrderStatus string `json:"orderStatus"`
	Auditor     string `json:"auditor"`
	ExpireDay   int    `json:"expireDay"`
	Reason      string `json:"reason"`
	Created     string `json:"created"`
}

type PermVO struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	DbId        int64  `json:"dbId"`
	DbHost      string `json:"dbHost"`
	DbName      string `json:"dbName"`
	AccessTable string `json:"accessTable"`
	PermType    string `json:"permType"`
	Created     string `json:"created"`
	Expired     string `json:"expired"`
}
