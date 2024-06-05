package propertymd

import "time"

const (
	PropertyHistoryTableName = "zproperty_history"
	PropertyFileTableName    = "zproperty_file"
	EtcdNodeTableName        = "zprop_etcd_node"
	EtcdAuthTableName        = "zprop_etcd_auth"
	PropDeployTableName      = "zprop_prop_deploy"
)

type EtcdNode struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	NodeId    string    `json:"nodeId"`
	Endpoints string    `json:"endpoints"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (*EtcdNode) TableName() string {
	return EtcdNodeTableName
}

type EtcdAuth struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	AppId    string    `json:"appId"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created" xorm:"created"`
}

func (*EtcdAuth) TableName() string {
	return EtcdAuthTableName
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
	Env         string    `json:"env"`
	Created     time.Time `json:"created" xorm:"created"`
}

func (*History) TableName() string {
	return PropertyHistoryTableName
}

type PropDeploy struct {
	Id           int64     `json:"id" xorm:"pk autoincr"`
	ContentId    int64     `json:"contentId"`
	Content      string    `json:"content"`
	Version      string    `json:"version"`
	NodeId       string    `json:"nodeId"`
	ContentAppId string    `json:"contentAppId"`
	ContentName  string    `json:"contentName"`
	Endpoints    string    `json:"endpoints"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Creator      string    `json:"creator"`
	Deleted      bool      `json:"delete"`
	Created      time.Time `json:"created" xorm:"created"`
}

func (*PropDeploy) TableName() string {
	return PropDeployTableName
}
