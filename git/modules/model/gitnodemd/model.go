package gitnodemd

import "time"

const (
	NodeTableName = "zgit_git_node"
)

type Node struct {
	Id int64 `json:"id" xorm:"pk autoincr"`
	// 名称标识
	Name string `json:"name"`
	// HttpHost 节点http地址
	HttpHost string `json:"httpHost"`
	// SshHost 节点ssh地址
	SshHost string    `json:"sshHost"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Node) TableName() string {
	return NodeTableName
}
