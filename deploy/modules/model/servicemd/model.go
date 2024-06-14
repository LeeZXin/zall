package servicemd

import (
	"github.com/LeeZXin/zall/pkg/deploy"
	"time"
)

const (
	ServiceTableName = "zservice_service"
)

type Service struct {
	Id          int64              `json:"id" xorm:"pk autoincr"`
	AppId       string             `json:"appId"`
	Name        string             `json:"name"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	Config      string             `json:"config"`
	Env         string             `json:"env"`
	IsEnabled   bool               `json:"isEnabled"`
	Probed      int64              `json:"probed"`
	Created     time.Time          `json:"created" xorm:"created"`
	Updated     time.Time          `json:"updated" xorm:"updated"`
}

func (*Service) TableName() string {
	return ServiceTableName
}
