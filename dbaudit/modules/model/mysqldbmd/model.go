package mysqldbmd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	MysqlDbTableName                  = "zdb_mysql_node"
	MysqlPermTableName                = "zdb_mysql_perm"
	MysqlPermApprovalOrderTableName   = "zdb_mysql_perm_approval_order"
	MysqlUpdateApprovalOrderTableName = "zdb_mysql_update_approval_order"
)

type Db struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	Name     string    `json:"name"`
	DbHost   string    `json:"dbHost"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
}

func (*Db) TableName() string {
	return MysqlDbTableName
}

type PermType int

const (
	ReadPermType  PermType = 1
	WritePermType PermType = 1 << 1

	ReadWritePermType = ReadPermType | WritePermType
)

func (t PermType) HasReadPermType() bool {
	return t&ReadPermType == ReadPermType
}

func (t PermType) HasWritePermType() bool {
	return t&WritePermType == WritePermType
}

func (t PermType) Readable() string {
	switch t {
	case ReadPermType:
		return i18n.GetByKey(i18n.DbReadPermType)
	case ReadWritePermType:
		return i18n.GetByKey(i18n.DbReadWritePermType)
	default:
		return i18n.GetByKey(i18n.DbUnknownPermType)
	}
}

func (t PermType) IsValid() bool {
	switch t {
	case ReadPermType, ReadWritePermType:
		return true
	default:
		return false
	}
}

type Perm struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	Account     string    `json:"account"`
	DbId        int64     `json:"dbId"`
	AccessBase  string    `json:"accessBase"`
	AccessTable string    `json:"accessTable"`
	PermType    PermType  `json:"permType"`
	Expired     time.Time `json:"expired"`
	Created     time.Time `json:"created" xorm:"created"`
}

func (*Perm) TableName() string {
	return MysqlPermTableName
}

func (p *Perm) IsExpired() bool {
	return p.Expired.Before(time.Now())
}

type PermOrderStatus int

const (
	PendingPermOrderStatus PermOrderStatus = iota + 1
	AgreePermOrderStatus
	DisagreePermOrderStatus
	CanceledPermOrderStatus
)

func (s PermOrderStatus) Readable() string {
	switch s {
	case PendingPermOrderStatus:
		return i18n.GetByKey(i18n.DbPendingPermOrderStatus)
	case AgreePermOrderStatus:
		return i18n.GetByKey(i18n.DbAgreePermOrderStatus)
	case DisagreePermOrderStatus:
		return i18n.GetByKey(i18n.DbDisagreePermOrderStatus)
	case CanceledPermOrderStatus:
		return i18n.GetByKey(i18n.DbCanceledPermOrderStatus)
	default:
		return i18n.GetByKey(i18n.DbUnknownPermOrderStatus)
	}
}

func (s PermOrderStatus) IsValid() bool {
	switch s {
	case PendingPermOrderStatus, AgreePermOrderStatus,
		DisagreePermOrderStatus, CanceledPermOrderStatus:
		return true
	default:
		return false
	}
}

type AccessTables []string

func (k *AccessTables) FromDB(content []byte) error {
	ret := make(AccessTables, 0)
	err := json.Unmarshal(content, &ret)
	if err != nil {
		return err
	}
	*k = ret
	return nil
}

func (k *AccessTables) ToDB() ([]byte, error) {
	return json.Marshal(k)
}

type PermApprovalOrder struct {
	Id             int64           `json:"id" xorm:"pk autoincr"`
	Account        string          `json:"account"`
	DbId           int64           `json:"dbId"`
	AccessBase     string          `json:"accessBase"`
	AccessTables   AccessTables    `json:"accessTables"`
	PermType       PermType        `json:"permType"`
	OrderStatus    PermOrderStatus `json:"orderStatus"`
	Auditor        string          `json:"auditor"`
	ExpireDay      int             `json:"expireDay"`
	Reason         string          `json:"reason"`
	DisagreeReason string          `json:"disagreeReason"`
	Created        time.Time       `json:"created" xorm:"created"`
	Updated        time.Time       `json:"updated" xorm:"updated"`
}

func (*PermApprovalOrder) TableName() string {
	return MysqlPermApprovalOrderTableName
}

type UpdateOrderStatus int

const (
	PendingUpdateOrderStatus UpdateOrderStatus = iota + 1
	AgreeUpdateOrderStatus
	DisagreeUpdateOrderStatus
	CanceledUpdateOrderStatus
	AskToExecuteUpdateOrderStatus
	ExecutedUpdateOrderStatus
)

func (s UpdateOrderStatus) IsValid() bool {
	switch s {
	case PendingUpdateOrderStatus, AgreeUpdateOrderStatus,
		DisagreeUpdateOrderStatus, CanceledUpdateOrderStatus,
		ExecutedUpdateOrderStatus, AskToExecuteUpdateOrderStatus:
		return true
	default:
		return false
	}
}

func (s UpdateOrderStatus) Readable() string {
	switch s {
	case PendingUpdateOrderStatus:
		return i18n.GetByKey(i18n.DbPendingUpdateOrderStatus)
	case AgreeUpdateOrderStatus:
		return i18n.GetByKey(i18n.DbAgreeUpdateOrderStatus)
	case DisagreeUpdateOrderStatus:
		return i18n.GetByKey(i18n.DbDisagreeUpdateOrderStatus)
	case CanceledUpdateOrderStatus:
		return i18n.GetByKey(i18n.DbCanceledUpdateOrderStatus)
	case ExecutedUpdateOrderStatus:
		return i18n.GetByKey(i18n.DbExecutedUpdateOrderStatus)
	default:
		return i18n.GetByKey(i18n.DbUnknownPermOrderStatus)
	}
}

type UpdateApprovalOrder struct {
	Id             int64             `json:"id" xorm:"pk autoincr"`
	Name           string            `json:"name"`
	Account        string            `json:"account"`
	DbId           int64             `json:"dbId"`
	AccessBase     string            `json:"accessBase"`
	UpdateCmd      string            `json:"updateCmd"`
	OrderStatus    UpdateOrderStatus `json:"orderStatus"`
	Auditor        string            `json:"auditor"`
	DisagreeReason string            `json:"disagreeReason"`
	ExecuteLog     string            `json:"executeLog"`
	Created        time.Time         `json:"created" xorm:"created"`
	Updated        time.Time         `json:"updated" xorm:"updated"`
}

func (*UpdateApprovalOrder) TableName() string {
	return MysqlUpdateApprovalOrderTableName
}
