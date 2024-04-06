package dbmd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	DbTableName            = "zdb_node"
	PermTableName          = "zdb_perm"
	ApprovalOrderTableName = "zdb_approval_order"
)

type DbType int

const (
	MysqlDbType DbType = iota + 1
	RedisDbType
	MongoDbType
)

func (t DbType) IsValid() bool {
	switch t {
	case MysqlDbType, RedisDbType, MongoDbType:
		return true
	default:
		return false
	}
}

func (t DbType) Readable() string {
	switch t {
	case MysqlDbType:
		return "mysql"
	case RedisDbType:
		return "redis"
	case MongoDbType:
		return "mongo"
	default:
		return "unknown"
	}
}

type Db struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	Name     string    `json:"name"`
	DbHost   string    `json:"dbHost"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	DbType   DbType    `json:"dbType"`
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
}

func (*Db) TableName() string {
	return DbTableName
}

type PermType int

const (
	ReadPermType PermType = 1
)

func (t PermType) HasReadPermType() bool {
	return t&ReadPermType == ReadPermType
}

func (t PermType) Readable() string {
	switch t {
	case ReadPermType:
		return i18n.GetByKey(i18n.DbReadPermType)
	default:
		return i18n.GetByKey(i18n.DbUnknownPermType)
	}
}

func (t PermType) IsValid() bool {
	switch t {
	case ReadPermType:
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
	return PermTableName
}

func (p *Perm) IsExpired() bool {
	return p.Expired.Before(time.Now())
}

type OrderStatus int

const (
	PendingOrderStatus OrderStatus = iota + 1
	AgreeOrderStatus
	DisagreeOrderStatus
	CanceledOrderStatus
)

func (s OrderStatus) Readable() string {
	switch s {
	case PendingOrderStatus:
		return i18n.GetByKey(i18n.DbPendingApprovalOrderStatus)
	case AgreeOrderStatus:
		return i18n.GetByKey(i18n.DbAgreeApprovalOrderStatus)
	case DisagreeOrderStatus:
		return i18n.GetByKey(i18n.DbDisagreeApprovalOrderStatus)
	case CanceledOrderStatus:
		return i18n.GetByKey(i18n.DbCanceledApprovalOrderStatus)
	default:
		return i18n.GetByKey(i18n.DbUnknownApprovalOrderStatus)
	}
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case PendingOrderStatus, AgreeOrderStatus, DisagreeOrderStatus, CanceledOrderStatus:
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

type ApprovalOrder struct {
	Id           int64        `json:"id" xorm:"pk autoincr"`
	Account      string       `json:"account"`
	DbId         int64        `json:"dbId"`
	AccessBase   string       `json:"accessBase"`
	AccessTables AccessTables `json:"accessTables"`
	PermType     PermType     `json:"permType"`
	OrderStatus  OrderStatus  `json:"orderStatus"`
	Auditor      string       `json:"auditor"`
	ExpireDay    int          `json:"expireDay"`
	Reason       string       `json:"reason"`
	Created      time.Time    `json:"created" xorm:"created"`
	Updated      time.Time    `json:"updated" xorm:"updated"`
}

func (*ApprovalOrder) TableName() string {
	return ApprovalOrderTableName
}
