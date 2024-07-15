package cfgsrv

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type UpdateSysCfgReqDTO struct {
	SysCfg
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateSysCfgReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateGitCfgReqDTO struct {
	GitCfg
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateGitCfgReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateEnvCfgReqDTO struct {
	Envs     []string            `json:"envs"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateEnvCfgReqDTO) IsValid() error {
	if len(r.Envs) == 0 {
		return util.InvalidArgsError()
	}
	for _, env := range r.Envs {
		if env == "" {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetGitCfgReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetGitCfgReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetEnvCfgReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetEnvCfgReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetGitRepoServerUrlReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetGitRepoServerUrlReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateGitRepoServerCfgReqDTO struct {
	GitRepoServerCfg
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateGitRepoServerCfgReqDTO) IsValid() error {
	if !util.GenIpPortPattern().MatchString(r.HttpHost) {
		return util.InvalidArgsError()
	}
	if !util.GenIpPortPattern().MatchString(r.SshHost) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetZonesCfgReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetZonesCfgReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateZonesCfgReqDTO struct {
	ZonesCfg
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateZonesCfgReqDTO) IsValid() error {
	for _, zone := range r.Zones {
		if !validZonesRegexp.MatchString(zone) {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
