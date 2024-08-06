package cfgsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/meta/modules/model/cfgmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type innerImpl struct{}

func (s *innerImpl) InitSysCfg() {
	ctx := context.Background()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret := SysCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Fatalf("init sys config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, DefaultSysCfg)
		if err != nil {
			logger.Logger.WithContext(ctx).Fatalf("init sys config with err: %v", err)
		}
	}
}

func (s *innerImpl) GetSysCfg(ctx context.Context) (SysCfg, error) {
	var cfg SysCfg
	err := getFromDB(ctx, &cfg)
	return cfg, err
}

func (s *innerImpl) InitGitCfg() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	ret := GitCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.Fatalf("init sys config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, DefaultGitCfg)
		if err != nil {
			logger.Logger.WithContext(ctx).Fatalf("init sys config with err: %v", err)
		}
	}
}

func (s *innerImpl) GetGitCfg() (GitCfg, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	var cfg GitCfg
	err := getFromDB(ctx, &cfg)
	return cfg, err
}

func (s *innerImpl) GetEnvCfg(ctx context.Context) ([]string, error) {
	var cfg EnvCfg
	err := getFromDB(ctx, &cfg)
	return cfg.Envs, err
}

func (s *innerImpl) InitEnvCfg() {
	ctx := context.Background()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret := EnvCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Fatalf("init env config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, DefaultEnvCfg)
		if err != nil {
			logger.Logger.WithContext(ctx).Fatalf("init env config with err: %v", err)
		}
	}
}

func (s *innerImpl) ContainsEnv(env string) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	envs, _ := s.GetEnvCfg(ctx)
	for _, s2 := range envs {
		if s2 == env {
			return true
		}
	}
	return false
}

// GetGitRepoServerCfg 获取git服务器地址 从缓存中获取
func (s *innerImpl) GetGitRepoServerCfg() (GitRepoServerCfg, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	var cfg GitRepoServerCfg
	err := getFromDB(ctx, &cfg)
	return cfg, err
}

func getFromDB(ctx context.Context, cfg util.KeyVal) error {
	b, err := cfgmd.GetByKey(ctx, cfg)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("%s not found", cfg.Key())
	}
	return nil
}

type outerImpl struct{}

// GetSysCfg 获取系统全局配置
func (s *outerImpl) GetSysCfg(ctx context.Context) (SysCfg, error) {
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

// UpdateSysCfg 编辑系统配置
func (s *outerImpl) UpdateSysCfg(ctx context.Context, reqDTO UpdateSysCfgReqDTO) (err error) {
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

func (s *outerImpl) GetGitCfg(ctx context.Context, reqDTO GetGitCfgReqDTO) (GitCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return GitCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return GitCfg{}, util.UnauthorizedError()
	}
	ret := GitCfg{}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return GitCfg{}, err
	}
	if !b {
		return *DefaultGitCfg, nil
	}
	return ret, nil
}

func (s *outerImpl) UpdateGitCfg(ctx context.Context, reqDTO UpdateGitCfgReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.CfgSrvKeysVO.UpdateGitCfg),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = cfgmd.UpdateByKey(ctx, &reqDTO.GitCfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// GetEnvCfg 所有人都可以获取 不校验权限 获取环境列表
func (s *outerImpl) GetEnvCfg(ctx context.Context, reqDTO GetEnvCfgReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var envConfig EnvCfg
	_, err := cfgmd.GetByKey(ctx, &envConfig)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := envConfig.Envs
	if ret == nil {
		ret = make([]string, 0)
	}
	return ret, nil
}

// UpdateEnvCfg 编辑环境配置
func (s *outerImpl) UpdateEnvCfg(ctx context.Context, reqDTO UpdateEnvCfgReqDTO) error {
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
func (s *outerImpl) GetGitRepoServerCfg(ctx context.Context, reqDTO GetGitRepoServerUrlReqDTO) (GitRepoServerCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return GitRepoServerCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return GitRepoServerCfg{}, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var cfg GitRepoServerCfg
	_, err := cfgmd.GetByKey(ctx, &cfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return GitRepoServerCfg{}, util.InternalError(err)
	}
	return cfg, nil
}

// UpdateGitRepoServerCfg 更新git服务器地址
func (s *outerImpl) UpdateGitRepoServerCfg(ctx context.Context, reqDTO UpdateGitRepoServerCfgReqDTO) error {
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
