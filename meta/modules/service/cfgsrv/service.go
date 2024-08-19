package cfgsrv

import (
	"context"
)

var (
	Inner InnerService
	Outer OuterService
)

func Init() {
	InitInner()
	InitOuter()
}

func InitInner() {
	if Inner == nil {
		Inner = new(innerImpl)
	}
}

func InitOuter() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

type InnerService interface {
	// GetSysCfg 获取系统配置
	GetSysCfg(context.Context) (SysCfg, error)
	// InitSysCfg 初始化系统配置
	InitSysCfg()
	// GetGitCfg 获取git配置
	GetGitCfg() (GitCfg, error)
	// InitGitCfg 初始化git配置
	InitGitCfg()
	// GetEnvCfg 获取环境配置
	GetEnvCfg(context.Context) ([]string, error)
	// InitEnvCfg 初始化环境配置
	InitEnvCfg()
	// ContainsEnv 检查是否包含环境
	ContainsEnv(string) bool
	// GetGitRepoServerCfg 获取git服务器地址
	GetGitRepoServerCfg() (GitRepoServerCfg, error)
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
	// 获取定时任务告警模板

}
