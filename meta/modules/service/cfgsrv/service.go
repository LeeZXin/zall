package cfgsrv

import (
	"context"
	"github.com/LeeZXin/zall/util"
)

var (
	Inner InnerService
	Outer OuterService
)

func Init() {
	if Inner == nil {
		Inner = &innerImpl{
			cfgCache: util.NewGoCache(),
		}
		Outer = new(outerImpl)
	}
}

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
	// UpdateSysCfg 编辑系统配置
	UpdateSysCfg(context.Context, UpdateSysCfgReqDTO) error
	// GetGitCfg 获取git配置
	GetGitCfg(context.Context, GetGitCfgReqDTO) (GitCfg, error)
	// UpdateGitCfg 编辑git配置
	UpdateGitCfg(context.Context, UpdateGitCfgReqDTO) error
	// GetEnvCfg 获取环境列表
	GetEnvCfg(context.Context, GetEnvCfgReqDTO) ([]string, error)
	// UpdateEnvCfg 编辑环境列表
	UpdateEnvCfg(context.Context, UpdateEnvCfgReqDTO) error
	// GetGitRepoServerCfg 获取git服务器地址
	GetGitRepoServerCfg(context.Context, GetGitRepoServerUrlReqDTO) (GitRepoServerCfg, error)
	// UpdateGitRepoServerCfg 更新git服务器地址
	UpdateGitRepoServerCfg(context.Context, UpdateGitRepoServerCfgReqDTO) error
}
