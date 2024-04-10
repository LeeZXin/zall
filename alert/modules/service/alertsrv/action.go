package alertsrv

import (
	"context"
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zall/pkg/sharding/handler"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/PaesslerAG/gval"
	"github.com/prometheus/common/model"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

var (
	taskExecutor *executor.Executor
	taskHandler  *handler.ShardingPeriodicalHandler
	httpClient   *http.Client
	limiter      *Limiter
)

func InitTask() {
	logger.Logger.Info("start alert task service")
	limiter = NewLimiter()
	httpClient = httputil.NewRetryableHttpClient()
	taskExecutor, _ = executor.NewExecutor(20, 1024, time.Minute, executor.CallerRunsStrategy)
	taskHandler, _ = handler.NewShardingPeriodicalHandler(&handler.Config{
		HeartbeatInterval:     5 * time.Second,
		HeartbeatHandler:      doHeartbeat,
		DeleteInstanceHandler: deleteInstance,
		TaskInterval:          5 * time.Second,
		TaskHandler:           doExecuteTask,
	})
	taskHandler.Start()
}

func doHeartbeat() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := alertmd.UpdateHeartbeatTime(ctx, common.GetInstanceId(), time.Now().UnixMilli())
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		err = alertmd.InsertInstance(ctx, common.GetInstanceId())
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func deleteInstance() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	// 删除数据
	alertmd.DeleteInstance(ctx, common.GetInstanceId())
}

func getInstanceIndex(ctx context.Context) (int64, int64) {
	instances, err := alertmd.GetValidInstances(ctx, time.Now().Add(-10*time.Second).UnixMilli())
	if err != nil {
		logger.Logger.Error(err)
		return int64(len(instances)), -1
	}
	if len(instances) == 0 {
		logger.Logger.Error("can not find instances")
		return int64(len(instances)), -1
	}
	for i, instance := range instances {
		if instance.InstanceId == common.GetInstanceId() {
			return int64(len(instances)), int64(i)
		}
	}
	logger.Logger.Error("can not find instances")
	return int64(len(instances)), -1
}

func doExecuteTask() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	total, index := getInstanceIndex(ctx)
	if total > 0 && index >= 0 {
		// 允许有500ms的误差
		err := alertmd.IterateConfig(ctx, time.Now().Add(500*time.Millisecond).UnixMilli(), true, func(cfg *alertmd.Config) error {
			if cfg.Id%total == index {
				handleConfig(cfg)
			}
			return nil
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func handleConfig(cfg *alertmd.Config) {
	if cfg.Content == nil || !cfg.Content.IsValid() {
		return
	}
	taskExecutor.Execute(func() {
		alertCfg := cfg.Content
		switch alertCfg.Source {
		case alert.MysqlType:
			result, err := alertCfg.Mysql.Execute()
			if err != nil {
				logger.Logger.Errorf("alert cfg: %v execute failed with err: %v", cfg.Id, err)
			} else if len(result) > 0 {
				evaluate, err := gval.Evaluate(alertCfg.Mysql.Condition, result)
				if err != nil {
					logger.Logger.Errorf("alert cfg: %v compile condition: %v failed with err: %v", cfg.Id, alertCfg.Mysql.Condition, err)
				} else if cast.ToBool(evaluate) {
					// 执行webhook
					executeWebhook(cfg, alertCfg, result)
				}
			}
		case alert.PromType:
			result, err := alertCfg.Prom.Execute(httpClient)
			if err != nil {
				logger.Logger.Errorf("alert cfg: %v execute failed with err: %v", cfg.Id, err)
			} else if result != nil {
				args := metric2Map(result)
				evaluate, err := gval.Evaluate(alertCfg.Prom.Condition, args)
				if err != nil {
					logger.Logger.Errorf("alert cfg: %v compile condition: %v failed with err: %v", cfg.Id, alertCfg.Mysql.Condition, err)
				} else if cast.ToBool(evaluate) {
					// 执行webhook
					executeWebhook(cfg, alertCfg, args)
				}
			}
		}
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		_, err := alertmd.UpdateNextTime(ctx, cfg.Id, time.Now().Add(time.Duration(cfg.IntervalSec)*time.Second).UnixMilli())
		if err != nil {
			logger.Logger.Errorf("alert cfg: %v update next time failed with err: %v", cfg.Id, err)
		}
	})
}

func metric2Map(m *model.Sample) map[string]string {
	ret := make(map[string]string)
	for k, v := range m.Metric {
		ret[string(k)] = string(v)
	}
	ret["metric_value"] = cast.ToString(int64(m.Value))
	return ret
}

func executeWebhook(cfg *alertmd.Config, alertCfg *alert.Alert, result map[string]string) {
	if limiter.TryPass(cfg.Id, time.Duration(cfg.SilenceSec)*time.Second) {
		// 执行webhook
		_, err := alertCfg.Api.DoRequest(httpClient, nil, result)
		if err != nil {
			logger.Logger.Errorf("alert cfg: %v webhook: %v failed with err: %v", cfg.Id, alertCfg.Api.Url, err)
		}
	}
}
