package deploymd

import (
	"github.com/LeeZXin/zall/pkg/deploy"
	"time"
)

const (
	DeployTableName = "zservice_deploy"
)

type Deploy struct {
	Id      int64          `json:"id" xorm:"pk autoincr"`
	AppId   string         `json:"appId"`
	Config  *deploy.Config `json:"config"`
	Created time.Time      `json:"created" xorm:"created"`
	Updated time.Time      `json:"updated" xorm:"updated"`
}

func (*Deploy) TableName() string {
	return DeployTableName
}
