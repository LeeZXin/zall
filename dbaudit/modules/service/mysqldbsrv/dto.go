package mysqldbsrv

import (
	"github.com/LeeZXin/zall/dbaudit/modules/model/mysqldbmd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"time"
)

type InsertDbReqDTO struct {
	Name     string              `json:"name"`
	DbHost   string              `json:"dbHost"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertDbReqDTO) IsValid() error {
	if !mysqldbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.DbHost) {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsUsernameValid(r.Username) {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateDbReqDTO struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	DbHost   string              `json:"dbHost"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateDbReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.DbHost) {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsUsernameValid(r.Username) {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteDbReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteDbReqDTO) IsValid() error {
	if r.Id <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDbReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDbReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DbDTO struct {
	Id       int64
	Name     string
	DbHost   string
	Username string
	Password string
	Created  time.Time
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
	Id     int64
	Name   string
	DbHost string
}

type ApplyDbPermReqDTO struct {
	DbId         int64                  `json:"dbId"`
	AccessBase   string                 `json:"accessBase"`
	AccessTables mysqldbmd.AccessTables `json:"accessTables"`
	Reason       string                 `json:"reason"`
	ExpireDay    int                    `json:"expireDay"`
	PermType     mysqldbmd.PermType     `json:"permType"`
	Operator     apisession.UserInfo    `json:"operator"`
}

func (r *ApplyDbPermReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if len(r.AccessTables) == 0 {
		return util.InvalidArgsError()
	}
	for _, table := range r.AccessTables {
		if !mysqldbmd.IsTableNameValid(table) {
			return util.InvalidArgsError()
		}
	}
	if !mysqldbmd.IsReasonValid(r.Reason) {
		return util.InvalidArgsError()
	}
	if r.ExpireDay <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.PermType.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPermApprovalOrderReqDTO struct {
	Cursor      int64               `json:"cursor"`
	Limit       int                 `json:"limit"`
	OrderStatus int                 `json:"orderStatus"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *ListPermApprovalOrderReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAppliedPermApprovalOrderReqDTO struct {
	Cursor      int64               `json:"cursor"`
	Limit       int                 `json:"limit"`
	OrderStatus int                 `json:"orderStatus"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *ListAppliedPermApprovalOrderReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AgreeDbPermReqDTO struct {
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AgreeDbPermReqDTO) IsValid() error {
	if r.OrderId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisagreeDbPermReqDTO struct {
	OrderId        int64               `json:"orderId"`
	DisagreeReason string              `json:"disagreeReason"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *DisagreeDbPermReqDTO) IsValid() error {
	if r.OrderId <= 0 {
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

type CancelDbPermReqDTO struct {
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CancelDbPermReqDTO) IsValid() error {
	if r.OrderId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteDbPermReqDTO struct {
	PermId   int64               `json:"permId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteDbPermReqDTO) IsValid() error {
	if r.PermId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDbPermReqDTO struct {
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDbPermReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDbPermByAccountReqDTO struct {
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDbPermByAccountReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PermApprovalOrderDTO struct {
	Id           int64
	Account      string
	DbId         int64
	DbHost       string
	DbName       string
	AccessBase   string
	AccessTables mysqldbmd.AccessTables
	PermType     mysqldbmd.PermType
	OrderStatus  mysqldbmd.PermOrderStatus
	Auditor      string
	ExpireDay    int
	Reason       string
	Created      time.Time
}

type PermDTO struct {
	Id          int64
	Account     string
	DbId        int64
	DbHost      string
	DbName      string
	AccessBase  string
	AccessTable string
	PermType    mysqldbmd.PermType
	Created     time.Time
	Expired     time.Time
}

type AllBasesReqDTO struct {
	DbId     int64               `json:"dbId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AllBasesReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AllTablesReqDTO struct {
	DbId       int64               `json:"dbId"`
	AccessBase string              `json:"accessBase"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *AllTablesReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !mysqldbmd.IsBaseNameValid(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SearchDbReqDTO struct {
	DbId       int64               `json:"dbId"`
	AccessBase string              `json:"accessBase"`
	Cmd        string              `json:"cmd"`
	Limit      int                 `json:"limit"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *SearchDbReqDTO) IsValid() error {
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

type ApplyDbUpdateReqDTO struct {
	Name       string              `json:"name"`
	DbId       int64               `json:"dbId"`
	AccessBase string              `json:"accessBase"`
	Cmd        string              `json:"cmd"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ApplyDbUpdateReqDTO) IsValid() error {
	if !mysqldbmd.IsUpdateOrderNameValid(r.Name) {
		return util.InvalidArgsError()
	}
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
	return nil
}

type ListUpdateApprovalOrderReqDTO struct {
	Cursor      int64               `json:"cursor"`
	Limit       int                 `json:"limit"`
	OrderStatus int                 `json:"orderStatus"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *ListUpdateApprovalOrderReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAppliedUpdateApprovalOrderReqDTO struct {
	Cursor      int64               `json:"cursor"`
	Limit       int                 `json:"limit"`
	OrderStatus int                 `json:"orderStatus"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *ListAppliedUpdateApprovalOrderReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateApprovalOrderDTO struct {
	Id             int64
	Name           string
	Account        string
	DbId           int64
	DbHost         string
	DbName         string
	AccessBase     string
	UpdateCmd      string
	OrderStatus    mysqldbmd.UpdateOrderStatus
	Auditor        string
	DisagreeReason string
	ExecuteLog     string
	Created        time.Time
}

type AgreeDbUpdateReqDTO struct {
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AgreeDbUpdateReqDTO) IsValid() error {
	if r.OrderId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisagreeDbUpdateReqDTO struct {
	OrderId        int64               `json:"orderId"`
	DisagreeReason string              `json:"disagreeReason"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *DisagreeDbUpdateReqDTO) IsValid() error {
	if r.OrderId <= 0 {
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

type CancelDbUpdateReqDTO struct {
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CancelDbUpdateReqDTO) IsValid() error {
	if r.OrderId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ExecuteDbUpdateReqDTO struct {
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ExecuteDbUpdateReqDTO) IsValid() error {
	if r.OrderId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AskToExecuteDbUpdateReqDTO struct {
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AskToExecuteDbUpdateReqDTO) IsValid() error {
	if r.OrderId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
