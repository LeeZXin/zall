package detectsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/pkg/detecttool"
	"github.com/LeeZXin/zall/tcpdetect/modules/model/detectmd"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

var (
	detectExecutor *executor.Executor
	//taskHandler    *handler.ShardingPeriodicalHandler
)

func InitDetect() {
	logger.Logger.Infof("start tcp detect service")
	detectExecutor, _ = executor.NewExecutor(10, 1024, time.Minute, executor.CallerRunsStrategy)
	//taskHandler, _ = handler.NewShardingPeriodicalHandler(&handler.Config{
	//	HeartbeatInterval:     8 * time.Second,
	//	HeartbeatHandler:      doHeartbeat,
	//	DeleteInstanceHandler: deleteInstance,
	//	TaskInterval:          10 * time.Second,
	//	TaskHandler:           doExecuteTask,
	//})
	//taskHandler.Start()
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

func deleteInstance() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	// 删除数据
	detectmd.DeleteInstance(ctx, common.GetInstanceId())
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
		err := detecttool.CheckTcp(fmt.Sprintf("%s:%d", detect.Ip, detect.Port))
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
