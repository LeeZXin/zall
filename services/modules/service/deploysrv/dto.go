package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"time"
)

type ListConfigReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListConfigReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateConfigReqDTO struct {
	ConfigId      int64                 `json:"configId"`
	Name          string                `json:"name"`
	Env           string                `json:"env"`
	ProcessConfig *deploy.ProcessConfig `json:"processConfig"`
	K8sConfig     *deploy.K8sConfig     `json:"k8sConfig"`
	Operator      apisession.UserInfo   `json:"operator"`
}

func (r *UpdateConfigReqDTO) IsValid() error {
	if !deploymd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ConfigDTO struct {
	Id            int64                 `json:"id"`
	AppId         string                `json:"appId"`
	Name          string                `json:"name"`
	ServiceType   deploy.ServiceType    `json:"serviceType"`
	ProcessConfig *deploy.ProcessConfig `json:"processConfig"`
	K8sConfig     *deploy.K8sConfig     `json:"k8sConfig"`
	Created       time.Time             `json:"created"`
}

type InsertConfigReqDTO struct {
	AppId         string                `json:"appId"`
	Name          string                `json:"name"`
	ServiceType   deploy.ServiceType    `json:"serviceType"`
	ProcessConfig *deploy.ProcessConfig `json:"processConfig"`
	K8sConfig     *deploy.K8sConfig     `json:"k8sConfig"`
	Env           string                `json:"env"`
	Operator      apisession.UserInfo   `json:"operator"`
}

func (r *InsertConfigReqDTO) IsValid() error {
	if !appmd.IsAppIdValid(r.AppId) {
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !deploymd.IsConfigNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	switch r.ServiceType {
	case deploy.ProcessServiceType:
		if r.ProcessConfig == nil || !r.ProcessConfig.IsValid() {
			return util.InvalidArgsError()
		}
	case deploy.K8sServiceType:
		if r.K8sConfig == nil || !r.K8sConfig.IsValid() {
			return util.InvalidArgsError()
		}
	default:
		return util.InvalidArgsError()
	}
	return nil
}

type InsertPlanReqDTO struct {
	Name     string              `json:"name"`
	TeamId   int64               `json:"teamId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertPlanReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !deploymd.IsPlanNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeployServiceWithoutPlanReqDTO struct {
	ConfigId       int64  `json:"configId"`
	Env            string `json:"env"`
	ProductVersion string `json:"productVersion"`
	Operator       string `json:"operator"`
}

func (r *DeployServiceWithoutPlanReqDTO) IsValid() error {
	if r.ConfigId <= 0 {
		return util.InvalidArgsError()
	}
	if r.ProductVersion == "" {
		return util.InvalidArgsError()
	}
	envs, _ := cfgsrv.Inner.GetEnvCfg(context.Background())
	contains, _ := listutil.Contains(envs, func(t string) (bool, error) {
		return t == r.Env, nil
	})
	if !contains {
		return util.InvalidArgsError()
	}
	if r.Operator == "" {
		return util.InvalidArgsError()
	}
	return nil
}
