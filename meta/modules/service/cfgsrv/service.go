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
}

type OuterService interface {
	GetSysCfg(context.Context, GetSysCfgReqDTO) (SysCfg, error)
	UpdateSysCfg(context.Context, UpdateSysCfgReqDTO) error
	GetGitCfg(context.Context, GetGitCfgReqDTO) (GitCfg, error)
	UpdateGitCfg(context.Context, UpdateGitCfgReqDTO) error
}
