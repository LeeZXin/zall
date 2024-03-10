package propmd

import "time"

const (
	PropHistoryTableName = "zprop_prop_history"
	PropContentTableName = "zprop_prop_content"
	EtcdNodeTableName    = "zprop_etcd_node"
	EtcdAuthTableName    = "zprop_etcd_auth"
	PropDeployTableName  = "zprop_prop_deploy"
)

type EtcdNode struct {
	Id        int64     `xorm:"pk autoincr"`
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
	Id       int64     `xorm:"pk autoincr"`
	AppId    string    `json:"appId"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created" xorm:"created"`
}

func (*EtcdAuth) TableName() string {
	return EtcdAuthTableName
}

type PropContent struct {
	Id      int64     `xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Name    string    `json:"name"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*PropContent) TableName() string {
	return PropContentTableName
}

type PropHistory struct {
	Id        int64     `xorm:"pk autoincr"`
	ContentId int64     `json:"contentId"`
	Content   string    `json:"content"`
	Version   string    `json:"version"`
	Created   time.Time `json:"created" xorm:"created"`
}

func (*PropHistory) TableName() string {
	return PropHistoryTableName
}

type PropDeploy struct {
	Id           int64     `xorm:"pk autoincr"`
	ContentId    int64     `json:"contentId"`
	Content      string    `json:"content"`
	Version      string    `json:"version"`
	NodeId       string    `json:"nodeId"`
	ContentAppId string    `json:"contentAppId"`
	ContentName  string    `json:"contentName"`
	Endpoints    string    `json:"endpoints"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Deleted      bool      `json:"delete"`
	Created      time.Time `json:"created" xorm:"created"`
}

func (*PropDeploy) TableName() string {
	return PropDeployTableName
}
