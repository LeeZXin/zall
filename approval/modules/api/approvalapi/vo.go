package approvalapi

import (
	"github.com/LeeZXin/zall/pkg/approval"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type AgreeApprovalReqVO struct {
	NotifyId int64 `json:"notifyId"`
}

type DisagreeApprovalReqVO struct {
	NotifyId int64 `json:"notifyId"`
}

type InsertCustomProcessReqVO struct {
	Pid     string              `json:"pid"`
	Name    string              `json:"name"`
	GroupId int64               `json:"groupId"`
	Process approval.ProcessCfg `json:"process"`
}

type UpdateCustomProcessReqVO struct {
	Id      int64               `json:"id"`
	Name    string              `json:"name"`
	GroupId int64               `json:"groupId"`
	Process approval.ProcessCfg `json:"process"`
}

type InsertCustomFlowReqVO struct {
	Pid   string        `json:"pid"`
	BizId string        `json:"bizId"`
	Kvs   []approval.Kv `json:"kvs"`
}

type InsertCustomFlowRespVO struct {
	ginutil.BaseResp
	ErrKeys []string `json:"errKeys"`
}

type CancelCustomFlowReqVO struct {
	FlowId int64 `json:"flowId"`
}
