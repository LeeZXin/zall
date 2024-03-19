package gitactionapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type InsertActionReqVO struct {
	Id            int64  `json:"id"`
	NodeId        int64  `json:"nodeId"`
	ActionContent string `json:"actionContent"`
}

type DeleteActionReqVO struct {
	Id int64 `json:"id"`
}

type ListActionReqVO struct {
	Id int64 `json:"id"`
}

type ListActionRespVO struct {
	ginutil.BaseResp
	Data []ActionVO `json:"data"`
}

type UpdateActionReqVO struct {
	Id            int64  `json:"id"`
	NodeId        int64  `json:"nodeId"`
	ActionContent string `json:"actionContent"`
}

type ActionVO struct {
	Id            int64  `json:"id"`
	NodeId        int64  `json:"nodeId"`
	ActionContent string `json:"actionContent"`
	Created       string `json:"created"`
}

type TriggerActionReqVO struct {
	Id   int64             `json:"id"`
	Args map[string]string `json:"args"`
}

type InsertNodeReqVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
}

type UpdateNodeReqVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
}

type DeleteNodeReqVO struct {
	Id int64 `json:"id"`
}

type NodeVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
}

type ListGitNodeRespVO struct {
	ginutil.BaseResp
	Data []NodeVO `json:"data"`
}
