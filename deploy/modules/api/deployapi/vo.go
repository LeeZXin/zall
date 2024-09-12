package deployapi

import (
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
)

type CreatePlanReqVO struct {
	Name            string `json:"name"`
	PipelineId      int64  `json:"pipelineId"`
	ArtifactVersion string `json:"artifactVersion"`
}

type ListPlanReqVO struct {
	AppId   string `json:"appId"`
	PageNum int    `json:"pageNum"`
	Env     string `json:"env"`
}

type PlanVO struct {
	Id              int64               `json:"id"`
	PipelineId      int64               `json:"pipelineId"`
	PipelineName    string              `json:"pipelineName"`
	Name            string              `json:"name"`
	ArtifactVersion string              `json:"artifactVersion"`
	PlanStatus      deploymd.PlanStatus `json:"planStatus"`
	Env             string              `json:"env"`
	Creator         util.User           `json:"creator"`
	Created         string              `json:"created"`
}

type PlanDetailVO struct {
	Id              int64               `json:"id"`
	PipelineId      int64               `json:"pipelineId"`
	PipelineName    string              `json:"pipelineName"`
	PipelineConfig  string              `json:"pipelineConfig"`
	Name            string              `json:"name"`
	ArtifactVersion string              `json:"artifactVersion"`
	PlanStatus      deploymd.PlanStatus `json:"planStatus"`
	Env             string              `json:"env"`
	Creator         util.User           `json:"creator"`
	Created         string              `json:"created"`
}

type CreatePipelineReqVO struct {
	AppId  string `json:"appId"`
	Name   string `json:"name"`
	Config string `json:"config"`
	Env    string `json:"env"`
}

type UpdatePipelineReqVO struct {
	PipelineId int64  `json:"pipelineId"`
	Name       string `json:"name"`
	Config     string `json:"config"`
}

type PipelineVO struct {
	Id     int64  `json:"id"`
	AppId  string `json:"appId"`
	Config string `json:"config"`
	Env    string `json:"env"`
	Name   string `json:"name"`
}

type SimplePipelineVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Env  string `json:"env"`
}

type SubStageVO struct {
	Id          int64                `json:"id"`
	Agent       string               `json:"agent"`
	AgentHost   string               `json:"agentHost"`
	StageStatus deploymd.StageStatus `json:"stageStatus"`
	ExecuteLog  string               `json:"executeLog"`
}

type StageVO struct {
	Name                             string          `json:"name"`
	Percent                          float64         `json:"percent"`
	Total                            int             `json:"total"`
	Done                             int             `json:"done"`
	IsAutomatic                      bool            `json:"isAutomatic"`
	HasError                         bool            `json:"hasError"`
	IsRunning                        bool            `json:"isRunning"`
	IsAllDone                        bool            `json:"isAllDone"`
	WaitInteract                     bool            `json:"waitInteract"`
	SubStages                        []SubStageVO    `json:"subStages"`
	Script                           string          `json:"script"`
	Confirm                          *deploy.Confirm `json:"confirm"`
	CanForceRedoUnSuccessAgentStages bool            `json:"canForceRedoUnSuccessAgentStages"`
}

type ConfirmInteractStageReqVO struct {
	PlanId     int64             `json:"planId"`
	StageIndex int               `json:"stageIndex"`
	Args       map[string]string `json:"args"`
}

type ForceRedoStageReqVO struct {
	PlanId     int64 `json:"planId"`
	StageIndex int   `json:"stageIndex"`
}

type CreateServiceSourceReqVO struct {
	Env    string `json:"env"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	ApiKey string `json:"apiKey"`
}

type UpdateServiceSourceReqVO struct {
	SourceId int64  `json:"sourceId"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	ApiKey   string `json:"apiKey"`
}

type ServiceSourceVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Env     string `json:"env"`
	Host    string `json:"host"`
	ApiKey  string `json:"apiKey"`
	Created string `json:"created"`
}

type SimpleServiceSourceVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SimpleBindServiceSourceVO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	BindId int64  `json:"bindId"`
	Env    string `json:"env"`
}

type PipelineVarsWithoutContentVO struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type PipelineVarsVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	AppId   string `json:"appId"`
	Env     string `json:"env"`
	Content string `json:"content"`
}

type CreatePipelineVarsReqVO struct {
	AppId   string `json:"appId"`
	Env     string `json:"env"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type UpdatePipelineVarsReqVO struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
}

type DoServiceStatusActionReqVO struct {
	BindId    int64  `json:"bindId"`
	ServiceId string `json:"serviceId"`
	Action    string `json:"action"`
}

type BindAppAndServiceSourceReqVO struct {
	AppId        string  `json:"appId"`
	SourceIdList []int64 `json:"sourceIdList"`
	Env          string  `json:"env"`
}
