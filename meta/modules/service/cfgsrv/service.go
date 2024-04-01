package cfgsrv

import (
	"context"
	"github.com/LeeZXin/zall/util"
)

var (
	Inner InnerService = &innerImpl{
		cfgCache: util.NewGoCache(),
	}
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	GetSysCfg(context.Context) (SysCfg, bool)
	InitSysCfg()
	GetGitCfg(context.Context) (GitCfg, bool)
	InitGitCfg()
	GetEnvCfg(context.Context) ([]string, bool)
	InitEnvCfg()
	ContainsEnv(string) bool
}

type OuterService interface {
	GetSysCfg(context.Context, GetSysCfgReqDTO) (SysCfg, error)
	UpdateSysCfg(context.Context, UpdateSysCfgReqDTO) error
	GetGitCfg(context.Context, GetGitCfgReqDTO) (GitCfg, error)
	UpdateGitCfg(context.Context, UpdateGitCfgReqDTO) error
	GetEnvCfg(context.Context, GetEnvCfgReqDTO) ([]string, error)
	UpdateEnvCfg(context.Context, UpdateEnvCfgReqDTO) error
}
