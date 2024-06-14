package deploymd

import (
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	StepTableName        = "zservice_deploy_step"
	ConfigTableName      = "zservice_deploy_config"
	PlanServiceTableName = "zservice_deploy_plan_service"
	DeployLogTableName   = "zservice_deploy_log"
	PlanTableName        = "zservice_deploy_plan"
	OpLogTableName       = "zservice_op_log"
)

// Config 部署配置
type Config struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Env     string    `json:"env"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Config) TableName() string {
	return ConfigTableName
}

type ServiceStatus int

const (
	PendingServiceStatus ServiceStatus = iota + 1
	RunningServiceStatus
	FinishServiceStatus
)

// PlanService 部署服务
type PlanService struct {
	Id int64 `json:"id" xorm:"pk autoincr"`
	// 发布计划id
	PlanId   int64 `json:"planId"`
	ConfigId int64 `json:"configId"`
	// 当前制品版本
	CurrProductVersion string `json:"productVersion"`
	// 上个制品版本
	LastProductVersion string        `json:"lastProductVersion"`
	DeployConfig       string        `json:"deployConfig"`
	ServiceStatus      ServiceStatus `json:"serviceStatus"`
	Created            time.Time     `json:"created" xorm:"created"`
}

func (*PlanService) TableName() string {
	return PlanServiceTableName
}

// DeployLog 部署日志
type DeployLog struct {
	Id       int64 `json:"id" xorm:"pk autoincr"`
	ConfigId int64 `json:"configId"`
	// 发布计划id
	PlanId         int64     `json:"planId"`
	AppId          string    `json:"appId"`
	ServiceConfig  string    `json:"serviceConfig"`
	ProductVersion string    `json:"productVersion"`
	Operator       string    `json:"operator"`
	DeployOutput   string    `json:"deployOutput"`
	Created        time.Time `json:"created" xorm:"created"`
}

func (*DeployLog) TableName() string {
	return DeployLogTableName
}

// Plan 发布计划
type Plan struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	Name     string    `json:"name"`
	IsClosed bool      `json:"isClosed"`
	TeamId   int64     `json:"teamId"`
	Creator  string    `json:"creator"`
	Env      string    `json:"env"`
	Expired  time.Time `json:"expired"`
	Created  time.Time `json:"created" xorm:"created"`
}

func (*Plan) TableName() string {
	return PlanTableName
}

func (p *Plan) IsExpired() bool {
	return p.Expired.Before(time.Now())
}

func (p *Plan) IsInvalid() bool {
	return p.IsClosed || p.IsExpired()
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

type StepStatus int

const (
	PendingStepStatus StepStatus = iota + 1
	RunningStepStatus
	SuccessStepStatus
	FailStepStatus
	RollbackStepStatus
)

type Step struct {
	Id          int64                                   `json:"id" xorm:"pk autoincr"`
	ServiceId   int64                                   `json:"serviceId"`
	Agent       string                                  `json:"agent"`
	InputArgs   *xormutil.Conversion[map[string]string] `json:"inputArgs"`
	StepIndex   int                                     `json:"stepIndex"`
	ExecuteLog  string                                  `json:"executeLog"`
	RollbackLog string                                  `json:"rollbackLog"`
	StepStatus  StepStatus                              `json:"stepStatus"`
	Created     time.Time                               `json:"created" xorm:"created"`
	Updated     time.Time                               `json:"updated" xorm:"updated"`
}

func (*Step) TableName() string {
	return StepTableName
}

func (s *Step) GetInputArgs() map[string]string {
	if s.InputArgs == nil || s.InputArgs.Data == nil {
		return make(map[string]string)
	}
	return s.InputArgs.Data
}
