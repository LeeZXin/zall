package dbsrv

import (
	"github.com/LeeZXin/zall/dbaudit/modules/model/dbmd"
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
	DbType   dbmd.DbType         `json:"dbType"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertDbReqDTO) IsValid() error {
	if !dbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.DbHost) {
		return util.InvalidArgsError()
	}
	if !dbmd.IsUsernameValid(r.Username) {
		return util.InvalidArgsError()
	}
	if !dbmd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	if !r.DbType.IsValid() {
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
	if !dbmd.IsDbNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !util.IpPortPattern.MatchString(r.DbHost) {
		return util.InvalidArgsError()
	}
	if !dbmd.IsUsernameValid(r.Username) {
		return util.InvalidArgsError()
	}
	if !dbmd.IsPasswordValid(r.Password) {
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
	DbType   dbmd.DbType
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
	DbId         int64               `json:"dbId"`
	AccessBase   string              `json:"accessBase"`
	AccessTables dbmd.AccessTables   `json:"accessTables"`
	Reason       string              `json:"reason"`
	ExpireDay    int                 `json:"expireDay"`
	PermType     dbmd.PermType       `json:"permType"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *ApplyDbPermReqDTO) IsValid() error {
	if r.DbId <= 0 {
		return util.InvalidArgsError()
	}
	if !dbmd.IsBaseNameValid(r.AccessBase) {
		return util.InvalidArgsError()
	}
	if len(r.AccessTables) == 0 {
		return util.InvalidArgsError()
	}
	for _, table := range r.AccessTables {
		if !dbmd.IsTableNameValid(table) {
			return util.InvalidArgsError()
		}
	}
	if !dbmd.IsReasonValid(r.Reason) {
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
	OrderStatus dbmd.OrderStatus    `json:"orderStatus"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *ListPermApprovalOrderReqDTO) IsValid() error {
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if !r.OrderStatus.IsValid() {
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
	OrderId  int64               `json:"orderId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DisagreeDbPermReqDTO) IsValid() error {
	if r.OrderId <= 0 {
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

type ApprovalOrderDTO struct {
	Id           int64
	Account      string
	DbId         int64
	DbHost       string
	DbName       string
	AccessBase   string
	AccessTables dbmd.AccessTables
	PermType     dbmd.PermType
	OrderStatus  dbmd.OrderStatus
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
	PermType    dbmd.PermType
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
	if !dbmd.IsBaseNameValid(r.AccessBase) {
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
	if !dbmd.IsBaseNameValid(r.AccessBase) {
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
