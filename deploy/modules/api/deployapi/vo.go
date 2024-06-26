package deployapi

import (
	"github.com/LeeZXin/zall/deploy/modules/model/deploymd"
	"github.com/LeeZXin/zall/pkg/deploy"
)

type CreatePlanReqVO struct {
	Name           string `json:"name"`
	ServiceId      int64  `json:"serviceId"`
	ProductVersion string `json:"productVersion"`
}

type ListPlanReqVO struct {
	AppId   string `json:"appId"`
	PageNum int    `json:"pageNum"`
	Env     string `json:"env"`
}

type PlanVO struct {
	Id             int64               `json:"id"`
	ServiceId      int64               `json:"serviceId"`
	ServiceName    string              `json:"serviceName"`
	Name           string              `json:"name"`
	ProductVersion string              `json:"productVersion"`
	PlanStatus     deploymd.PlanStatus `json:"planStatus"`
	Env            string              `json:"env"`
	Creator        string              `json:"creator"`
	Created        string              `json:"created"`
}

type CreateServiceReqVO struct {
	AppId  string `json:"appId"`
	Name   string `json:"name"`
	Config string `json:"config"`
	Env    string `json:"env"`
}

type UpdateServiceReqVO struct {
	ServiceId int64  `json:"serviceId"`
	Name      string `json:"name"`
	Config    string `json:"config"`
}

type ServiceVO struct {
	Id          int64              `json:"id"`
	AppId       string             `json:"appId"`
	Config      string             `json:"config"`
	Env         string             `json:"env"`
	Name        string             `json:"name"`
	ServiceType deploy.ServiceType `json:"serviceType"`
}

type SimpleServiceVO struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	Env         string             `json:"env"`
}
