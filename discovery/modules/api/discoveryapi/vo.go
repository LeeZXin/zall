package discoveryapi

import "github.com/LeeZXin/zsf/services/lb"

type ListDiscoverySourceReqVO struct {
	Env string `json:"env"`
}

type ListBindDiscoverySourceReqVO struct {
	AppId string `json:"appId"`
	Env   string `json:"env"`
}

type DiscoverySourceVO struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type CreateDiscoverySourceReqVO struct {
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type UpdateDiscoverySourceReqVO struct {
	SourceId  int64    `json:"sourceId"`
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
}

type SimpleDiscoverySourceVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ServiceVO struct {
	lb.Server
	InstanceId string `json:"instanceId"`
}

type MarkAsDownServiceReqVO struct {
	BindId     int64  `json:"bindId"`
	InstanceId string `json:"instanceId"`
}

type MarkAsUpServiceReqVO struct {
	BindId     int64  `json:"bindId"`
	InstanceId string `json:"instanceId"`
}

type BindAppAndDiscoverySourceReqVO struct {
	AppId        string  `json:"appId"`
	SourceIdList []int64 `json:"sourceIdList"`
	Env          string  `json:"env"`
}

type SimpleBindDiscoverySourceVO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	BindId int64  `json:"bindId"`
	Env    string `json:"env"`
}
