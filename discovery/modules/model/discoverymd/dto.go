package discoverymd

import "github.com/LeeZXin/zsf/services/lb"

type InsertEtcdNodeReqDTO struct {
	AppId     string
	Name      string
	Endpoints string
	Username  string
	Password  string
	Env       string
}

type UpdateEtcdNodeReqDTO struct {
	Id        int64
	Name      string
	Endpoints string
	Username  string
	Password  string
}

type ListEtcdNodeReqDTO struct {
	AppId string
	Env   string
	Cols  []string
}

type InsertDownServiceDTO struct {
	SourceId   int64
	AppId      string
	Service    lb.Server
	InstanceId string
}
