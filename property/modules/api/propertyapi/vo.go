package propertyapi

import "github.com/LeeZXin/zall/util"

type CreatePropertySourceReqVO struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
	Name      string   `json:"name"`
}

type UpdatePropertySourceReqVO struct {
	SourceId  int64    `json:"sourceId"`
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type PropertySourceVO struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type SimplePropertySourceVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
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

type DeployHistoryReqVO struct {
	HistoryId    int64   `json:"historyId"`
	SourceIdList []int64 `json:"sourceIdList"`
}

type PageHistoryReqVO struct {
	FileId  int64 `json:"fileId"`
	PageNum int   `json:"pageNum"`
}

type FileVO struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"`
	Env   string `json:"env"`
}

type HistoryVO struct {
	Id          int64     `json:"id"`
	FileName    string    `json:"fileName"`
	FileId      int64     `json:"fileId"`
	Content     string    `json:"content"`
	Version     string    `json:"version"`
	Created     string    `json:"created"`
	Creator     util.User `json:"creator"`
	LastVersion string    `json:"lastVersion"`
	Env         string    `json:"env"`
}

type DeployVO struct {
	NodeName  string    `json:"nodeName"`
	Endpoints string    `json:"endpoints"`
	Created   string    `json:"created"`
	Creator   util.User `json:"creator"`
}

type BindAppAndPropertySourceReqVO struct {
	AppId        string  `json:"appId"`
	SourceIdList []int64 `json:"sourceIdList"`
	Env          string  `json:"env"`
}

type SearchFromSourceResult struct {
	Version string `json:"version"`
	Content string `json:"content"`
	Exist   bool   `json:"exist"`
}
