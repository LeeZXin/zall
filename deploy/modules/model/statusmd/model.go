package statusmd

import (
	"github.com/LeeZXin/zall/pkg/status"
	"time"
)

const (
	ServiceTableName = "zallet_service"
)

type Service struct {
	Id            int64        `json:"id" xorm:"pk autoincr"`
	ServiceId     string       `json:"serviceId"`
	Pid           int          `json:"pid"`
	InstanceId    string       `json:"instanceId"`
	App           string       `json:"app"`
	AppYaml       *status.Yaml `json:"appYaml"`
	ServiceStatus string       `json:"serviceStatus"`
	ErrLog        string       `json:"errLog"`
	EventTime     int64        `json:"eventTime"`
	AgentHost     string       `json:"agentHost"`
	AgentToken    string       `json:"agentToken"`
	Env           string       `json:"env"`
	CpuPercent    int          `json:"cpuPercent"`
	MemPercent    int          `json:"memPercent"`
	Created       time.Time    `json:"created" xorm:"created"`
}

func (*Service) TableName() string {
	return ServiceTableName
}
