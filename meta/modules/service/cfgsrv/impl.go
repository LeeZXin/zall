package cfgsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/cfgmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/patrickmn/go-cache"
	"time"
)

type innerImpl struct {
	cfgCache *cache.Cache
}

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

func (s *innerImpl) GetSysCfg(ctx context.Context) (SysCfg, bool) {
	cfg := new(SysCfg)
	v, b := s.cfgCache.Get(cfg.Key())
	if b {
		return v.(SysCfg), true
	}
	b = getFromDB(ctx, cfg)
	if b {
		s.cfgCache.Set(cfg.Key(), *cfg, time.Minute)
	}
	return *cfg, b
}

func (s *innerImpl) InitGitCfg() {
	ctx := context.Background()
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret := GitCfg{}
	b, err := cfgmd.GetByKey(ctx, &ret)
	if err != nil {
		logger.Logger.WithContext(ctx).Fatalf("init sys config with err: %v", err)
	}
	if !b {
		err = cfgmd.InsertCfg(ctx, DefaultGitCfg)
		if err != nil {
			logger.Logger.WithContext(ctx).Fatalf("init sys config with err: %v", err)
		}
	}
}

func (s *innerImpl) GetGitCfg(ctx context.Context) (GitCfg, bool) {
	cfg := new(GitCfg)
	v, b := s.cfgCache.Get(cfg.Key())
	if b {
		return v.(GitCfg), true
	}
	b = getFromDB(ctx, cfg)
	if b {
		s.cfgCache.Set(cfg.Key(), *cfg, time.Minute)
	}
	return *cfg, b
}

func (s *innerImpl) GetEnvCfg(ctx context.Context) ([]string, bool) {
	cfg := new(EnvCfg)
	v, b := s.cfgCache.Get(cfg.Key())
	if b {
		return v.(EnvCfg).Envs, true
	}
	b = getFromDB(ctx, cfg)
	if b {
		s.cfgCache.Set(cfg.Key(), *cfg, time.Minute)
	}
	return cfg.Envs, b
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

func getFromDB(ctx context.Context, cfg util.KeyVal) bool {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	b, err := cfgmd.GetByKey(ctx, cfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false
	}
	return b
}

type outerImpl struct{}

func (s *outerImpl) GetSysCfg(ctx context.Context, reqDTO GetSysCfgReqDTO) (SysCfg, error) {
	if err := reqDTO.IsValid(); err != nil {
		return SysCfg{}, err
	}
	if !reqDTO.Operator.IsAdmin {
		return SysCfg{}, util.UnauthorizedError()
	}
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

func (s *outerImpl) UpdateSysCfg(ctx context.Context, reqDTO UpdateSysCfgReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.CfgSrvKeysVO.UpdateSysCfg),
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

// GetEnvCfg 所有人都可以获取 不校验权限
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
	return envConfig.Envs, nil
}

func (s *outerImpl) UpdateEnvCfg(ctx context.Context, reqDTO UpdateEnvCfgReqDTO) (err error) {
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
	cfg := EnvCfg{
		Envs: reqDTO.Envs,
	}
	_, err = cfgmd.UpdateByKey(ctx, &cfg)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}
