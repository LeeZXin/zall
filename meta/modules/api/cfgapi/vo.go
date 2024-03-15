package cfgapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type GetSysCfgRespVO struct {
	ginutil.BaseResp
	Cfg cfgsrv.SysCfg `json:"cfg"`
}

type UpdateSysCfgReqVO struct {
	cfgsrv.SysCfg
}

type UpdateGitCfgReqVO struct {
	cfgsrv.GitCfg
}

type GetGitCfgRespVO struct {
	ginutil.BaseResp
	Cfg cfgsrv.GitCfg `json:"cfg"`
}

type UpdateEnvCfgReqVO struct {
	Envs []string `json:"envs"`
}

type GetEnvCfgRespVO struct {
	ginutil.BaseResp
	Cfg []string `json:"cfg"`
}
