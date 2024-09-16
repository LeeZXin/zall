package discoverysrv

import (
	"github.com/LeeZXin/zall/discovery/modules/model/discoverymd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/services/lb"
)

type CreateDiscoverySourceReqDTO struct {
	Name      string              `json:"name"`
	Endpoints []string            `json:"endpoints"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Env       string              `json:"env"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *CreateDiscoverySourceReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if len(r.Endpoints) == 0 {
		return util.InvalidArgsError()
	}
	for _, endpoint := range r.Endpoints {
		if !util.GenIpPortPattern().MatchString(endpoint) {
			return util.InvalidArgsError()
		}
	}
	if !discoverymd.IsDiscoverySourceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteDiscoverySourceReqDTO struct {
	SourceId int64               `json:"sourceId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteDiscoverySourceReqDTO) IsValid() error {
	if r.SourceId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateDiscoverySourceReqDTO struct {
	SourceId  int64               `json:"sourceId"`
	Name      string              `json:"name"`
	Endpoints []string            `json:"endpoints"`
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateDiscoverySourceReqDTO) IsValid() error {
	if r.SourceId <= 0 {
		return util.InvalidArgsError()
	}
	if !discoverymd.IsDiscoverySourceNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if len(r.Endpoints) == 0 {
		return util.InvalidArgsError()
	}
	for _, endpoint := range r.Endpoints {
		if !util.GenIpPortPattern().MatchString(endpoint) {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDiscoverySourceReqDTO struct {
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDiscoverySourceReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAllDiscoverySourceReqDTO struct {
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllDiscoverySourceReqDTO) IsValid() error {
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListBindDiscoverySourceReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListBindDiscoverySourceReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DiscoverySourceDTO struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Env       string   `json:"env"`
}

type SimpleDiscoverySourceDTO struct {
	Id   int64
	Name string
}

type SimpleBindDiscoverySourceDTO struct {
	Id     int64
	Name   string
	BindId int64
	Env    string
}

type ListDiscoveryServiceReqDTO struct {
	BindId   int64               `json:"bindId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDiscoveryServiceReqDTO) IsValid() error {
	if r.BindId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ServiceDTO struct {
	lb.Server
	InstanceId string
}

type MarkAsDownServiceReqDTO struct {
	BindId     int64               `json:"bindId"`
	InstanceId string              `json:"instanceId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *MarkAsDownServiceReqDTO) IsValid() error {
	if r.BindId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.InstanceId) == 0 || len(r.InstanceId) > 32 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type MarkAsUpServiceReqDTO struct {
	BindId     int64               `json:"bindId"`
	InstanceId string              `json:"instanceId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *MarkAsUpServiceReqDTO) IsValid() error {
	if r.BindId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.InstanceId) == 0 || len(r.InstanceId) > 32 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type BindAppAndDiscoverySourceReqDTO struct {
	AppId        string              `json:"appId"`
	SourceIdList []int64             `json:"sourceIdList"`
	Env          string              `json:"env"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *BindAppAndDiscoverySourceReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	for _, i := range r.SourceIdList {
		if i <= 0 {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !cfgsrv.ContainsEnv(r.Env) {
		return util.InvalidArgsError()
	}
	return nil
}
