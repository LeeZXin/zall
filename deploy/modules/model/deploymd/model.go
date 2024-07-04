package deploymd

import (
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	StageTableName        = "zservice_deploy_stage"
	DeployLogTableName    = "zservice_deploy_log"
	PlanTableName         = "zservice_deploy_plan"
	OpLogTableName        = "zservice_op_log"
	PipelineTableName     = "zservice_pipeline"
	PipelineVarsTableName = "zservice_pipeline_vars"
	SourceTableName       = "zservice_source"
)

type PlanStatus int

const (
	PendingPlanStatus PlanStatus = iota + 1
	RunningPlanStatus
	SuccessfulPlanStatus
	ClosedPlanStatus
)

func (s PlanStatus) IsFinalStatus() bool {
	switch s {
	case SuccessfulPlanStatus, ClosedPlanStatus:
		return true
	default:
		return false
	}
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
	Id             int64      `json:"id" xorm:"pk autoincr"`
	AppId          string     `json:"appId"`
	PipelineId     int64      `json:"pipelineId"`
	PipelineName   string     `json:"pipelineName"`
	Name           string     `json:"name"`
	ProductVersion string     `json:"productVersion"`
	PlanStatus     PlanStatus `json:"planStatus"`
	Env            string     `json:"env"`
	Creator        string     `json:"creator"`
	PipelineConfig string     `json:"pipelineConfig"`
	Created        time.Time  `json:"created" xorm:"created"`
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
	SuccessfulStageStatus
	FailedStageStatus
)

type Stage struct {
	Id          int64                                   `json:"id" xorm:"pk autoincr"`
	PlanId      int64                                   `json:"planId"`
	Agent       string                                  `json:"agent"`
	InputArgs   *xormutil.Conversion[map[string]string] `json:"inputArgs"`
	StageIndex  int                                     `json:"stageIndex"`
	ExecuteLog  string                                  `json:"executeLog"`
	StageStatus StageStatus                             `json:"stageStatus"`
	TaskId      string                                  `json:"taskId"`
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

type Pipeline struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Name    string    `json:"name"`
	Config  string    `json:"config"`
	Env     string    `json:"env"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Pipeline) TableName() string {
	return PipelineTableName
}

type ServiceSource struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Name    string    `json:"name"`
	AppId   string    `json:"appId"`
	Env     string    `json:"env"`
	Host    string    `json:"host"`
	ApiKey  string    `json:"apiKey"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*ServiceSource) TableName() string {
	return SourceTableName
}

type PipelineVars struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Env     string    `json:"env"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*PipelineVars) TableName() string {
	return PipelineVarsTableName
}
