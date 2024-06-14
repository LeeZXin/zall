package serviceapi

import "github.com/LeeZXin/zall/pkg/deploy"

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

type ListServiceReqVO struct {
	AppId   string `json:"appId"`
	Env     string `json:"env"`
	PageNum int    `json:"pageNum"`
}

type ServiceVO struct {
	Id          int64              `json:"id"`
	AppId       string             `json:"appId"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	Config      string             `json:"config"`
	Env         string             `json:"env"`
	IsEnabled   bool               `json:"isEnabled"`
	Name        string             `json:"name"`
}
