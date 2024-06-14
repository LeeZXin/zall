package servicesrv

import (
	"context"
	"crypto/tls"
	"github.com/LeeZXin/zall/deploy/modules/model/servicemd"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/pkg/detecttool"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/lease"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/http"
	"time"
)

var (
	taskExecutor *executor.Executor
	taskEnv      string
	httpClient   *http.Client
)

func initTask() {
	taskEnv = static.GetString("probe.env")
	if taskEnv == "" {
		logger.Logger.Fatal("probe task started with empty env")
	}
	poolSize := static.GetInt("probe.poolSize")
	if poolSize <= 0 {
		poolSize = 20
	}
	queueSize := static.GetInt("probe.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	httpClient = &http.Client{
		Transport: &httputil.RetryableRoundTripper{
			Delegated: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				TLSHandshakeTimeout: time.Second,
				MaxIdleConns:        10,
				IdleConnTimeout:     30 * time.Second,
				MaxConnsPerHost:     10,
			},
		},
		Timeout: time.Second,
	}
	logger.Logger.Infof("start probe task service with env: %s poolSize: %v queueSize: %v", taskEnv, poolSize, queueSize)
	taskExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	leaser, _ := lease.NewDbLease(
		"probe-lock-"+taskEnv,
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
					time.Sleep(10 * time.Second)
				}
			},
			Leaser: leaser,
			// 抢锁失败 空转等待时间
			WaitDuration: 30 * time.Second,
			// 锁过期时间有20秒 每8秒续命 至少2次续命成功的机会
			RenewDuration: 8 * time.Second,
			GrantCallback: func(err error, b bool) {
				if err != nil {
					logger.Logger.Errorf("probe task grant lease failed with err: %v", err)
					return
				}
				if b {
					logger.Logger.Infof("probe task grant lease success: %v", common.GetInstanceId())
				}
			},
		},
	)
	quit.AddShutdownHook(quit.ShutdownHook(stopFunc), true)
}

func doExecuteTask(runCtx context.Context) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	err := servicemd.IterateService(ctx, taskEnv, func(probe *servicemd.Service) error {
		rerr := runCtx.Err()
		if rerr == nil {
			return executeProbe(probe)
		}
		return rerr
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func executeProbe(probe *servicemd.Service) error {
	return taskExecutor.Execute(func() {
		if probe != nil {
			if run(probe.Config) {
				ctx, closer := xormstore.Context(context.Background())
				defer closer.Close()
				servicemd.UpdateProbed(ctx, probe.Id, time.Now().UnixMilli())
			}
		}
	})
}

func run(s *deploy.Service) bool {
	if s == nil {
		return false
	}
	p := s.Probe
	switch p.Type {
	case deploy.HttpProbeType:
		resp, err := httpClient.Get(p.Http.Url)
		if err == nil {
			resp.Body.Close()
		}
		return err == nil && resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest
	case deploy.TcpProbeType:
		return detecttool.CheckTcp(p.Tcp.Addr) == nil
	default:
		return false
	}
}
