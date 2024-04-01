package deploymd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	TeamConfigTableName     = "zservice_team_config"
	ConfigTableName         = "zservice_deploy_config"
	ServiceTableName        = "zservice_deploy_service"
	DeployLogTableName      = "zservice_deploy_log"
	PlanTableName           = "zservice_deploy_plan"
	PlanItemTableName       = "zservice_deploy_plan_item"
	ApprovalTableName       = "zservice_deploy_approval"
	ApprovalNotifyTableName = "zservice_deploy_approval_notify"
	ProbeInstanceTableName  = "zservice_probe_instance"
	OpLogTableName          = "zservice_op_log"
)

// Config 部署配置
type Config struct {
	Id          int64              `json:"id" xorm:"pk autoincr"`
	AppId       string             `json:"appId"`
	Name        string             `json:"name"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	Content     string             `json:"content"`
	Created     time.Time          `json:"created" xorm:"created"`
	Updated     time.Time          `json:"updated" xorm:"updated"`
}

func (*Config) TableName() string {
	return ConfigTableName
}

type ActiveStatus int

const (
	AbnormalStatus ActiveStatus = iota + 1
	StartingStatus
	StartedStatus
	StoppingStatus
	StoppedStatus
)

func (s ActiveStatus) Readable() string {
	switch s {
	case AbnormalStatus:
		return i18n.GetByKey(i18n.ServiceAbnormalStatus)
	case StartingStatus:
		return i18n.GetByKey(i18n.ServiceStartingStatus)
	case StartedStatus:
		return i18n.GetByKey(i18n.ServiceStartedStatus)
	case StoppingStatus:
		return i18n.GetByKey(i18n.ServiceStoppingStatus)
	case StoppedStatus:
		return i18n.GetByKey(i18n.ServiceStoppedStatus)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownStatus)
	}
}

// Service 部署服务
type Service struct {
	Id       int64 `json:"id" xorm:"pk autoincr"`
	ConfigId int64 `json:"configId"`
	// 当前制品版本
	CurrProductVersion string `json:"productVersion"`
	// 上个制品版本
	LastProductVersion string             `json:"lastProductVersion"`
	ServiceType        deploy.ServiceType `json:"serviceType"`
	ServiceConfig      string             `json:"serviceConfig"`
	ActiveStatus       ActiveStatus       `json:"activeStatus"`
	StartTime          int64              `json:"startTime"`
	ProbeTime          int64              `json:"probeTime"`
	Created            time.Time          `json:"created" xorm:"created"`
	Updated            time.Time          `json:"updated" xorm:"updated"`
}

func (*Service) TableName() string {
	return ServiceTableName
}

// DeployLog 部署日志
type DeployLog struct {
	Id       int64 `json:"id" xorm:"pk autoincr"`
	ConfigId int64 `json:"configId"`
	// 发布计划id
	PlanId         int64              `json:"planId"`
	AppId          string             `json:"appId"`
	ServiceType    deploy.ServiceType `json:"serviceType"`
	ServiceConfig  string             `json:"serviceConfig"`
	ProductVersion string             `json:"productVersion"`
	Operator       string             `json:"operator"`
	DeployOutput   string             `json:"deployOutput"`
	Created        time.Time          `json:"created" xorm:"created"`
}

func (*DeployLog) TableName() string {
	return DeployLogTableName
}

type PlanStatus int

const (
	RunningPlanStatus PlanStatus = iota + 1
	ClosedPlanStatus
)

func (t PlanStatus) Readable() string {
	switch t {
	case RunningPlanStatus:
		return i18n.GetByKey(i18n.ServiceRunningPlanStatus)
	case ClosedPlanStatus:
		return i18n.GetByKey(i18n.ServiceClosedPlanStatus)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownPlanStatus)
	}
}

type PlanType int

const (
	AddServiceAfterPlanCreatingType PlanType = iota + 1
	AddServiceBeforePlanCreatingType
)

func (t PlanType) IsValid() bool {
	switch t {
	case AddServiceAfterPlanCreatingType, AddServiceBeforePlanCreatingType:
		return true
	default:
		return false
	}
}

func (t PlanType) Readable() string {
	switch t {
	case AddServiceAfterPlanCreatingType:
		return i18n.GetByKey(i18n.ServiceAddServiceAfterPlanCreatingType)
	case AddServiceBeforePlanCreatingType:
		return i18n.GetByKey(i18n.ServiceAddServiceBeforePlanCreatingType)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownPlanType)
	}
}

// Plan 发布计划
type Plan struct {
	Id         int64      `json:"id" xorm:"pk autoincr"`
	Name       string     `json:"name"`
	PlanType   PlanType   `json:"planType"`
	PlanStatus PlanStatus `json:"planStatus"`
	TeamId     int64      `json:"teamId"`
	Creator    string     `json:"creator"`
	Expired    time.Time  `json:"expired"`
	Created    time.Time  `json:"created" xorm:"created"`
	Updated    time.Time  `json:"updated" xorm:"updated"`
}

func (*Plan) TableName() string {
	return PlanTableName
}

func (p *Plan) IsExpired() bool {
	return p.Expired.Before(time.Now())
}

func (p *Plan) IsClosed() bool {
	return p.PlanStatus == ClosedPlanStatus || p.IsExpired()
}

type PlanItemStatus int

const (
	WaitItemStatus PlanItemStatus = iota + 1
	DeployedItemStatus
	RollbackItemStatus
	ClosedItemStatus
)

func (s PlanItemStatus) Readable() string {
	switch s {
	case WaitItemStatus:
		return i18n.GetByKey(i18n.ServiceWaitPlanItemStatus)
	case DeployedItemStatus:
		return i18n.GetByKey(i18n.ServiceDeployedPlanItemStatus)
	case RollbackItemStatus:
		return i18n.GetByKey(i18n.ServiceRollbackPlanItemStatus)
	case ClosedItemStatus:
		return i18n.GetByKey(i18n.ServiceClosedPlanItemStatus)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownPlanItemStatus)
	}
}

type PlanItem struct {
	Id                 int64          `json:"id" xorm:"pk autoincr"`
	PlanId             int64          `json:"planId"`
	ConfigId           int64          `json:"configId"`
	ProductVersion     string         `json:"productVersion"`
	LastProductVersion string         `json:"lastProductVersion"`
	ItemStatus         PlanItemStatus `json:"itemStatus"`
	Created            time.Time      `json:"created" xorm:"created"`
	Updated            time.Time      `json:"updated" xorm:"updated"`
}

func (*PlanItem) TableName() string {
	return PlanItemTableName
}

type TeamConfig struct {
	Id      int64                `json:"id" xorm:"pk autoincr"`
	TeamId  int64                `json:"teamId"`
	Content *deploy.NormalConfig `json:"content"`
	Created time.Time            `json:"created" xorm:"created"`
	Updated time.Time            `json:"updated" xorm:"updated"`
}

func (*TeamConfig) TableName() string {
	return TeamConfigTableName
}

type ProbeInstance struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	InstanceId    string    `json:"instanceId"`
	HeartbeatTime int64     `json:"heartbeatTime"`
	Created       time.Time `json:"created" xorm:"created"`
}

func (*ProbeInstance) TableName() string {
	return ProbeInstanceTableName
}

type Op int

const (
	RestartServiceOp Op = iota + 1
	StopServiceOp
)

func (o Op) Readable() string {
	switch o {
	case RestartServiceOp:
		return i18n.GetByKey(i18n.ServiceRestartOp)
	case StopServiceOp:
		return i18n.GetByKey(i18n.ServiceStopOp)
	default:
		return i18n.GetByKey(i18n.ServiceUnknownOp)
	}
}

type OpLog struct {
	Id             int64     `json:"id" xorm:"pk autoincr"`
	ConfigId       int64     `json:"configId"`
	Op             Op        `json:"op"`
	Operator       string    `json:"operator"`
	ScriptOutput   string    `json:"scriptOutput"`
	ProductVersion string    `json:"productVersion"`
	Created        time.Time `json:"created" xorm:"created"`
}

func (*OpLog) TableName() string {
	return OpLogTableName
}

type DeployItem struct {
	ConfigId       int64  `json:"configId"`
	ProductVersion string `json:"productVersion"`
}

func (r *DeployItem) IsValid() bool {
	return r.ConfigId > 0 && len(r.ProductVersion) > 0 && len(r.ProductVersion) <= 64
}

type DeployItems []DeployItem

func (r *DeployItems) FromDB(content []byte) error {
	if r == nil {
		*r = make([]DeployItem, 0)
	}
	return json.Unmarshal(content, r)
}

func (r *DeployItems) ToDB() ([]byte, error) {
	return json.Marshal(r)
}

type PlanApproval struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	Name        string      `json:"name"`
	TeamId      int64       `json:"teamId"`
	DeployItems DeployItems `json:"deployItems"`
	Creator     string      `json:"creator"`
	Created     time.Time   `json:"created" xorm:"created"`
	Updated     time.Time   `json:"updated" xorm:"updated"`
}

func (*PlanApproval) TableName() string {
	return ApprovalTableName
}

type NotifyStatus int

const (
	WaitStatus NotifyStatus = iota + 1
	AgreeStatus
	DisagreeStatus
)

type ApprovalNotify struct {
	Id           int64        `json:"id" xorm:"pk autoincr"`
	Aid          int64        `json:"aid"`
	Account      string       `json:"account"`
	NotifyStatus NotifyStatus `json:"notifyStatus"`
	Created      time.Time    `json:"created" xorm:"created"`
	Updated      time.Time    `json:"updated" xorm:"updated"`
}

func (*ApprovalNotify) TableName() string {
	return ApprovalNotifyTableName
}
