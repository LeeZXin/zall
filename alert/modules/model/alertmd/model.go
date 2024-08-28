package alertmd

import (
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	ConfigTableName  = "zalert_config"
	ExecuteTableName = "zalert_execute"
)

type Config struct {
	Id          int64                             `json:"id" xorm:"pk autoincr"`
	Name        string                            `json:"name"`
	AppId       string                            `json:"appId"`
	Content     *xormutil.Conversion[alert.Alert] `json:"content"`
	IntervalSec int                               `json:"intervalSec"`
	IsEnabled   bool                              `json:"isEnabled"`
	Env         string                            `json:"env"`
	Creator     string                            `json:"creator"`
	Created     time.Time                         `json:"created" xorm:"created"`
	Updated     time.Time                         `json:"updated" xorm:"updated"`
}

func (*Config) TableName() string {
	return ConfigTableName
}

func (c *Config) GetContent() alert.Alert {
	if c.Content == nil {
		return alert.Alert{}
	}
	return c.Content.Data
}

type Execute struct {
	Id         int64  `json:"id" xorm:"pk autoincr"`
	ConfigId   int64  `json:"configId"`
	IsEnabled  bool   `json:"isEnabled"`
	NextTime   int64  `json:"nextTime"`
	Env        string `json:"env"`
	RunVersion int64  `json:"runVersion"`
}

func (*Execute) TableName() string {
	return ExecuteTableName
}
