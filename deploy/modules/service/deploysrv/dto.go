package deploysrv

import (
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
	"gopkg.in/yaml.v3"
	"time"
)

type CreatePlanReqDTO struct {
	Name           string              `json:"name"`
	PipelineId     int64               `json:"pipelineId"`
	ProductVersion string              `json:"productVersion"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *CreatePlanReqDTO) IsValid() error {
	if !deploymd.IsProductVersionValid(r.ProductVersion) {
		return util.InvalidArgsError()
	}
	if !deploymd.IsPlanNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.PipelineId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ClosePlanReqDTO struct {
	PlanId   int64               `json:"planId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ClosePlanReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDeployLogReqDTO struct {
	ConfigId int64               `json:"configId"`
	Env      string              `json:"env"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDeployLogReqDTO) IsValid() error {
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListOpLogReqDTO struct {
	ConfigId int64               `json:"configId"`
	Env      string              `json:"env"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListOpLogReqDTO) IsValid() error {
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Cursor < 0 || r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	return nil
}

type OpLogDTO struct {
	Op             deploymd.Op
	Operator       string
	ScriptOutput   string
	ProductVersion string
	Created        time.Time
}

type ListPlanReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPlanReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetPlanDetailReqDTO struct {
	PlanId   int64               `json:"planId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetPlanDetailReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListStagesReqDTO struct {
	PlanId   int64               `json:"planId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListStagesReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type KillStageReqDTO struct {
	PlanId     int64               `json:"planId"`
	StageIndex int                 `json:"stageIndex"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *KillStageReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if r.StageIndex < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ConfirmInteractStageReqDTO struct {
	PlanId     int64               `json:"planId"`
	StageIndex int                 `json:"stageIndex"`
	Args       map[string]string   `json:"args"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ConfirmInteractStageReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if r.StageIndex < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RedoAgentStageReqDTO struct {
	StageId  int64               `json:"stageId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RedoAgentStageReqDTO) IsValid() error {
	if r.StageId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ForceRedoNotSuccessfulAgentStagesReqDTO struct {
	PlanId     int64               `json:"planId"`
	StageIndex int                 `json:"stageIndex"`
	Args       map[string]string   `json:"args"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ForceRedoNotSuccessfulAgentStagesReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if r.StageIndex < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PlanDTO struct {
	Id             int64
	PipelineId     int64
	PipelineName   string
	Name           string
	ProductVersion string
	PlanStatus     deploymd.PlanStatus
	Env            string
	Creator        string
	Created        time.Time
}

type PlanDetailDTO struct {
	Id             int64
	PipelineId     int64
	PipelineName   string
	PipelineConfig string
	Name           string
	ProductVersion string
	PlanStatus     deploymd.PlanStatus
	Env            string
	Creator        string
	Created        time.Time
}

type StartPlanReqDTO struct {
	PlanId   int64               `json:"planId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *StartPlanReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreatePipelineReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Config   string              `json:"config"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
	pipeline deploy.Pipeline
}

func (r *CreatePipelineReqDTO) IsValid() error {
	if !deploymd.IsPipelineNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Config), &r.pipeline)
	if err != nil || !r.pipeline.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdatePipelineReqDTO struct {
	PipelineId int64               `json:"pipelineId"`
	Name       string              `json:"name"`
	Config     string              `json:"config"`
	Operator   apisession.UserInfo `json:"operator"`
	pipeline   deploy.Pipeline
}

func (r *UpdatePipelineReqDTO) IsValid() error {
	if r.PipelineId <= 0 {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Config), &r.pipeline)
	if err != nil || !r.pipeline.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeletePipelineReqDTO struct {
	PipelineId int64               `json:"pipelineId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *DeletePipelineReqDTO) IsValid() error {
	if r.PipelineId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPipelineReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPipelineReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPipelineWhenCreatePlanReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPipelineWhenCreatePlanReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PipelineDTO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	AppId  string `json:"appId"`
	Config string `json:"config"`
	Env    string `json:"env"`
}

type SimplePipelineDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Env  string `json:"env"`
}

type SubStageDTO struct {
	Id          int64
	Agent       string
	AgentHost   string
	StageStatus deploymd.StageStatus
	ExecuteLog  string
}

type StageDTO struct {
	Name                             string
	Percent                          float64
	Total                            int
	Done                             int
	IsAutomatic                      bool
	HasError                         bool
	IsRunning                        bool
	IsAllDone                        bool
	WaitInteract                     bool
	SubStages                        []SubStageDTO
	Script                           string
	Confirm                          *deploy.Confirm
	CanForceRedoUnSuccessAgentStages bool
}

type ListServiceSourceReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListServiceSourceReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ServiceSourceDTO struct {
	Id      int64
	Name    string
	AppId   string
	Env     string
	Hosts   []string
	ApiKey  string
	Created time.Time
}

type CreateServiceSourceReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Name     string              `json:"name"`
	Hosts    []string            `json:"hosts"`
	ApiKey   string              `json:"apiKey"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateServiceSourceReqDTO) IsValid() error {
	if !deploymd.IsServiceSourceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Hosts) == 0 {
		return util.InvalidArgsError()
	}
	for _, host := range r.Hosts {
		if !util.IpPortPattern.MatchString(host) {
			return util.InvalidArgsError()
		}
	}
	return nil
}

type UpdateServiceSourceReqDTO struct {
	SourceId int64               `json:"sourceId"`
	Name     string              `json:"name"`
	Hosts    []string            `json:"hosts"`
	ApiKey   string              `json:"apiKey"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateServiceSourceReqDTO) IsValid() error {
	if r.SourceId <= 0 {
		return util.InvalidArgsError()
	}
	if !deploymd.IsServiceSourceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Hosts) == 0 {
		return util.InvalidArgsError()
	}
	for _, host := range r.Hosts {
		if !util.IpPortPattern.MatchString(host) {
			return util.InvalidArgsError()
		}
	}
	return nil
}

type DeleteServiceSourceReqDTO struct {
	SourceId int64               `json:"sourceId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteServiceSourceReqDTO) IsValid() error {
	if r.SourceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
