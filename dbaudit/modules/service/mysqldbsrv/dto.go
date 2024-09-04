package mysqldbsrv

import (
	"github.com/LeeZXin/zall/dbaudit/modules/model/mysqldbmd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateDbReqDTO struct {
	Name     string              `json:"name"`
	Config   mysqldbmd.Config    `json:"config"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateDbReqDTO) IsValid() error {
	if !mysqldbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Config.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateDbReqDTO struct {
	DbId     int64               `json:"dbId"`
	Name     string              `json:"name"`
	Config   mysqldbmd.Config    `json:"config"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateDbReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Config.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteDbReqDTO struct {
	DbId     int64               `json:"dbId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteDbReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDbReqDTO struct {
	PageNum  int                 `json:"pageNum"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDbReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.Name) > 0 && !mysqldbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DbDTO struct {
	Id      int64
	Name    string
	Config  mysqldbmd.Config
	Created time.Time
}

type ListSimpleDbReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListSimpleDbReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SimpleDbDTO struct {
	Id   int64
	Name string
}

type ApplyReadPermReqDTO struct {
	DbId         int64               `json:"dbId"`
	AccessBase   string              `json:"accessBase"`
	AccessTables string              `json:"accessTables"`
	ApplyReason  string              `json:"reason"`
	ExpireDay    int                 `json:"expireDay"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *ApplyReadPermReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if len(r.AccessTables) == 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsReasonValid(r.ApplyReason) {
		return util.InvalidArgsError()
	}
	switch r.ExpireDay {
	case 1, 30, 90, 180, 365:
	default:
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListReadPermApplyByDbaReqDTO struct {
	DbId        int64                         `json:"dbId"`
	PageNum     int                           `json:"pageNum"`
	ApplyStatus mysqldbmd.ReadPermApplyStatus `json:"applyStatus"`
	Operator    apisession.UserInfo           `json:"operator"`
}

func (r *ListReadPermApplyByDbaReqDTO) IsValid() error {
	if r.DbId < 0 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListReadPermApplyByOperatorReqDTO struct {
	PageNum     int                           `json:"pageNum"`
	ApplyStatus mysqldbmd.ReadPermApplyStatus `json:"applyStatus"`
	Operator    apisession.UserInfo           `json:"operator"`
}

func (r *ListReadPermApplyByOperatorReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.ApplyStatus.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AgreeReadPermApplyReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AgreeReadPermApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisagreeReadPermApplyReqDTO struct {
	ApplyId        int64               `json:"applyId"`
	DisagreeReason string              `json:"disagreeReason"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *DisagreeReadPermApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsReasonValid(r.DisagreeReason) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CancelReadPermApplyReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CancelReadPermApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetReadPermApplyReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetReadPermApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteReadPermByDbaReqDTO struct {
	PermId   int64               `json:"permId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteReadPermByDbaReqDTO) IsValid() error {
	if r.PermId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListReadPermByOperatorReqDTO struct {
	PageNum  int                 `json:"pageNum"`
	DbId     int64               `json:"dbId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListReadPermByOperatorReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	// 允许dbId为0 0则标识所有数据库
	if r.DbId < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListReadPermByDbaReqDTO struct {
	PageNum  int                 `json:"pageNum"`
	DbId     int64               `json:"dbId"`
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListReadPermByDbaReqDTO) IsValid() error {
	if r.DbId < 0 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if r.Account != "" && !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ReadPermApplyDTO struct {
	Id             int64
	Account        string
	DbId           int64
	DbName         string
	AccessBase     string
	AccessTables   string
	ApplyStatus    mysqldbmd.ReadPermApplyStatus
	Auditor        string
	ExpireDay      int
	ApplyReason    string
	DisagreeReason string
	Created        time.Time
	Updated        time.Time
}

type ReadPermDTO struct {
	Id          int64
	Account     string
	DbId        int64
	DbName      string
	AccessBase  string
	AccessTable string
	Created     time.Time
	Expired     time.Time
	ApplyId     int64
}

type ListAuthorizedBaseReqDTO struct {
	DbId     int64               `json:"dbId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAuthorizedBaseReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAuthorizedTableReqDTO struct {
	DbId       int64               `json:"dbId"`
	AccessBase string              `json:"accessBase"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ListAuthorizedTableReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) || defaultBases.Contains(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetCreateSqlReqDTO struct {
	DbId        int64               `json:"dbId"`
	AccessBase  string              `json:"accessBase"`
	AccessTable string              `json:"accessTable"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *ShowTableIndexReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) || defaultBases.Contains(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsTableNameValid(r.AccessTable) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ShowTableIndexReqDTO struct {
	DbId        int64               `json:"dbId"`
	AccessBase  string              `json:"accessBase"`
	AccessTable string              `json:"accessTable"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *GetCreateSqlReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) || defaultBases.Contains(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsTableNameValid(r.AccessTable) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAuthorizedDbReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAuthorizedDbReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ExecuteSelectSqlReqDTO struct {
	DbId       int64               `json:"dbId"`
	AccessBase string              `json:"accessBase"`
	Cmd        string              `json:"cmd"`
	Limit      int                 `json:"limit"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ExecuteSelectSqlReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) || defaultBases.Contains(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if len(r.Cmd) == 0 || len(r.Cmd) > 10240 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Limit < 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	return nil
}

type ExplainDbReqDTO struct {
	DbId       int64               `json:"dbId"`
	AccessBase string              `json:"accessBase"`
	Cmd        string              `json:"cmd"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ExplainDbReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if len(r.Cmd) == 0 || len(r.Cmd) > 2048 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ApplyDataUpdateReqDTO struct {
	DbId                            int64               `json:"dbId"`
	AccessBase                      string              `json:"accessBase"`
	Cmd                             string              `json:"cmd"`
	ApplyReason                     string              `json:"applyReason"`
	ExecuteImmediatelyAfterApproval bool                `json:"executeImmediatelyAfterApproval"`
	Operator                        apisession.UserInfo `json:"operator"`
}

func (r *ApplyDataUpdateReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if len(r.Cmd) == 0 || len(r.Cmd) > 10240 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsReasonValid(r.ApplyReason) {
		return util.InvalidArgsError()
	}
	return nil
}

type ExplainDataUpdateReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ExplainDataUpdateReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDataUpdateApplyByDbaReqDTO struct {
	PageNum     int                             `json:"pageNum"`
	DbId        int64                           `json:"dbId"`
	Operator    apisession.UserInfo             `json:"operator"`
	ApplyStatus mysqldbmd.DataUpdateApplyStatus `json:"applyStatus"`
}

func (r *ListDataUpdateApplyByDbaReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.DbId < 0 {
		return util.InvalidArgsError()
	}
	if !r.ApplyStatus.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDataUpdateApplyByOperatorReqDTO struct {
	PageNum     int                             `json:"pageNum"`
	ApplyStatus mysqldbmd.DataUpdateApplyStatus `json:"applyStatus"`
	Operator    apisession.UserInfo             `json:"operator"`
}

func (r *ListDataUpdateApplyByOperatorReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.ApplyStatus.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DataUpdateApplyDTO struct {
	Id                              int64
	Account                         string
	DbId                            int64
	DbName                          string
	AccessBase                      string
	UpdateCmd                       string
	ApplyStatus                     mysqldbmd.DataUpdateApplyStatus
	Auditor                         string
	Executor                        string
	ApplyReason                     string
	DisagreeReason                  string
	ExecuteLog                      string
	ExecuteImmediatelyAfterApproval bool
	Created                         time.Time
	Updated                         time.Time
}

type AgreeDbUpdateReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AgreeDbUpdateReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisagreeDataUpdateApplyReqDTO struct {
	ApplyId        int64               `json:"applyId"`
	DisagreeReason string              `json:"disagreeReason"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *DisagreeDataUpdateApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsReasonValid(r.DisagreeReason) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CancelDataUpdateApplyReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CancelDataUpdateApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ExecuteDataUpdateApplyReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ExecuteDataUpdateApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AskToExecuteDataUpdateApplyReqDTO struct {
	ApplyId  int64               `json:"applyId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AskToExecuteDataUpdateApplyReqDTO) IsValid() error {
	if r.ApplyId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
