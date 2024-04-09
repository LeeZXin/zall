package detectsrv

import (
	"context"
	"github.com/LeeZXin/zall/pkg/detecttool"
	"github.com/LeeZXin/zall/tcpdetect/modules/model/detectmd"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/http/httptask"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/url"
	"time"
)

var (
	detectExecutor *executor.Executor
	heartbeatTask  *taskutil.PeriodicalTask
	executeTask    *taskutil.PeriodicalTask
)

func InitDetect() {
	logger.Logger.Infof("start tcp detect service")
	detectExecutor, _ = executor.NewExecutor(20, 1024, time.Minute, executor.CallerRunsStrategy)
	// 触发心跳任务
	go doHeartbeat()
	heartbeatTask, _ = taskutil.NewPeriodicalTask(8*time.Second, doHeartbeat)
	heartbeatTask.Start()
	quit.AddShutdownHook(func() {
		// 停止心跳任务
		heartbeatTask.Stop()
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		// 删除数据
		detectmd.DeleteInstance(ctx, common.GetInstanceId())
	}, true)
	time.Sleep(time.Second)
	// 执行任务
	executeTask, _ = taskutil.NewPeriodicalTask(10*time.Second, doExecuteTask)
	executeTask.Start()
	quit.AddShutdownHook(executeTask.Stop, true)
	// 清除过期实例
	httptask.AppendHttpTask("clearTcpDetectInvalidInstances", func(_ []byte, _ url.Values) {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		err := detectmd.DeleteInValidInstances(ctx, time.Now().Add(-30*time.Second).UnixMilli())
		if err != nil {
			logger.Logger.Error(err)
		}
	})
	// 清除过期日志
	httptask.AppendHttpTask("clearTcpDetectExpiredLog", func(_ []byte, _ url.Values) {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		err := detectmd.DeleteLogByTime(ctx, time.Now().Add(-3*24*time.Hour))
		if err != nil {
			logger.Logger.Error(err)
		}
	})

}

func doHeartbeat() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := detectmd.UpdateHeartbeatTime(ctx, common.GetInstanceId(), time.Now().UnixMilli())
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		err = detectmd.InsertInstance(ctx, common.GetInstanceId())
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func doExecuteTask() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	total, index := getInstanceIndex(ctx)
	if total > 0 && index >= 0 {
		err := detectmd.IterateDetect(ctx, func(detect *detectmd.TcpDetect) error {
			if detect.Id%total == index {
				doDetect(detect)
			}
			return nil
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func doDetect(detect *detectmd.TcpDetect) {
	detectExecutor.Execute(func() {
		err := detecttool.CheckTcp(detect.Ip, detect.Port)
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		if err == nil {
			err = detectmd.UpdateDetectHeartbeatTime(ctx, detect.Id, time.Now().UnixMilli())
			if err != nil {
				logger.Logger.Error(err)
			}
			err = detectmd.InsertLog(ctx, detectmd.InsertLogReqDTO{
				DetectId: detect.Id,
				Ip:       detect.Ip,
				Port:     detect.Port,
				Valid:    true,
			})
		} else {
			err = detectmd.InsertLog(ctx, detectmd.InsertLogReqDTO{
				DetectId: detect.Id,
				Ip:       detect.Ip,
				Port:     detect.Port,
				Valid:    false,
			})
		}
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}

func getInstanceIndex(ctx context.Context) (int64, int64) {
	instances, err := detectmd.GetValidInstances(ctx, time.Now().Add(-10*time.Second).UnixMilli())
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
