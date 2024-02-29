package gitnodemd

import "time"

const (
	GitNodeTableName = "git_node"
)

type GitNode struct {
	Id int64 `json:"id" xorm:"pk autoincr"`
	// 节点id
	NodeId string `json:"nodeId"`
	// HttpHost 节点http地址 多个 逗号隔开
	HttpHost string `json:"httpHost"`
	// SshHost 节点ssh地址 多个 逗号隔开
	SshHost string    `json:"sshHost"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*GitNode) TableName() string {
	return GitNodeTableName
}
