package deploymd

import "github.com/LeeZXin/zall/pkg/deploy"

type InsertConfigReqDTO struct {
	AppId       string
	Name        string
	ServiceType deploy.ServiceType
	Content     string
	Env         string
}

type UpdateConfigReqDTO struct {
	ConfigId int64
	Name     string
	Content  string
	Env      string
}

type InsertServiceReqDTO struct {
	ConfigId           int64
	CurrProductVersion string
	ServiceType        deploy.ServiceType
	ServiceConfig      string
	Env                string
	ActiveStatus       ActiveStatus
}

type UpdateServiceReqDTO struct {
	ConfigId           int64
	CurrProductVersion string
	LastProductVersion string
	ServiceConfig      string
	Env                string
}

type InsertPlanReqDTO struct {
	Name       string
	PlanStatus PlanStatus
	TeamId     int64
	Creator    string
	Env        string
}

type InsertLogReqDTO struct {
	ConfigId       int64
	AppId          string
	ServiceType    deploy.ServiceType
	ServiceConfig  string
	ProductVersion string
	Env            string
	DeployOutput   string
	Operator       string
}
