package deploymd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	StageTableName         = "zservice_deploy_stage"
	DeployLogTableName     = "zservice_deploy_log"
	PlanTableName          = "zservice_deploy_plan"
	OpLogTableName         = "zservice_op_log"
	ServiceTableName       = "zservice_service"
	DeployServiceTableName = "zservice_deploy_service"
)

type PlanStatus int

const (
	PendingPlanStatus PlanStatus = iota + 1
	RunningPlanStatus
	ClosedPlanStatus
)

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
	Id             int64           `json:"id" xorm:"pk autoincr"`
	AppId          string          `json:"appId"`
	ServiceId      int64           `json:"serviceId"`
	Name           string          `json:"name"`
	ProductVersion string          `json:"productVersion"`
	PlanStatus     PlanStatus      `json:"planStatus"`
	Env            string          `json:"env"`
	Creator        string          `json:"creator"`
	ServiceConfig  *deploy.Service `json:"serviceConfig"`
	Created        time.Time       `json:"created" xorm:"created"`
}

func (*Plan) TableName() string {
	return PlanTableName
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

type StageStatus int

const (
	PendingStageStatus StageStatus = iota + 1
	RunningStageStatus
	SuccessStageStatus
	FailStageStatus
	RollbackStageStatus
)

type Stage struct {
	Id          int64                                   `json:"id" xorm:"pk autoincr"`
	PlanId      int64                                   `json:"planId"`
	Agent       string                                  `json:"agent"`
	InputArgs   *xormutil.Conversion[map[string]string] `json:"inputArgs"`
	StageIndex  int                                     `json:"stageIndex"`
	ExecuteLog  string                                  `json:"executeLog"`
	RollbackLog string                                  `json:"rollbackLog"`
	StageStatus StageStatus                             `json:"stageStatus"`
	Created     time.Time                               `json:"created" xorm:"created"`
	Updated     time.Time                               `json:"updated" xorm:"updated"`
}

func (*Stage) TableName() string {
	return StageTableName
}

func (s *Stage) GetInputArgs() map[string]string {
	if s.InputArgs == nil || s.InputArgs.Data == nil {
		return make(map[string]string)
	}
	return s.InputArgs.Data
}

type Service struct {
	Id          int64              `json:"id" xorm:"pk autoincr"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	AppId       string             `json:"appId"`
	Name        string             `json:"name"`
	Config      string             `json:"config"`
	Env         string             `json:"env"`
	Created     time.Time          `json:"created" xorm:"created"`
	Updated     time.Time          `json:"updated" xorm:"updated"`
}

func (*Service) TableName() string {
	return ServiceTableName
}

type K8sConfig struct {
	AgentHost       string `json:"agentHost"`
	AgentToken      string `json:"agentToken"`
	GetStatusScript string `json:"getStatusScript"`
}

type ProcessConfig struct {
	Host       string `json:"host"`
	AgentHost  string `json:"agentHost"`
	AgentToken string `json:"agentToken"`
}

type DeployServiceConfig struct {
	Type    deploy.ServiceType `json:"type"`
	Process *ProcessConfig     `json:"process,omitempty"`
	K8s     *K8sConfig         `json:"k8s,omitempty"`
}

func (c *DeployServiceConfig) FromDB(content []byte) error {
	if c == nil {
		*c = DeployServiceConfig{}
	}
	return json.Unmarshal(content, c)
}

func (c *DeployServiceConfig) ToDB() ([]byte, error) {
	return json.Marshal(c)
}

type DeployService struct {
	Id             int64                `json:"id" xorm:"pk autoincr"`
	PlanId         int64                `json:"planId"`
	ServiceId      int64                `json:"serviceId"`
	Config         *DeployServiceConfig `json:"process"`
	Probed         int64                `json:"probed"`
	IsPlanDone     bool                 `json:"isPlanDone"`
	ProbeFailTimes int                  `json:"probeFailTimes"`
	Created        time.Time            `json:"created" xorm:"created"`
}

func (*DeployService) TableName() string {
	return DeployServiceTableName
}
