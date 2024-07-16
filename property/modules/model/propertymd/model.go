package propertymd

import "time"

const (
	PropertyHistoryTableName = "zproperty_history"
	PropertyFileTableName    = "zproperty_file"
	EtcdNodeTableName        = "zproperty_etcd_node"
	DeployTableName          = "zproperty_deploy"
)

type EtcdNode struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	Name      string    `json:"name"`
	AppId     string    `json:"appId"`
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

type File struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Name    string    `json:"name"`
	Env     string    `json:"env"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*File) TableName() string {
	return PropertyFileTableName
}

type History struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	FileId      int64     `json:"fileId"`
	Content     string    `json:"content"`
	Version     string    `json:"version"`
	LastVersion string    `json:"lastVersion"`
	Creator     string    `json:"creator"`
	Created     time.Time `json:"created" xorm:"created"`
}

func (*History) TableName() string {
	return PropertyHistoryTableName
}

type Deploy struct {
	Id        int64 `json:"id" xorm:"pk autoincr"`
	HistoryId int64 `json:"historyId"`
	// 方便查询删除
	FileId int64 `json:"fileId"`
	// 方便查询删除
	AppId string `json:"appId"`

	NodeName  string    `json:"nodeName"`
	Endpoints string    `json:"endpoints"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Creator   string    `json:"creator"`
	Created   time.Time `json:"created" xorm:"created"`
}

func (*Deploy) TableName() string {
	return DeployTableName
}
