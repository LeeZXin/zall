package cfgapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
)

type UpdateSysCfgReqVO struct {
	cfgsrv.SysCfg
}

type UpdateGitCfgReqVO struct {
	cfgsrv.GitCfg
}

type UpdateEnvCfgReqVO struct {
	Envs []string `json:"envs"`
}
