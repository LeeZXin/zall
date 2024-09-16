package discoverymd

import (
	"time"
)

const (
	EtcdNodeTableName        = "zdiscovery_etcd_node"
	AppEtcdNodeBindTableName = "zdiscovery_app_etcd_node_bind"
)

type EtcdNode struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name"`
	Env       string    `json:"env"`
	Endpoints string    `json:"endpoints"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (*EtcdNode) TableName() string {
	return EtcdNodeTableName
}

type AppEtcdNodeBind struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	NodeId  int64     `json:"nodeId"`
	Env     string    `json:"env"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*AppEtcdNodeBind) TableName() string {
	return AppEtcdNodeBindTableName
}
