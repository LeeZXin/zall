package alertsrv

import (
	"context"
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/pkg/alert"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/hashicorp/go-bexpr"
	"github.com/prometheus/common/model"
	"net/http"
	"time"
)

var (
	taskExecutor  *executor.Executor
	heartbeatTask *taskutil.PeriodicalTask
	executeTask   *taskutil.PeriodicalTask
	httpClient    *http.Client
)

func InitTask() {
	logger.Logger.Info("start alert task service")
	httpClient = httputil.NewRetryableHttpClient()
	taskExecutor, _ = executor.NewExecutor(20, 1024, time.Minute, executor.CallerRunsStrategy)
	// 触发心跳任务
	heartbeatTask, _ = taskutil.NewPeriodicalTask(8*time.Second, doHeartbeat)
	heartbeatTask.Start()
	quit.AddShutdownHook(func() {
		// 停止心跳任务
		heartbeatTask.Stop()
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		// 删除数据
		alertmd.DeleteInstance(ctx, common.GetInstanceId())
	}, true)
	time.Sleep(time.Second)
	// 执行任务
	executeTask, _ = taskutil.NewPeriodicalTask(5*time.Second, doExecuteTask)
	executeTask.Start()
	quit.AddShutdownHook(executeTask.Stop, true)
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
		err := alertmd.IterateConfig(ctx, time.Now().UnixMilli(), true, func(cfg *alertmd.Config) error {
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
				ev, err := bexpr.CreateEvaluator(alertCfg.Mysql.Condition)
				if err != nil {
					logger.Logger.Errorf("alert cfg: %v compile condition: %v failed with err: %v", cfg.Id, alertCfg.Mysql.Condition, err)
				} else {
					evr, err := ev.Evaluate(result)
					if err == nil && evr {
						// 执行webhook
						executeWebhook(cfg, alertCfg, result)
					}
				}
			}
		case alert.PromType:
			result, err := alertCfg.Prom.Execute(httpClient)
			if err != nil {
				logger.Logger.Errorf("alert cfg: %v execute failed with err: %v", cfg.Id, err)
			} else if result != nil {
				// 执行webhook
				executeWebhook(cfg, alertCfg, metric2Map(result.Metric))
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

func metric2Map(m model.Metric) map[string]string {
	ret := make(map[string]string)
	for k, v := range m {
		ret[string(k)] = string(v)
	}
	return ret
}

func executeWebhook(cfg *alertmd.Config, alertCfg *alert.Alert, result map[string]string) {
	// todo 告警聚合

	// 执行webhook
	_, err := alertCfg.Api.DoRequest(httpClient, nil, result)
	if err != nil {
		logger.Logger.Errorf("alert cfg: %v webhook: %v failed with err: %v", cfg.Id, alertCfg.Api.Url, err)
	}
}
