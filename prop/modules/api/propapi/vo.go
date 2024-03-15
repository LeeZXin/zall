package propapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InsertEtcdNodeReqVO struct {
	NodeId    string   `json:"nodeId"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type DeleteEtcdNodeReqVO struct {
	NodeId string `json:"nodeId"`
	Env    string `json:"env"`
}

type UpdateEtcdNodeReqVO struct {
	NodeId    string   `json:"nodeId"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type EtcdNodeVO struct {
	NodeId    string   `json:"nodeId"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type ListEtcdNodeRespVO struct {
	ginutil.BaseResp
	Data []EtcdNodeVO
}

type ListSimpleEtcdNodeRespVO struct {
	ginutil.BaseResp
	Data []string
}

type InsertContentReqVO struct {
	AppId   string `json:"appId"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Env     string `json:"env"`
}

type UpdateContentReqVO struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
	Env     string `json:"env"`
}

type DeleteContentReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type ListContentReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type DeployContentReqVO struct {
	Id           int64    `json:"id"`
	Version      string   `json:"version"`
	EtcdNodeList []string `json:"etcdNodeList"`
	Env          string   `json:"env"`
}

type ListHistoryReqVO struct {
	ContentId int64  `json:"contentId"`
	Version   string `json:"version"`
	Cursor    int64  `json:"cursor"`
	Limit     int    `json:"limit"`
	Env       string `json:"env"`
}

type ListDeployReqVO struct {
	ContentId int64  `json:"contentId"`
	NodeId    string `json:"nodeId"`
	Version   string `json:"version"`
	Cursor    int64  `json:"cursor"`
	Limit     int    `json:"limit"`
	Env       string `json:"env"`
}

type PropContentVO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
}

type ListContentRespVO struct {
	ginutil.BaseResp
	Data []PropContentVO `json:"data"`
}

type HistoryVO struct {
	ContentId int64  `json:"contentId"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	Created   string `json:"created"`
	Creator   string `json:"creator"`
}

type ListHistoryRespVO struct {
	ginutil.BaseResp
	Cursor int64       `json:"cursor"`
	Data   []HistoryVO `json:"data"`
}

type DeployVO struct {
	ContentId int64  `json:"contentId"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	NodeId    string `json:"nodeId"`
	Created   string `json:"created"`
	Creator   string `json:"creator"`
}

type ListDeployRespVO struct {
	ginutil.BaseResp
	Cursor int64      `json:"cursor"`
	Data   []DeployVO `json:"data"`
}

type GrantAuthReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type GetAuthReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type GetAuthRespVO struct {
	ginutil.BaseResp
	Username string `json:"username"`
	Password string `json:"password"`
}

type ListSimpleEtcdNodeReqVO struct {
	Env string `json:"env"`
}

type ListEtcdNodeReqVO struct {
	Env string `json:"env"`
}
