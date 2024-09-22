package alertsrv

import (
	"context"
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/notify/modules/service/notifysrv"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zall/pkg/commonhook"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/lease"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/PaesslerAG/gval"
	"github.com/prometheus/common/model"
	"github.com/spf13/cast"
	"time"
)

var (
	taskExecutor *executor.Executor
)

func InitTask() {
	poolSize := static.GetInt("alert.poolSize")
	if poolSize <= 0 {
		poolSize = 20
	}
	queueSize := static.GetInt("alert.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	taskExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	leaser, _ := lease.NewDbLease(
		"alert-lock",
		common.GetInstanceId(),
		"zall_lock",
		xormstore.GetEngine(),
		20*time.Second,
	)
	stopFunc, _ := taskutil.RunMainLoopTask(
		taskutil.MainLoopTaskOpts{
			Handler: func(ctx context.Context) {
				for ctx.Err() == nil {
					doExecuteTask(ctx)
					time.Sleep(5 * time.Second)
				}
			},
			Leaser: leaser,
			// 抢锁失败 空转等待时间
			WaitDuration: 30 * time.Second,
			// 锁过期时间有20秒 每8秒续命 至少2次续命成功的机会
			RenewDuration: 8 * time.Second,
			GrantCallback: func(err error, b bool) {
				if err != nil {
					logger.Logger.Errorf("alert task %s grant lease failed with err: %v", common.GetInstanceId(), err)
					return
				}
				if b {
					logger.Logger.Infof("alert task grant lease success: %v", common.GetInstanceId())
				}
			},
		},
	)
	quit.AddShutdownHook(quit.ShutdownHook(stopFunc), true)
}

func doExecuteTask(runCtx context.Context) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	err := alertmd.IterateExecute(ctx, time.Now().UnixMilli(), func(cfg *alertmd.Execute) error {
		if err := runCtx.Err(); err != nil {
			return err
		}
		return handleExecute(cfg)
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func handleExecute(execute *alertmd.Execute) error {
	return taskExecutor.Execute(func() {
		ctx, closer := xormstore.Context(context.Background())
		cfg, b, err := alertmd.GetConfigById(ctx, execute.ConfigId)
		if err == nil && b {
			interval := time.Duration(cfg.IntervalSec) * time.Second
			endTime := time.UnixMilli(execute.NextTime)
			if resetNextTime(ctx, &cfg, interval.Milliseconds()+execute.NextTime, execute.RunVersion) {
				closer.Close()
				alertCfg := cfg.GetContent()
				switch alertCfg.SourceType {
				case alert.MysqlType:
					result, err := alertCfg.Mysql.Execute(endTime)
					if err != nil {
						logger.Logger.Errorf("alert cfg: %v execute failed with err: %v", cfg.Id, err)
					} else if len(result) > 0 {
						evaluate, err := gval.Evaluate(alertCfg.Mysql.Condition, result)
						if err != nil {
							logger.Logger.Errorf("alert cfg: %v compile condition: %v failed with err: %v", cfg.Id, alertCfg.Mysql.Condition, err)
						} else if cast.ToBool(evaluate) {
							// 执行hook
							executeHook(cfg, result, "", endTime.Format(time.DateTime))
						}
					}
				case alert.PromType:
					result, err := alertCfg.Prom.Execute(endTime)
					if err != nil {
						logger.Logger.Errorf("alert cfg: %v execute failed with err: %v", cfg.Id, err)
					} else {
						var args map[string]any
						if result != nil {
							args = metric2Map(result)
						} else {
							args = map[string]any{
								"metric_value": 0,
							}
						}
						evaluate, err := gval.Evaluate(alertCfg.Prom.Condition, args)
						if err != nil {
							logger.Logger.Errorf("alert cfg: %v compile condition: %v failed with err: %v", cfg.Id, alertCfg.Prom.Condition, err)
						} else if cast.ToBool(evaluate) {
							// 执行hook
							executeHook(cfg, args, "", endTime.Format(time.DateTime))
						}
					}
				case alert.LokiType:
					startTime, result, err := alertCfg.Loki.Execute(endTime)
					if err != nil {
						logger.Logger.Errorf("alert cfg: %v execute failed with err: %v", cfg.Id, err)
					} else {
						params := map[string]any{
							"Value": result,
						}
						evaluate, err := gval.Evaluate(alertCfg.Loki.Condition, params)
						if err != nil {
							logger.Logger.Errorf("alert cfg: %v compile condition: %v failed with err: %v", cfg.Id, alertCfg.Loki.Condition, err)
						} else if cast.ToBool(evaluate) {
							// 执行hook
							executeHook(cfg, params, startTime.Format(time.DateTime), endTime.Format(time.DateTime))
						}
					}
				case alert.HttpType:
					err := alertCfg.Http.Execute()
					if err != nil {
						// 执行hook
						executeHook(cfg, nil, "", endTime.Format(time.DateTime))
					}
				case alert.TcpType:
					err := alertCfg.Tcp.Execute()
					if err != nil {
						// 执行hook
						executeHook(cfg, nil, "", endTime.Format(time.DateTime))
					}
				}
			}
		}
		closer.Close()
	})
}

func resetNextTime(ctx context.Context, cfg *alertmd.Config, nextTime, runVersion int64) bool {
	b, err := alertmd.UpdateExecuteNextTime(ctx, cfg.Id, nextTime, runVersion)
	if err != nil {
		logger.Logger.Error(err)
	}
	return err == nil && b
}

func metric2Map(m *model.Sample) map[string]any {
	ret := make(map[string]any)
	for k, v := range m.Metric {
		ret[string(k)] = string(v)
	}
	ret["metric_value"] = int64(m.Value)
	return ret
}

func executeHook(cfg alertmd.Config, result map[string]any, startTime, endTime string) {
	content := cfg.GetContent()
	app, team := getAppAndTeam(cfg.AppId)
	req := event.AlertEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		EventTime: time.Now().Format(time.DateTime),
		StartTime: startTime,
		EndTime:   endTime,
		Result:    result,
	}
	switch content.HookType {
	case commonhook.NotifyType:
		notifysrv.SendNotificationByTplId(content.HookCfg.NotifyTplId, req)
	case commonhook.WebhookType:
		webhook.TriggerWebhook(content.HookCfg.HookUrl, content.HookCfg.Secret, &req)
	}
}

func getAppAndTeam(appId string) (appmd.App, teammd.Team) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.Error(err)
		return appmd.App{}, teammd.Team{}
	}
	if !b {
		return appmd.App{}, teammd.Team{}
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.Error(err)
		return appmd.App{}, teammd.Team{}
	}
	if !b {
		return appmd.App{}, teammd.Team{}
	}
	return app, team
}
