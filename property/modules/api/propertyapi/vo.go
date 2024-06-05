package propertyapi

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

type CreateFileReqVO struct {
	AppId   string `json:"appId"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Env     string `json:"env"`
}

type NewVersionReqVO struct {
	FileId      int64  `json:"fileId"`
	Content     string `json:"content"`
	LastVersion string `json:"lastVersion"`
}

type DeleteContentReqVO struct {
	Id  int64  `json:"id"`
	Env string `json:"env"`
}

type DeployContentReqVO struct {
	Id           int64    `json:"id"`
	Version      string   `json:"version"`
	EtcdNodeList []string `json:"etcdNodeList"`
	Env          string   `json:"env"`
}

type PageHistoryReqVO struct {
	FileId  int64 `json:"fileId"`
	PageNum int   `json:"pageNum"`
}

type ListDeployReqVO struct {
	ContentId int64  `json:"contentId"`
	NodeId    string `json:"nodeId"`
	Version   string `json:"version"`
	Cursor    int64  `json:"cursor"`
	Limit     int    `json:"limit"`
	Env       string `json:"env"`
}

type FileVO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
	Env   string `json:"env"`
}

type HistoryVO struct {
	Id          int64  `json:"id"`
	FileId      int64  `json:"fileId"`
	Content     string `json:"content"`
	Version     string `json:"version"`
	Created     string `json:"created"`
	Creator     string `json:"creator"`
	LastVersion string `json:"lastVersion"`
}

type DeployVO struct {
	ContentId int64  `json:"contentId"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	NodeId    string `json:"nodeId"`
	Created   string `json:"created"`
	Creator   string `json:"creator"`
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
