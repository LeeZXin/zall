package cfgsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/meta/modules/model/cfgmd"
	"github.com/LeeZXin/zall/pkg/login"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

func InitSysCfg() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret := SysCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.Fatalf("init sys config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, &SysCfg{
			DisableSelfRegisterUser: false,
			AllowUserCreateTeam:     true,
		})
		if err != nil {
			logger.Logger.Fatalf("init sys config with err: %v", err)
		}
	}
}

func InitLoginCfg() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret := LoginCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.Fatalf("init login config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, &LoginCfg{
			Cfg: login.Cfg{
				AccountPassword: login.AccountPassword{
					IsEnabled: true,
				},
			},
		})
		if err != nil {
			logger.Logger.Fatalf("init login config with err: %v", err)
		}
	}
}

func InitGitCfg() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret := GitCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.Fatalf("init sys config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, &GitCfg{
			LfsJwtExpiry: 3600,
			LfsJwtSecret: idutil.RandomUuid(),
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Fatalf("init sys config with err: %v", err)
		}
	}
}

func GetGitCfgFromDB(ctx context.Context) (GitCfg, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var cfg GitCfg
	err := getFromDB(ctx, &cfg)
	return cfg, err
}

func GetEnvCfgFromDB(ctx context.Context) ([]string, error) {
	var cfg EnvCfg
	err := getFromDB(ctx, &cfg)
	return cfg.Envs, err
}

func InitEnvCfg() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret := EnvCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Fatalf("init env config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, &EnvCfg{
			Envs: []string{
				"prd",
			},
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Fatalf("init env config with err: %v", err)
		}
	}
}

func ContainsEnv(env string) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	envs, _ := GetEnvCfgFromDB(ctx)
	for _, s2 := range envs {
		if s2 == env {
			return true
		}
	}
	return false
}

// GetGitRepoServerCfgFromDB 获取git服务器地址
func GetGitRepoServerCfgFromDB(ctx context.Context) (GitRepoServerCfg, error) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var cfg GitRepoServerCfg
	err := getFromDB(ctx, &cfg)
	return cfg, err
}

func getFromDB(ctx context.Context, cfg util.KeyVal) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := cfgmd.GetByKey(ctx, cfg)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("%s not found", cfg.Key())
	}
	return nil
}

// GetSysCfg 获取系统全局配置
func GetSysCfg(ctx context.Context) (SysCfg, error) {
	ret := SysCfg{}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SysCfg{}, err
	}
	if !b {
		return *DefaultSysCfg, nil
	}
	return ret, nil
}

// GetLoginCfgFromDB 获取登录配置
func GetLoginCfgFromDB(ctx context.Context) (LoginCfg, error) {
	ret := LoginCfg{}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return LoginCfg{}, err
	}
	if !b {
		return LoginCfg{}, nil
	}
	return ret, nil
}

// UpdateSysCfg 编辑系统配置
func UpdateSysCfg(ctx context.Context, reqDTO UpdateSysCfgReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = cfgmd.UpdateByKey(ctx, &reqDTO.SysCfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func GetGitCfg(ctx context.Context, reqDTO GetGitCfgReqDTO) (GitCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return GitCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return GitCfg{}, util.UnauthorizedError()
	}
	cfg, err := GetGitCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return GitCfg{}, util.InternalError(err)
	}
	return cfg, nil
}

func UpdateGitCfg(ctx context.Context, reqDTO UpdateGitCfgReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := cfgmd.UpdateByKey(ctx, &reqDTO.GitCfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetEnvCfg 所有人都可以获取 不校验权限 获取环境列表
func GetEnvCfg(ctx context.Context, reqDTO GetEnvCfgReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	cfg, err := GetEnvCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return cfg, nil
}

// UpdateEnvCfg 编辑环境配置
func UpdateEnvCfg(ctx context.Context, reqDTO UpdateEnvCfgReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := cfgmd.UpdateByKey(ctx, &reqDTO.EnvCfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetGitRepoServerCfg 获取git服务器地址
func GetGitRepoServerCfg(ctx context.Context, reqDTO GetGitRepoServerUrlReqDTO) (GitRepoServerCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return GitRepoServerCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return GitRepoServerCfg{}, util.UnauthorizedError()
	}
	cfg, err := GetGitRepoServerCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return GitRepoServerCfg{}, util.InternalError(err)
	}
	return cfg, nil
}

// UpdateGitRepoServerCfg 更新git服务器地址
func UpdateGitRepoServerCfg(ctx context.Context, reqDTO UpdateGitRepoServerCfgReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var b bool
	b, err := cfgmd.ExistByKey(ctx, reqDTO.GitRepoServerCfg.Key())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		_, err = cfgmd.UpdateByKey(ctx, &reqDTO.GitRepoServerCfg)
	} else {
		err = cfgmd.InsertCfg(ctx, &reqDTO.GitRepoServerCfg)
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetLoginCfg 获取登录配置
func GetLoginCfg(ctx context.Context) (LoginCfg, error) {
	cfg, err := GetLoginCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return LoginCfg{}, util.InternalError(err)
	}
	cfg.Cfg.Wework.EraseSensitiveVar()
	cfg.Cfg.Feishu.EraseSensitiveVar()
	return cfg, nil
}

// GetLoginCfgBySa 超级管理员获取登录配置
func GetLoginCfgBySa(ctx context.Context, reqDTO GetLoginCfgBySaReqDTO) (LoginCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return LoginCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return LoginCfg{}, util.UnauthorizedError()
	}
	cfg, err := GetLoginCfgFromDB(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return LoginCfg{}, util.InternalError(err)
	}
	return cfg, nil
}

// UpdateLoginCfg 更新登录
func UpdateLoginCfg(ctx context.Context, reqDTO UpdateLoginCfgReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var b bool
	b, err := cfgmd.ExistByKey(ctx, reqDTO.LoginCfg.Key())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		_, err = cfgmd.UpdateByKey(ctx, &reqDTO.LoginCfg)
	} else {
		err = cfgmd.InsertCfg(ctx, &reqDTO.LoginCfg)
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}
