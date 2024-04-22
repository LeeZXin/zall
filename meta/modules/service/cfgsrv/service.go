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
	// GetGitRepoServerCfg 获取git服务器地址 从缓存中获取
	GetGitRepoServerCfg(context.Context) (GitRepoServerCfg, bool)
}

type OuterService interface {
	// GetSysCfg 获取系统全局配置
	GetSysCfg(context.Context) (SysCfg, error)
	UpdateSysCfg(context.Context, UpdateSysCfgReqDTO) error
	GetGitCfg(context.Context, GetGitCfgReqDTO) (GitCfg, error)
	UpdateGitCfg(context.Context, UpdateGitCfgReqDTO) error
	GetEnvCfg(context.Context, GetEnvCfgReqDTO) ([]string, error)
	UpdateEnvCfg(context.Context, UpdateEnvCfgReqDTO) error
	// GetGitRepoServerCfg 获取git服务器地址
	GetGitRepoServerCfg(context.Context, GetGitRepoServerUrlReqDTO) (GitRepoServerCfg, error)
	// UpdateGitRepoServerUrl 更新git服务器地址
	UpdateGitRepoServerUrl(context.Context, UpdateGitRepoServerUrlReqDTO) error
}
