package discoveryapi

import "github.com/LeeZXin/zsf/services/lb"

type ListDiscoverySourceReqVO struct {
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
	AppId     string   `json:"appId"`
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
	Up         bool   `json:"up"`
	InstanceId string `json:"instanceId"`
}

type DeregisterServiceReqVO struct {
	SourceId   int64  `json:"sourceId"`
	InstanceId string `json:"instanceId"`
}

type ReRegisterServiceReqVO struct {
	SourceId   int64  `json:"sourceId"`
	InstanceId string `json:"instanceId"`
}

type DeleteDownServiceReqVO struct {
	SourceId   int64  `json:"sourceId"`
	InstanceId string `json:"instanceId"`
}
