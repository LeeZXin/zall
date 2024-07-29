package mysqldbapi

import (
	"github.com/LeeZXin/zall/dbaudit/modules/model/mysqldbmd"
	"github.com/LeeZXin/zall/dbaudit/modules/service/mysqldbsrv/command"
)

type CreateDbReqVO struct {
	Name   string           `json:"name"`
	Config mysqldbmd.Config `json:"config"`
}

type UpdateDbReqVO struct {
	DbId   int64            `json:"dbId"`
	Name   string           `json:"name"`
	Config mysqldbmd.Config `json:"config"`
}

type DbVO struct {
	Id      int64            `json:"id"`
	Name    string           `json:"name"`
	Config  mysqldbmd.Config `json:"config"`
	Created string           `json:"created"`
}

type SimpleDbVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ApplyDbPermReqVO struct {
	DbId         int64  `json:"dbId"`
	AccessBase   string `json:"accessBase"`
	AccessTables string `json:"accessTables"`
	ApplyReason  string `json:"applyReason"`
	ExpireDay    int    `json:"expireDay"`
}

type DisagreeReadPermReqVO struct {
	ApplyId        int64  `json:"applyId"`
	DisagreeReason string `json:"disagreeReason"`
}

type ListDbPermByDbaReqVO struct {
	PageNum int    `json:"pageNum"`
	Account string `json:"account"`
	DbId    int64  `json:"dbId"`
}

type ListDbReqVO struct {
	PageNum int    `json:"pageNum"`
	Name    string `json:"name"`
}

type ListReadPermApplyByDbaReqVO struct {
	DbId        int64                         `json:"dbId"`
	PageNum     int                           `json:"pageNum"`
	ApplyStatus mysqldbmd.ReadPermApplyStatus `json:"applyStatus"`
}

type listReadPermApplyByOperatorReqVO struct {
	PageNum     int                           `json:"pageNum"`
	ApplyStatus mysqldbmd.ReadPermApplyStatus `json:"applyStatus"`
}

type ReadPermApplyVO struct {
	Id             int64                         `json:"id"`
	Account        string                        `json:"account"`
	DbId           int64                         `json:"dbId"`
	DbName         string                        `json:"dbName"`
	AccessBase     string                        `json:"accessBase"`
	AccessTables   string                        `json:"accessTables"`
	ApplyStatus    mysqldbmd.ReadPermApplyStatus `json:"applyStatus"`
	Auditor        string                        `json:"auditor"`
	ExpireDay      int                           `json:"expireDay"`
	ApplyReason    string                        `json:"applyReason"`
	DisagreeReason string                        `json:"disagreeReason"`
	Created        string                        `json:"created"`
	Updated        string                        `json:"updated"`
}

type ListDataUpdateApplyByDbaReqVO struct {
	PageNum     int                             `json:"pageNum"`
	DbId        int64                           `json:"dbId"`
	ApplyStatus mysqldbmd.DataUpdateApplyStatus `json:"applyStatus"`
}

type ListDataUpdateApplyByOperatorReqVO struct {
	PageNum     int                             `json:"pageNum"`
	ApplyStatus mysqldbmd.DataUpdateApplyStatus `json:"applyStatus"`
}

type DataUpdateApplyVO struct {
	Id               int64                           `json:"id"`
	Account          string                          `json:"account"`
	DbId             int64                           `json:"dbId"`
	DbName           string                          `json:"dbName"`
	AccessBase       string                          `json:"accessBase"`
	UpdateCmd        string                          `json:"updateCmd"`
	ApplyStatus      mysqldbmd.DataUpdateApplyStatus `json:"applyStatus"`
	Executor         string                          `json:"executor"`
	Auditor          string                          `json:"auditor"`
	ApplyReason      string                          `json:"applyReason"`
	DisagreeReason   string                          `json:"disagreeReason"`
	ExecuteLog       string                          `json:"executeLog"`
	ExecuteWhenApply bool                            `json:"executeWhenApply"`
	Created          string                          `json:"created"`
	Updated          string                          `json:"updated"`
	IsUnExecuted     bool                            `json:"isUnExecuted"`
}

type ReadPermVO struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	DbId        int64  `json:"dbId"`
	DbName      string `json:"dbName"`
	AccessBase  string `json:"accessBase"`
	AccessTable string `json:"accessTable"`
	Created     string `json:"created"`
	Expired     string `json:"expired"`
	ApplyId     int64  `json:"applyId"`
}

type ExecuteSelectSqlReqVO struct {
	DbId       int64  `json:"dbId"`
	AccessBase string `json:"accessBase"`
	Cmd        string `json:"cmd"`
	Limit      int    `json:"limit"`
}

type ExecuteSelectSqlResultVO struct {
	Columns []string            `json:"columns"`
	Data    []map[string]string `json:"data"`
}

type ApplyDataUpdateReqVO struct {
	DbId             int64  `json:"dbId"`
	AccessBase       string `json:"accessBase"`
	Cmd              string `json:"cmd"`
	ExecuteWhenApply bool   `json:"executeWhenApply"`
	ApplyReason      string `json:"applyReason"`
}

type ApplyDbUpdateResultVO struct {
	Results []command.ValidateUpdateResult `json:"results"`
	AllPass bool                           `json:"allPass"`
}

type DisagreeDataUpdateReqVO struct {
	ApplyId        int64  `json:"applyId"`
	DisagreeReason string `json:"disagreeReason"`
}

type ListReadPermByOperatorReqVO struct {
	DbId    int64 `json:"dbId"`
	PageNum int   `json:"pageNum"`
}

type GetCreateTableSqlReqVO struct {
	DbId        int64  `json:"dbId"`
	AccessBase  string `json:"accessBase"`
	AccessTable string `json:"accessTable"`
}

type ShowTableIndexReqVO struct {
	DbId        int64  `json:"dbId"`
	AccessBase  string `json:"accessBase"`
	AccessTable string `json:"accessTable"`
}

type ShowTableIndexRespVO struct {
	Columns []string            `json:"columns"`
	Data    []map[string]string `json:"data"`
}
