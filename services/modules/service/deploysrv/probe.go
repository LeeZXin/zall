package deploysrv

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/detecttool"
	"github.com/LeeZXin/zall/pkg/sharding/handler"
	"github.com/LeeZXin/zall/services/modules/model/deploymd"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/http"
	"time"
)

var (
	probeEnv      string
	httpClient    *http.Client
	probeExecutor *executor.Executor
	taskHandler   *handler.ShardingPeriodicalHandler
)

func InitProbeTask() {
	probeExecutor, _ = executor.NewExecutor(20, 10, time.Minute, executor.CallerRunsStrategy)
	probeEnv = static.GetString("probe.env")
	if probeEnv == "" {
		logger.Logger.Fatal("probe task started with empty env")
	}
	logger.Logger.Infof("start probe service with env: %s", probeEnv)
	httpClient = httputil.NewRetryableHttpClient()
	taskHandler, _ = handler.NewShardingPeriodicalHandler(&handler.Config{
		HeartbeatInterval:     5 * time.Second,
		HeartbeatHandler:      doHeartbeat,
		DeleteInstanceHandler: deleteInstance,
		TaskInterval:          10 * time.Second,
		TaskHandler:           probeAction,
	})
	taskHandler.Start()
}

func probeAction() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	total, index := getInstanceIndex(ctx)
	if total > 0 && index >= 0 {
		err := deploymd.IterateService(ctx, probeEnv,
			[]deploymd.ActiveStatus{deploymd.StartedStatus, deploymd.AbnormalStatus},
			func(service *deploymd.Service) error {
				if service.Id%total == index {
					doProbe(service)
				}
				return nil
			})
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func deleteInstance() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	deploymd.DeleteProbeInstance(ctx, common.GetInstanceId(), probeEnv)
}

func doProbe(service *deploymd.Service) {
	probeExecutor.Execute(func() {
		switch service.ServiceType {
		case deploy.ProcessServiceType:
			var p deploy.ProcessConfig
			err := json.Unmarshal([]byte(service.ServiceConfig), &p)
			if err != nil {
				logger.Logger.Errorf("configId: %d, unmarshal process config failed: %v", service.ConfigId, err)
				return
			}
			if !p.IsValid() {
				logger.Logger.Errorf("configId: %d, process config is invalid", service.ConfigId)
				return
			}
			switch p.DetectConfig.DetectType {
			case deploy.TcpDetectType:
				err := detecttool.CheckTcp(p.DetectConfig.TcpDetect.Ip, p.DetectConfig.TcpDetect.Port)
				if err != nil {
					// 探针失败
					updateActiveStatus(context.Background(), service.ConfigId, deploymd.AbnormalStatus, service.ActiveStatus)
				} else {
					// 探针成功
					updateActiveStatus(context.Background(), service.ConfigId, deploymd.StartedStatus, service.ActiveStatus)
				}
			case deploy.HttpGetDetectType:
				ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancelFunc()
				err := httputil.Get(ctx, httpClient, p.DetectConfig.HttpGetDetect.Url, nil, nil)
				if err != nil {
					// 探针失败
					updateActiveStatus(ctx, service.ConfigId, deploymd.AbnormalStatus, service.ActiveStatus)
				} else {
					// 探针成功
					updateActiveStatus(ctx, service.ConfigId, deploymd.StartedStatus, service.ActiveStatus)
				}
			}
		case deploy.K8sServiceType:

		}
	})
}

func updateActiveStatus(ctx context.Context, configId int64, newStatus, oldStatus deploymd.ActiveStatus) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := deploymd.UpdateServiceActiveStatusAndProbeTimeWithOldStatus(ctx, configId, probeEnv, newStatus, oldStatus, time.Now().UnixMilli())
	if err != nil {
		logger.Logger.Error(err)
	}
}

func doHeartbeat() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	now := time.Now().UnixMilli()
	b, err := deploymd.UpdateProbeInstanceHeartbeatTime(ctx, common.GetInstanceId(), probeEnv, now)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		err = deploymd.InsertProbeInstance(ctx, common.GetInstanceId(), probeEnv)
		if err != nil {
			logger.Logger.Error(err)
			return
		}
	}
}

func getInstanceIndex(ctx context.Context) (int64, int64) {
	instances, err := deploymd.GetValidProbeInstances(ctx, probeEnv, time.Now().Add(-10*time.Second).UnixMilli())
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
