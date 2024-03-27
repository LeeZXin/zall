package deploymd

import (
	"github.com/LeeZXin/zall/pkg/deploy"
	"time"
)

const (
	TeamConfigTableName = "zservice_team_config"
	ConfigTableName     = "zservice_deploy_config"
	ServiceTableName    = "zservice_deploy_service"
	LogTableName        = "zservice_deploy_log"
	PlanTableName       = "zservice_deploy_plan"
)

// Config 部署配置
type Config struct {
	Id          int64              `json:"id" xorm:"pk autoincr"`
	AppId       string             `json:"appId"`
	Name        string             `json:"name"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	Content     string             `json:"content"`
	Created     time.Time          `json:"created" xorm:"created"`
	Updated     time.Time          `json:"updated" xorm:"updated"`
}

func (*Config) TableName() string {
	return ConfigTableName
}

type ActiveStatus int

const (
	OfflineStatus ActiveStatus = iota
	OnlineStatus
	ShutdownStatus
)

// Service 部署服务
type Service struct {
	Id       int64 `json:"id" xorm:"pk autoincr"`
	ConfigId int64 `json:"configId"`
	// 当前制品版本
	CurrProductVersion string `json:"productVersion"`
	// 上个制品版本
	LastProductVersion string             `json:"lastProductVersion"`
	ServiceType        deploy.ServiceType `json:"serviceType"`
	ServiceConfig      string             `json:"serviceConfig"`
	ActiveStatus       ActiveStatus       `json:"activeStatus"`
	Created            time.Time          `json:"created" xorm:"created"`
	Updated            time.Time          `json:"updated" xorm:"updated"`
}

func (*Service) TableName() string {
	return ServiceTableName
}

// Log 部署日志
type Log struct {
	Id             int64              `json:"id" xorm:"pk autoincr"`
	ConfigId       int64              `json:"configId"`
	AppId          string             `json:"appId"`
	ServiceType    deploy.ServiceType `json:"serviceType"`
	ServiceConfig  string             `json:"serviceConfig"`
	ProductVersion string             `json:"productVersion"`
	Operator       string             `json:"operator"`
	DeployOutput   string             `json:"deployOutput"`
	Created        time.Time          `json:"created" xorm:"created"`
}

func (*Log) TableName() string {
	return LogTableName
}

type PlanStatus int

const (
	Created PlanStatus = iota
	Running
	Canceled
)

// Plan 发布计划
type Plan struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	Name       string     `json:"name"`
	PlanStatus PlanStatus `json:"planStatus"`
	TeamId     int64      `json:"teamId"`
	Creator    string     `json:"creator"`
	Created    time.Time  `json:"created" xorm:"created"`
	Updated    time.Time  `json:"updated" xorm:"updated"`
}

func (*Plan) TableName() string {
	return PlanTableName
}

type TeamConfig struct {
	Id      int64                `json:"id" xorm:"pk autoincr"`
	TeamId  int64                `json:"teamId"`
	Content *deploy.NormalConfig `json:"content"`
	Created time.Time            `json:"created" xorm:"created"`
	Updated time.Time            `json:"updated" xorm:"updated"`
}

func (*TeamConfig) TableName() string {
	return TeamConfigTableName
}
