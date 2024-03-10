package detectmd

import "time"

const (
	TcpDetectTableName = "ztcp_detect"
	InstanceTableName  = "ztcp_instance"
	LogTableName       = "ztcp_detect_log"
)

type TcpDetect struct {
	Id            int64     `xorm:"pk autoincr"`
	Ip            string    `json:"ip"`
	Port          int       `json:"port"`
	Name          string    `json:"name"`
	Enabled       bool      `json:"enabled"`
	HeartbeatTime int64     `json:"heartbeatTime"`
	Created       time.Time `json:"created" xorm:"created"`
	Updated       time.Time `json:"updated" xorm:"updated"`
}

func (*TcpDetect) TableName() string {
	return TcpDetectTableName
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

type DetectLog struct {
	Id       int64     `xorm:"pk autoincr"`
	DetectId int64     `json:"detectId"`
	Ip       string    `json:"ip"`
	Port     int       `json:"port"`
	Valid    bool      `json:"valid"`
	Created  time.Time `json:"created" xorm:"created"`
}

func (*DetectLog) TableName() string {
	return LogTableName
}
