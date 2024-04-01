package deploymd

import (
	"github.com/LeeZXin/zall/pkg/deploy"
	"time"
)

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
	StartTime          int64
}

type UpdateServiceReqDTO struct {
	ConfigId           int64
	CurrProductVersion string
	LastProductVersion string
	ServiceConfig      string
	Env                string
	ActiveStatus       ActiveStatus
	StartTime          int64
	ProbeTime          int64
}

type InsertPlanReqDTO struct {
	Name       string
	PlanStatus PlanStatus
	PlanType   PlanType
	TeamId     int64
	Creator    string
	Env        string
	Expired    time.Time
}

type InsertDeployLogReqDTO struct {
	ConfigId       int64
	AppId          string
	PlanId         int64
	ServiceType    deploy.ServiceType
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

type InsertPlanApprovalReqDTO struct {
	Name        string
	TeamId      int64
	Env         string
	DeployItems DeployItems
	Creator     string
}

type InsertApprovalNotifyReqDTO struct {
	Aid          int64
	Account      string
	NotifyStatus NotifyStatus
}

type InsertPlanItemReqDTO struct {
	PlanId             int64
	ConfigId           int64
	LastProductVersion string
	ProductVersion     string
	ItemStatus         PlanItemStatus
}

type ListPlanReqDTO struct {
	TeamId int64
	Cursor int64
	Limit  int
	Env    string
}
