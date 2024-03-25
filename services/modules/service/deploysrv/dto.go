package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
)

type GetDeployReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetDeployReqDTO) IsValid() error {
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

type UpdateDeployReqDTO struct {
	AppId    string              `json:"appId"`
	Env      string              `json:"env"`
	Config   deploy.Config       `json:"config"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateDeployReqDTO) IsValid() error {
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
	if !r.Config.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
