package deploymd

import (
	"time"
)

type InsertConfigReqDTO struct {
	AppId   string
	Name    string
	Content string
	Env     string
}

type UpdateConfigReqDTO struct {
	ConfigId int64
	Name     string
	Content  string
}

type InsertPlanServiceReqDTO struct {
	ConfigId           int64
	CurrProductVersion string
	LastProductVersion string
	DeployConfig       string
	Status             ServiceStatus
	PlanId             int64
}

type UpdateServiceReqDTO struct {
	ConfigId           int64
	CurrProductVersion string
	LastProductVersion string
	ServiceConfig      string
	Env                string
	ActiveStatus       ServiceStatus
	StartTime          int64
	ProbeTime          int64
}

type InsertPlanReqDTO struct {
	Name     string
	IsClosed bool
	TeamId   int64
	Creator  string
	Env      string
	Expired  time.Time
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
	TeamId   int64
	PageNum  int
	PageSize int
	Env      string
}

type InsertDeployStepReqDTO struct {
	ServiceId  int64
	StepIndex  int
	Agent      string
	StepStatus StepStatus
}
