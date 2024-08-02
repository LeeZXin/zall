package discoverymd

import (
	"github.com/LeeZXin/zsf/services/lb"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	EtcdNodeTableName        = "zdiscovery_etcd_node"
	DownServiceTableName     = "zdiscovery_down_service"
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

type DownService struct {
	Id          int64                           `json:"id" xorm:"pk autoincr"`
	SourceId    int64                           `json:"sourceId"`
	AppId       string                          `json:"appId"`
	DownService *xormutil.Conversion[lb.Server] `json:"downService"`
	InstanceId  string                          `json:"instanceId"`
	Created     time.Time                       `json:"created" xorm:"created"`
}

func (*DownService) TableName() string {
	return DownServiceTableName
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
