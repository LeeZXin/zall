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

type GetSysCfgReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetSysCfgReqDTO) IsValid() error {
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

type GetGitCfgReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetGitCfgReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
