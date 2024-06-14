package servicesrv

import (
	"github.com/LeeZXin/zall/deploy/modules/model/servicemd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
	"gopkg.in/yaml.v3"
)

type CreateServiceReqDTO struct {
	AppId    string              `json:"appId"`
	Name     string              `json:"name"`
	Config   string              `json:"config"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
	service  deploy.Service
}

func (r *CreateServiceReqDTO) IsValid() error {
	if !servicemd.IsServiceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Config), &r.service)
	if err != nil || !r.service.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Name      string              `json:"name"`
	Config    string              `json:"config"`
	Operator  apisession.UserInfo `json:"operator"`
	service   deploy.Service
}

func (r *UpdateServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	err := yaml.Unmarshal([]byte(r.Config), &r.service)
	if err != nil || !r.service.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DeleteServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListServiceReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListServiceReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.Inner.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type EnableServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *EnableServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DisableServiceReqDTO struct {
	ServiceId int64               `json:"serviceId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DisableServiceReqDTO) IsValid() error {
	if r.ServiceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ServiceDTO struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	ServiceType deploy.ServiceType `json:"serviceType"`
	AppId       string             `json:"appId"`
	Config      string             `json:"config"`
	Env         string             `json:"env"`
	IsEnabled   bool               `json:"isEnabled"`
}
