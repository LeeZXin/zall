package prommd

import (
	"time"
)

const (
	ScrapeTableName = "zprom_scrape"
)

type TargetType int

const (
	DiscoveryTargetType TargetType = iota + 1
	HostTargetType
)

func (t TargetType) IsValid() bool {
	switch t {
	case DiscoveryTargetType, HostTargetType:
		return true
	default:
		return false
	}
}

type Scrape struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	Endpoint   string     `json:"endpoint"`
	AppId      string     `json:"appId"`
	Target     string     `json:"target"`
	TargetType TargetType `json:"targetType"`
	Env        string     `json:"env"`
	Created    time.Time  `json:"created" xorm:"created"`
	Updated    time.Time  `json:"updated" xorm:"updated"`
}

func (*Scrape) TableName() string {
	return ScrapeTableName
}
