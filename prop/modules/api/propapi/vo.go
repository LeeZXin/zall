package propapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InsertEtcdNodeReqVO struct {
	NodeId    string   `json:"nodeId"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type DeleteEtcdNodeReqVO struct {
	NodeId string `json:"nodeId"`
}

type UpdateEtcdNodeReqVO struct {
	NodeId    string   `json:"nodeId"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
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
}

type UpdateContentReqVO struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
}

type DeleteContentReqVO struct {
	Id int64 `json:"id"`
}

type ListContentReqVO struct {
	AppId string `json:"appId"`
}

type DeployContentReqVO struct {
	Id           int64    `json:"id"`
	Version      string   `json:"version"`
	EtcdNodeList []string `json:"etcdNodeList"`
}

type ListHistoryReqVO struct {
	ContentId int64  `json:"contentId"`
	Version   string `json:"version"`
	Cursor    int64  `json:"cursor"`
	Limit     int    `json:"limit"`
}

type ListDeployReqVO struct {
	ContentId int64  `json:"contentId"`
	NodeId    string `json:"nodeId"`
	Version   string `json:"version"`
	Cursor    int64  `json:"cursor"`
	Limit     int    `json:"limit"`
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
}

type ListDeployRespVO struct {
	ginutil.BaseResp
	Cursor int64      `json:"cursor"`
	Data   []DeployVO `json:"data"`
}

type GrantAuthReqVO struct {
	AppId string `json:"appId"`
}
