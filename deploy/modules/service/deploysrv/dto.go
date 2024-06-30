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
	ServiceId      int64               `json:"serviceId"`
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
	if r.ServiceId <= 0 {
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

type DeployServiceWithoutPlanReqDTO struct {
	AppId    string `json:"appId"`
	Env      string `json:"env"`
	Product  string `json:"product"`
	Operator string `json:"operator"`
}

func (r *DeployServiceWithoutPlanReqDTO) IsValid() error {
	if r.Product == "" {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if r.Operator == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type DeployServiceReqDTO struct {
	ConfigId       int64               `json:"configId"`
	Env            string              `json:"env"`
	ProductVersion string              `json:"productVersion"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *DeployServiceReqDTO) IsValid() error {
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.ProductVersion == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type StopServiceReqDTO struct {
	ConfigId int64               `json:"configId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *StopServiceReqDTO) IsValid() error {
	if r.ConfigId <= 0 {
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

type DeployServiceWithPlanReqDTO struct {
	ItemId   int64               `json:"itemId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeployServiceWithPlanReqDTO) IsValid() error {
	if r.ItemId <= 0 {
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

type RollbackServiceWithPlanReqDTO struct {
	ItemId   int64               `json:"itemId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RollbackServiceWithPlanReqDTO) IsValid() error {
	if r.ItemId <= 0 {
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
	ServiceId      int64
	ServiceName    string
	Name           string
	ProductVersion string
	PlanStatus     deploymd.PlanStatus
	Env            string
	Creator        string
	Created        time.Time
}

type PlanDetailDTO struct {
	Id             int64
	ServiceId      int64
	ServiceName    string
	ServiceConfig  string
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

type CreateServiceReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Config   string              `json:"config"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
	service  deploy.Service
}

func (r *CreateServiceReqDTO) IsValid() error {
	if !deploymd.IsServiceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Config), &r.service)
	if err != nil || !r.service.IsValid() {
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

type UpdateServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Name      string              `json:"name"`
	Config    string              `json:"config"`
	Operator  apisession.UserInfo `json:"operator"`
	service   deploy.Service
}

func (r *UpdateServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Config), &r.service)
	if err != nil || !r.service.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DeleteServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListServiceReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListServiceReqDTO) IsValid() error {
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

type ListServiceWhenCreatePlanReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListServiceWhenCreatePlanReqDTO) IsValid() error {
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

type ServiceDTO struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	AppId       string             `json:"appId"`
	Config      string             `json:"config"`
	Env         string             `json:"env"`
	ServiceType deploy.ServiceType `json:"serviceType"`
}

type SimpleServiceDTO struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	Env         string             `json:"env"`
	ServiceType deploy.ServiceType `json:"serviceType"`
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
