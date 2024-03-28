package deploymd

import (
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	TeamConfigTableName    = "zservice_team_config"
	ConfigTableName        = "zservice_deploy_config"
	ServiceTableName       = "zservice_deploy_service"
	DeployLogTableName     = "zservice_deploy_log"
	PlanTableName          = "zservice_deploy_plan"
	ProbeInstanceTableName = "zservice_probe_instance"
	OpLogTableName         = "zservice_op_log"
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
	AbnormalStatus ActiveStatus = iota + 1
	StartingStatus
	StartedStatus
	StoppingStatus
	StoppedStatus
)

func (s ActiveStatus) Readable() string {
	switch s {
	case AbnormalStatus:
		return i18n.GetByKey(i18n.ServiceAbnormalStatus)
	case StartingStatus:
		return i18n.GetByKey(i18n.ServiceStartingStatus)
	case StartedStatus:
		return i18n.GetByKey(i18n.ServiceStartedStatus)
	case StoppingStatus:
		return i18n.GetByKey(i18n.ServiceStoppingStatus)
	case StoppedStatus:
		return i18n.GetByKey(i18n.ServiceStoppedStatus)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownStatus)
	}
}

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
	StartTime          int64              `json:"startTime"`
	ProbeTime          int64              `json:"probeTime"`
	Created            time.Time          `json:"created" xorm:"created"`
	Updated            time.Time          `json:"updated" xorm:"updated"`
}

func (*Service) TableName() string {
	return ServiceTableName
}

// DeployLog 部署日志
type DeployLog struct {
	Id       int64 `json:"id" xorm:"pk autoincr"`
	ConfigId int64 `json:"configId"`
	// 发布计划id
	PlanId         int64              `json:"planId"`
	AppId          string             `json:"appId"`
	ServiceType    deploy.ServiceType `json:"serviceType"`
	ServiceConfig  string             `json:"serviceConfig"`
	ProductVersion string             `json:"productVersion"`
	Operator       string             `json:"operator"`
	DeployOutput   string             `json:"deployOutput"`
	Created        time.Time          `json:"created" xorm:"created"`
}

func (*DeployLog) TableName() string {
	return DeployLogTableName
}

type PlanStatus int

const (
	CreatedPlanStatus PlanStatus = iota + 1
	RunningPlanStatus
	CanceledPlanStatus
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

type ProbeInstance struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	InstanceId    string    `json:"instanceId"`
	HeartbeatTime int64     `json:"heartbeatTime"`
	Created       time.Time `json:"created" xorm:"created"`
}

func (*ProbeInstance) TableName() string {
	return ProbeInstanceTableName
}

type Op int

const (
	RestartServiceOp Op = iota + 1
	StopServiceOp
)

func (o Op) Readable() string {
	switch o {
	case RestartServiceOp:
		return i18n.GetByKey(i18n.ServiceRestartOp)
	case StopServiceOp:
		return i18n.GetByKey(i18n.ServiceStopOp)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownOp)
	}
}

type OpLog struct {
	Id             int64     `json:"id" xorm:"pk autoincr"`
	ConfigId       int64     `json:"configId"`
	Op             Op        `json:"op"`
	Operator       string    `json:"operator"`
	ScriptOutput   string    `json:"scriptOutput"`
	ProductVersion string    `json:"productVersion"`
	Created        time.Time `json:"created" xorm:"created"`
}

func (*OpLog) TableName() string {
	return OpLogTableName
}
