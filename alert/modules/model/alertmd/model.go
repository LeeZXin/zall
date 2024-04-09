package alertmd

import (
	"github.com/LeeZXin/zall/pkg/alert"
	"time"
)

const (
	ConfigTableName   = "zalert_config"
	InstanceTableName = "zalert_instance"
)

type Config struct {
	Id          int64        `json:"id" xorm:"pk autoincr"`
	Name        string       `json:"name"`
	AppId       string       `json:"appId"`
	Content     *alert.Alert `json:"content"`
	IntervalSec int          `json:"intervalSec"`
	SilenceSec  int          `json:"silenceSec"`
	Enabled     bool         `json:"enabled"`
	NextTime    int64        `json:"nextTime"`
	Created     time.Time    `json:"created" xorm:"created"`
	Updated     time.Time    `json:"updated" xorm:"updated"`
}

func (*Config) TableName() string {
	return ConfigTableName
}

type Instance struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	InstanceId    string    `json:"instanceId"`
	HeartbeatTime int64     `json:"heartbeatTime"`
	Created       time.Time `json:"created" xorm:"created"`
}

func (*Instance) TableName() string {
	return InstanceTableName
}
