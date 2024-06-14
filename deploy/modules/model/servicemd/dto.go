package servicemd

import "github.com/LeeZXin/zall/pkg/deploy"

type InsertServiceReqDTO struct {
	AppId       string
	Name        string
	ServiceType deploy.ServiceType
	Config      string
	Env         string
	Probed      int64
	IsEnabled   bool
}

type UpdateServiceReqDTO struct {
	ServiceId   int64
	ServiceType deploy.ServiceType
	Name        string
	Config      string
}

type ListServiceReqDTO struct {
	AppId    string
	Env      string
	PageNum  int
	PageSize int
}
