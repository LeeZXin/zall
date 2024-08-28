package alertsrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

// CreateConfig 新增配置
func CreateConfig(ctx context.Context, reqDTO CreateConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		cfg, err2 := alertmd.InsertConfig(ctx, alertmd.InsertConfigReqDTO{
			Name:        reqDTO.Name,
			Alert:       reqDTO.Alert,
			AppId:       reqDTO.AppId,
			IntervalSec: reqDTO.IntervalSec,
			IsEnabled:   false,
			Env:         reqDTO.Env,
			Creator:     reqDTO.Operator.Account,
		})
		if err2 != nil {
			return err2
		}
		return alertmd.InsertExecute(ctx, alertmd.InsertExecuteReqDTO{
			ConfigId:  cfg.Id,
			IsEnabled: false,
			NextTime:  time.Now().Add(time.Duration(reqDTO.IntervalSec) * time.Second).UnixMilli(),
			Env:       reqDTO.Env,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateConfig 修改配置
func UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = alertmd.UpdateConfig(ctx, alertmd.UpdateConfigReqDTO{
		Id:          reqDTO.Id,
		Name:        reqDTO.Name,
		Alert:       reqDTO.Alert,
		IntervalSec: reqDTO.IntervalSec,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteConfig 删除配置
func DeleteConfig(ctx context.Context, reqDTO DeleteConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := alertmd.DeleteConfigById(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		_, err2 = alertmd.DeleteExecuteByConfigId(ctx, reqDTO.Id)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListConfig 展示配置
func ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]ConfigDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return nil, 0, err
	}
	configs, total, err := alertmd.ListConfig(ctx, alertmd.ListConfigReqDTO{
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		AppId:    reqDTO.AppId,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data := listutil.MapNe(configs, func(t alertmd.Config) ConfigDTO {
		return ConfigDTO{
			Id:          t.Id,
			Name:        t.Name,
			AppId:       t.AppId,
			Content:     t.GetContent(),
			IntervalSec: t.IntervalSec,
			IsEnabled:   t.IsEnabled,
			Creator:     t.Creator,
			Env:         t.Env,
		}
	})
	return data, total, nil
}

func EnableConfig(ctx context.Context, reqDTO EnableConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err2 := alertmd.EnableExecute(ctx, reqDTO.Id, time.Now().UnixMilli())
		if err2 != nil {
			return err2
		}
		if b {
			b, err2 = alertmd.EnableConfig(ctx, reqDTO.Id)
			if err2 != nil {
				return err2
			}
			if !b {
				return errors.New("failed")
			}
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func DisableConfig(ctx context.Context, reqDTO DisableConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, _, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err2 := alertmd.DisableConfig(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		if b {
			b, err2 = alertmd.DisableExecute(ctx, reqDTO.Id)
			if err2 != nil {
				return err2
			}
			if !b {
				return errors.New("failed")
			}
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkAppDevelopPermByConfigId(ctx context.Context, id int64, operator apisession.UserInfo) (alertmd.Config, appmd.App, teammd.Team, error) {
	cfg, b, err := alertmd.GetConfigById(ctx, id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return alertmd.Config{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return alertmd.Config{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkAppDevelopPermByAppId(ctx, cfg.AppId, operator)
	return cfg, app, team, err
}

func checkAppDevelopPermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return app, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return app, team, util.InternalError(err)
	}
	if !b {
		return app, team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.GetAppPerm(appId).CanDevelop {
		return app, team, nil
	}
	return app, team, util.UnauthorizedError()
}
