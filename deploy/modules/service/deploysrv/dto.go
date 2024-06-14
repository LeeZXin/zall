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

type ListConfigReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListConfigReqDTO) IsValid() error {
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

type UpdateConfigReqDTO struct {
	ConfigId int64               `json:"configId"`
	Name     string              `json:"name"`
	Content  string              `json:"content"`
	Operator apisession.UserInfo `json:"operator"`
	deploy   deploy.Deploy
}

func (r *UpdateConfigReqDTO) IsValid() error {
	if !deploymd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Content), &r.deploy)
	if err != nil || !r.deploy.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteConfigReqDTO struct {
	ConfigId int64               `json:"configId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteConfigReqDTO) IsValid() error {
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ConfigDTO struct {
	Id      int64  `json:"id"`
	AppId   string `json:"appId"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Env     string `json:"env"`
}

type CreateConfigReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Content  string              `json:"content"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
	deploy   deploy.Deploy
}

func (r *CreateConfigReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !deploymd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Content), &r.deploy)
	if err != nil || !r.deploy.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreatePlanReqDTO struct {
	Name        string              `json:"name"`
	TeamId      int64               `json:"teamId"`
	Env         string              `json:"env"`
	ExpireHours int                 `json:"expireHours"`
	Operator    apisession.UserInfo `json:"operator"`
}

func (r *CreatePlanReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !deploymd.IsPlanNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.ExpireHours <= 0 || r.ExpireHours > 240 {
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

type AddPlanServiceReqDTO struct {
	PlanId             int64               `json:"planId"`
	ConfigId           int64               `json:"configId"`
	LastProductVersion string              `json:"lastProductVersion"`
	CurrProductVersion string              `json:"currProductVersion"`
	Operator           apisession.UserInfo `json:"operator"`
}

func (r *AddPlanServiceReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if r.ConfigId <= 0 || r.CurrProductVersion == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeletePendingPlanServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DeletePendingPlanServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPlanServiceReqDTO struct {
	PlanId   int64               `json:"planId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPlanServiceReqDTO) IsValid() error {
	if r.PlanId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PlanServiceDTO struct {
	Id                 int64
	AppId              string
	ConfigId           int64
	ConfigName         string
	CurrProductVersion string
	LastProductVersion string
	ServiceStatus      deploymd.ServiceStatus
	Created            time.Time
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

type RestartServiceReqDTO struct {
	ConfigId int64               `json:"configId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RestartServiceReqDTO) IsValid() error {
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
	TeamId   int64               `json:"teamId"`
	Env      string              `json:"env"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListPlanReqDTO) IsValid() error {
	if r.TeamId <= 0 {
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

type PlanDTO struct {
	Id       int64
	Name     string
	TeamId   int64
	Creator  string
	IsClosed bool
	Expired  time.Time
	Created  time.Time
}

type StartPlanServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *StartPlanServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type FinishPlanServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *FinishPlanServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ConfirmPlanServiceStepReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Index     int                 `json:"index"`
	Input     map[string]string   `json:"input"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ConfirmPlanServiceStepReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Index < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RollbackPlanServiceStepReqDTO struct {
	StepId   int64               `json:"stepId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RollbackPlanServiceStepReqDTO) IsValid() error {
	if r.StepId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
