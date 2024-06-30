package deploymd

import "github.com/LeeZXin/zall/pkg/deploy"

type InsertPlanReqDTO struct {
	Name           string
	PlanStatus     PlanStatus
	AppId          string
	ServiceId      int64
	ProductVersion string
	Creator        string
	Env            string
	ServiceConfig  string
}

type InsertDeployLogReqDTO struct {
	ConfigId       int64
	AppId          string
	PlanId         int64
	ServiceConfig  string
	ProductVersion string
	Env            string
	DeployOutput   string
	Operator       string
}

type InsertOpLogReqDTO struct {
	ConfigId       int64
	Operator       string
	ScriptOutput   string
	Env            string
	Op             Op
	ProductVersion string
}

type ListDeployLogReqDTO struct {
	ConfigId int64
	Cursor   int64
	Limit    int
	Env      string
}

type ListOpLogReqDTO struct {
	ConfigId int64
	Cursor   int64
	Limit    int
	Env      string
}

type ListPlanReqDTO struct {
	AppId    string
	PageNum  int
	PageSize int
	Env      string
}

type InsertDeployStageReqDTO struct {
	PlanId      int64
	StageIndex  int
	Agent       string
	TaskId      string
	StageStatus StageStatus
}

type InsertServiceReqDTO struct {
	AppId       string
	Name        string
	Config      string
	Env         string
	ServiceType deploy.ServiceType
}

type UpdateServiceReqDTO struct {
	ServiceId   int64
	Name        string
	Config      string
	ServiceType deploy.ServiceType
}

type ListServiceReqDTO struct {
	AppId string
	Env   string
	Cols  []string
}
