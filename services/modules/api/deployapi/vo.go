package deployapi

import (
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
)

type ListConfigReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type UpdateConfigReqVO struct {
	ConfigId      int64                 `json:"configId"`
	Name          string                `json:"name"`
	Env           string                `json:"env"`
	ProcessConfig *deploy.ProcessConfig `json:"processConfig"`
	K8sConfig     *deploy.K8sConfig     `json:"k8sConfig"`
}

type ConfigVO struct {
	Id            int64                 `json:"id"`
	AppId         string                `json:"appId"`
	Name          string                `json:"name"`
	ServiceType   string                `json:"serviceType"`
	ProcessConfig *deploy.ProcessConfig `json:"processConfig,omitempty"`
	K8sConfig     *deploy.K8sConfig     `json:"k8sConfig,omitempty"`
	Created       string                `json:"created"`
}

type InsertConfigReqVO struct {
	AppId         string                `json:"appId"`
	Name          string                `json:"name"`
	ServiceType   deploy.ServiceType    `json:"serviceType"`
	ProcessConfig *deploy.ProcessConfig `json:"processConfig,omitempty"`
	K8sConfig     *deploy.K8sConfig     `json:"k8SConfig,omitempty"`
	Env           string                `json:"env"`
}

type InsertPlanReqVO struct {
	Name        string                `json:"name"`
	TeamId      int64                 `json:"teamId"`
	Env         string                `json:"env"`
	PlanType    deploymd.PlanType     `json:"planType"`
	DeployItems []deploymd.DeployItem `json:"deployItems"`
	ExpireHours int                   `json:"expireHours"`
}

type ClosePlanReqVO struct {
	PlanId int64  `json:"planId"`
	Env    string `json:"env"`
}

type InsertPlanItemReqVO struct {
	PlanId      int64                 `json:"planId"`
	DeployItems []deploymd.DeployItem `json:"deployItems"`
	Env         string                `json:"env"`
}

type ClosePlanItemReqVO struct {
	ItemId int64  `json:"itemId"`
	Env    string `json:"env"`
}

type ListPlanItemReqVO struct {
	PlanId int64  `json:"planId"`
	Env    string `json:"env"`
}

type DeployServiceWithoutPlanReqVO struct {
	ConfigId       int64  `json:"configId"`
	Env            string `json:"env"`
	ProductVersion string `json:"productVersion"`
	Operator       string `json:"operator"`
	AppId          string `json:"appId"`
}

type DeployServiceReqVO struct {
	ConfigId       int64  `json:"configId"`
	Env            string `json:"env"`
	ProductVersion string `json:"productVersion"`
}

type DeployServiceWithPlanReqVO struct {
	ItemId int64  `json:"itemId"`
	Env    string `json:"env"`
}

type RollbackServiceWithPlanReqVO struct {
	ItemId int64  `json:"itemId"`
	Env    string `json:"env"`
}

type StopServiceReqVO struct {
	ConfigId int64  `json:"configId"`
	Env      string `json:"env"`
}

type RestartServiceReqVO struct {
	ConfigId int64  `json:"configId"`
	Env      string `json:"env"`
}

type ListServiceReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type ServiceVO struct {
	CurrProductVersion string                `json:"currProductVersion"`
	LastProductVersion string                `json:"lastProductVersion"`
	ServiceType        string                `json:"serviceType"`
	ProcessConfig      *deploy.ProcessConfig `json:"processConfig,omitempty"`
	K8sConfig          *deploy.K8sConfig     `json:"k8SConfig,omitempty"`
	ActiveStatus       string                `json:"activeStatus"`
	StartTime          string                `json:"startTime"`
	ProbeTime          string                `json:"probeTime"`
	Created            string                `json:"created"`
}

type ListDeployLogReqVO struct {
	ConfigId int64  `json:"configId"`
	Cursor   int64  `json:"cursor"`
	Limit    int    `json:"limit"`
	Env      string `json:"env"`
}

type DeployLogVO struct {
	ServiceType    string `json:"serviceType"`
	ServiceConfig  string `json:"serviceConfig"`
	ProductVersion string `json:"productVersion"`
	Operator       string `json:"operator"`
	DeployOutput   string `json:"deployOutput"`
	Created        string `json:"created"`
	PlanId         int64  `json:"planId"`
}

type ListOpLogReqVO struct {
	ConfigId int64  `json:"configId"`
	Cursor   int64  `json:"cursor"`
	Limit    int    `json:"limit"`
	Env      string `json:"env"`
}

type OpLogVO struct {
	Op             string `json:"op"`
	Operator       string `json:"operator"`
	ScriptOutput   string `json:"scriptOutput"`
	ProductVersion string `json:"productVersion"`
	Created        string `json:"created"`
}

type ListPlanReqVO struct {
	TeamId int64  `json:"teamId"`
	Cursor int64  `json:"cursor"`
	Limit  int    `json:"limit"`
	Env    string `json:"env"`
}

type PlanVO struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	PlanType   string `json:"planType"`
	PlanStatus string `json:"planStatus"`
	TeamId     int64  `json:"teamId"`
	Creator    string `json:"creator"`
	Expired    string `json:"expired"`
	Created    string `json:"created"`
}

type PlanItemVO struct {
	Id                 int64  `json:"id"`
	AppId              string `json:"appId"`
	ConfigId           int64  `json:"configId"`
	ConfigName         string `json:"configName"`
	ProductVersion     string `json:"productVersion"`
	LastProductVersion string `json:"lastProductVersion"`
	ItemStatus         string `json:"itemStatus"`
	Created            string `json:"created"`
}
