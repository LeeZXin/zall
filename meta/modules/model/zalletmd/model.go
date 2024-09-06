package zalletmd

import "time"

const (
	ZalletNodeTableName = "zallet_node"
)

type ZalletNode struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	NodeId     string    `json:"nodeId"`
	Name       string    `json:"name"`
	AgentHost  string    `json:"agentHost"`
	AgentToken string    `json:"agentToken"`
	Created    time.Time `json:"created" xorm:"created"`
	Updated    time.Time `json:"updated" xorm:"updated"`
}

func (*ZalletNode) TableName() string {
	return ZalletNodeTableName
}
