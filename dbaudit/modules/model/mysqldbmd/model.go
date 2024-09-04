package mysqldbmd

import (
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	DbTableName              = "zdb_mysql_db"
	ReadPermTableName        = "zdb_mysql_read_perm"
	ReadPermApplyTableName   = "zdb_mysql_read_perm_apply"
	DataUpdateApplyTableName = "zdb_mysql_data_update_apply"
)

type Config struct {
	WriteNode Node `json:"writeNode"`
	ReadNode  Node `json:"readNode"`
}

type Node struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Node) IsValid() bool {
	return len(c.Host) > 0 && util.GenIpPortPattern().MatchString(c.Host) && len(c.Username) > 0
}

func (c *Config) IsValid() bool {
	return c.WriteNode.IsValid() && c.ReadNode.IsValid()
}

type Db struct {
	Id      int64                        `json:"id" xorm:"pk autoincr"`
	Name    string                       `json:"name"`
	Config  *xormutil.Conversion[Config] `json:"config"`
	Created time.Time                    `json:"created" xorm:"created"`
	Updated time.Time                    `json:"updated" xorm:"updated"`
}

func (*Db) TableName() string {
	return DbTableName
}

type ReadPerm struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	Account     string    `json:"account"`
	DbId        int64     `json:"dbId"`
	AccessBase  string    `json:"accessBase"`
	AccessTable string    `json:"accessTable"`
	ApplyId     int64     `json:"applyId"`
	Expired     time.Time `json:"expired"`
	Created     time.Time `json:"created" xorm:"created"`
}

func (*ReadPerm) TableName() string {
	return ReadPermTableName
}

func (p *ReadPerm) IsExpired() bool {
	return p.Expired.Before(time.Now())
}

type ReadPermApplyStatus int

const (
	AllReadPermApplyStatus ReadPermApplyStatus = iota
	PendingReadPermApplyStatus
	AgreeReadPermApplyStatus
	DisagreeReadPermApplyStatus
	CanceledReadPermApplyStatus
)

func (s ReadPermApplyStatus) IsValid() bool {
	switch s {
	case AllReadPermApplyStatus, PendingReadPermApplyStatus, AgreeReadPermApplyStatus,
		DisagreeReadPermApplyStatus, CanceledReadPermApplyStatus:
		return true
	default:
		return false
	}
}

type ReadPermApply struct {
	Id             int64               `json:"id" xorm:"pk autoincr"`
	Account        string              `json:"account"`
	DbId           int64               `json:"dbId"`
	DbName         string              `json:"dbName"`
	AccessBase     string              `json:"accessBase"`
	AccessTables   string              `json:"accessTables"`
	ApplyStatus    ReadPermApplyStatus `json:"applyStatus"`
	Auditor        string              `json:"auditor"`
	ExpireDay      int                 `json:"expireDay"`
	ApplyReason    string              `json:"applyReason"`
	DisagreeReason string              `json:"disagreeReason"`
	Created        time.Time           `json:"created" xorm:"created"`
	Updated        time.Time           `json:"updated" xorm:"updated"`
}

func (*ReadPermApply) TableName() string {
	return ReadPermApplyTableName
}

type DataUpdateApplyStatus int

const (
	AllDataUpdateApplyStatus DataUpdateApplyStatus = iota
	PendingDataUpdateApplyStatus
	AgreeDataUpdateApplyStatus
	DisagreeDataUpdateApplyStatus
	CanceledDataUpdateApplyStatus
	AskToExecuteDataUpdateApplyStatus
	ExecutedDataUpdateApplyStatus
)

func (s DataUpdateApplyStatus) IsUnExecuted() bool {
	switch s {
	case PendingDataUpdateApplyStatus, AgreeDataUpdateApplyStatus, AskToExecuteDataUpdateApplyStatus:
		return true
	default:
		return false
	}
}

func (s DataUpdateApplyStatus) IsExecutable() bool {
	switch s {
	case AgreeDataUpdateApplyStatus, AskToExecuteDataUpdateApplyStatus:
		return true
	default:
		return false
	}
}

func (s DataUpdateApplyStatus) IsValid() bool {
	switch s {
	case PendingDataUpdateApplyStatus, AgreeDataUpdateApplyStatus,
		DisagreeDataUpdateApplyStatus, CanceledDataUpdateApplyStatus,
		ExecutedDataUpdateApplyStatus, AskToExecuteDataUpdateApplyStatus, AllDataUpdateApplyStatus:
		return true
	default:
		return false
	}
}

type DataUpdateApply struct {
	Id                              int64                 `json:"id" xorm:"pk autoincr"`
	Account                         string                `json:"account"`
	DbId                            int64                 `json:"dbId"`
	AccessBase                      string                `json:"accessBase"`
	UpdateCmd                       string                `json:"updateCmd"`
	ApplyStatus                     DataUpdateApplyStatus `json:"applyStatus"`
	Auditor                         string                `json:"auditor"`
	Executor                        string                `json:"executor"`
	ApplyReason                     string                `json:"applyReason"`
	DisagreeReason                  string                `json:"disagreeReason"`
	ExecuteLog                      string                `json:"executeLog"`
	ExecuteImmediatelyAfterApproval bool                  `json:"executeImmediatelyAfterApproval"`
	Created                         time.Time             `json:"created" xorm:"created"`
	Updated                         time.Time             `json:"updated" xorm:"updated"`
}

func (*DataUpdateApply) TableName() string {
	return DataUpdateApplyTableName
}
