package prommd

import (
	"github.com/LeeZXin/zall/pkg/i18n"
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

func (t TargetType) Readable() string {
	switch t {
	case DiscoveryTargetType:
		return i18n.GetByKey(i18n.PromScrapeDiscoveryTargetType)
	case HostTargetType:
		return i18n.GetByKey(i18n.PromScrapeHostTargetType)
	default:
		return i18n.GetByKey(i18n.PromScrapeUnknownTargetType)
	}
}

type Scrape struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	ServerUrl  string     `json:"serverUrl"`
	AppId      string     `json:"appId"`
	Target     string     `json:"target"`
	TargetType TargetType `json:"targetType"`
	Created    time.Time  `json:"created" xorm:"created"`
	Updated    time.Time  `json:"updated" xorm:"updated"`
}

func (*Scrape) TableName() string {
	return ScrapeTableName
}
