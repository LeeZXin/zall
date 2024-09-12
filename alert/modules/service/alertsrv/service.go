package alertsrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppAlertConfigTopic, func(data any) {
			req, ok := data.(event.AppAlertConfigEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppAlertConfigCreateAction:
							return cfg.AppAlertConfig.Create
						case event.AppAlertConfigUpdateAction:
							return cfg.AppAlertConfig.Update
						case event.AppAlertConfigDeleteAction:
							return cfg.AppAlertConfig.Delete
						case event.AppAlertConfigEnableAction:
							return cfg.AppAlertConfig.Enable
						case event.AppAlertConfigDisableAction:
							return cfg.AppAlertConfig.Disable
						default:
							return false
						}
					}
					return false
				})
			}
		})
	})
}

// CreateConfig 新增配置
func CreateConfig(ctx context.Context, reqDTO CreateConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	app, team, err := checkAppDevelopPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator)
	if err != nil {
		return err
	}
	var cfg alertmd.Config
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		cfg, err2 = alertmd.InsertConfig(ctx, alertmd.InsertConfigReqDTO{
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
	notifyEvent(team, app, cfg, reqDTO.Operator, event.AppAlertConfigCreateAction)
	return nil
}

// UpdateConfig 修改配置
func UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	cfg, app, team, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
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
	notifyEvent(team, app, cfg, reqDTO.Operator, event.AppAlertConfigUpdateAction)
	return nil
}

// DeleteConfig 删除配置
func DeleteConfig(ctx context.Context, reqDTO DeleteConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	cfg, app, team, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
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
	notifyEvent(team, app, cfg, reqDTO.Operator, event.AppAlertConfigDeleteAction)
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
	accounts := listutil.MapNe(configs, func(t alertmd.Config) string {
		return t.Creator
	})
	userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, accounts...)
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
			Creator:     userMap[t.Creator],
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
	cfg, app, team, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
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
	notifyEvent(team, app, cfg, reqDTO.Operator, event.AppAlertConfigEnableAction)
	return nil
}

func DisableConfig(ctx context.Context, reqDTO DisableConfigReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	cfg, app, team, err := checkAppDevelopPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator)
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
	notifyEvent(team, app, cfg, reqDTO.Operator, event.AppAlertConfigDisableAction)
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

func notifyEvent(team teammd.Team, app appmd.App, cfg alertmd.Config, operator apisession.UserInfo, action event.AppAlertConfigEventAction) {
	initPsub()
	psub.Publish(event.AppAlertConfigTopic, event.AppAlertConfigEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action: action,
		Name:   cfg.Name,
		Env:    cfg.Env,
	})
}
