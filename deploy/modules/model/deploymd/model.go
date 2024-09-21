package deploymd

import (
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

const (
	StageTableName         = "zservice_deploy_stage"
	PlanTableName          = "zservice_deploy_plan"
	PipelineTableName      = "zservice_pipeline"
	PipelineVarsTableName  = "zservice_pipeline_vars"
	SourceTableName        = "zservice_source"
	AppSourceBindTableName = "zservice_app_source_bind"
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

// Plan 发布计划
type Plan struct {
	Id              int64      `json:"id" xorm:"pk autoincr"`
	AppId           string     `json:"appId"`
	PipelineId      int64      `json:"pipelineId"`
	PipelineName    string     `json:"pipelineName"`
	Name            string     `json:"name"`
	ArtifactVersion string     `json:"artifactVersion"`
	PlanStatus      PlanStatus `json:"planStatus"`
	Env             string     `json:"env"`
	Creator         string     `json:"creator"`
	PipelineConfig  string     `json:"pipelineConfig"`
	Created         time.Time  `json:"created" xorm:"created"`
}

func (*Plan) TableName() string {
	return PlanTableName
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
	AppId       string                                  `json:"appId"`
	Agent       string                                  `json:"agent"`
	AgentHost   string                                  `json:"agentHost"`
	AgentToken  string                                  `json:"agentToken"`
	InputArgs   *xormutil.Conversion[map[string]string] `json:"inputArgs"`
	StageIndex  int                                     `json:"stageIndex"`
	ExecuteLog  string                                  `json:"executeLog"`
	StageStatus StageStatus                             `json:"stageStatus"`
	Script      string                                  `json:"script"`
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
	Id         int64     `json:"id" xorm:"pk autoincr"`
	Name       string    `json:"name"`
	Env        string    `json:"env"`
	Datasource string    `json:"datasource"`
	Created    time.Time `json:"created" xorm:"created"`
	Updated    time.Time `json:"updated" xorm:"updated"`
}

func (*ServiceSource) TableName() string {
	return SourceTableName
}

type AppServiceSourceBind struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	SourceId int64     `json:"sourceId"`
	AppId    string    `json:"appId"`
	Env      string    `json:"env"`
	Created  time.Time `json:"created" xorm:"created"`
}

func (*AppServiceSourceBind) TableName() string {
	return AppSourceBindTableName
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
