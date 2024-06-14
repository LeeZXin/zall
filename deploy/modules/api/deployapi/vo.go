package deployapi

import (
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
)

type UpdateConfigReqVO struct {
	ConfigId int64  `json:"configId"`
	Name     string `json:"name"`
	Content  string `json:"content"`
}

type ConfigVO struct {
	Id      int64  `json:"id"`
	AppId   string `json:"appId"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Env     string `json:"env"`
}

type CreateConfigReqVO struct {
	AppId   string `json:"appId"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Env     string `json:"env"`
}

type CreatePlanReqVO struct {
	Name        string `json:"name"`
	TeamId      int64  `json:"teamId"`
	Env         string `json:"env"`
	ExpireHours int    `json:"expireHours"`
}

type AddDeployPlanServiceReqVO struct {
	PlanId             int64  `json:"planId"`
	ConfigId           int64  `json:"configId"`
	LastProductVersion string `json:"lastProductVersion"`
	CurrProductVersion string `json:"currProductVersion"`
}

type ListPlanReqVO struct {
	TeamId  int64  `json:"teamId"`
	PageNum int    `json:"pageNum"`
	Env     string `json:"env"`
}

type PlanVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	IsClosed bool   `json:"isClosed"`
	TeamId   int64  `json:"teamId"`
	Creator  string `json:"creator"`
	Expired  string `json:"expired"`
	Created  string `json:"created"`
}

type PlanServiceVO struct {
	Id                 int64                  `json:"id"`
	AppId              string                 `json:"appId"`
	ConfigId           int64                  `json:"configId"`
	ConfigName         string                 `json:"configName"`
	CurrProductVersion string                 `json:"currProductVersion"`
	LastProductVersion string                 `json:"lastProductVersion"`
	ServiceStatus      deploymd.ServiceStatus `json:"serviceStatus"`
	Created            string                 `json:"created"`
}
