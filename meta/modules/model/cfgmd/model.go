package cfgmd

import "time"

const (
	SysCfgTableName = "zall_sys_cfg"
)

type SysCfg struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	CfgKey  string    `json:"cfgKey"`
	Content string    `json:"content"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*SysCfg) TableName() string {
	return SysCfgTableName
}
