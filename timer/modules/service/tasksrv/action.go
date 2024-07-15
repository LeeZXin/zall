package tasksrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
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
	httpClient = httputil.NewRetryableHttpClient()
	taskEnv = static.GetString("timer.env")
	if taskEnv == "" {
		logger.Logger.Fatal("timer task started with empty env")
	}
	poolSize := static.GetInt("timer.poolSize")
	if poolSize <= 0 {
		poolSize = 10
	}
	queueSize := static.GetInt("timer.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	logger.Logger.Infof("start timer task service with env: %s poolSize: %v queueSize: %v", taskEnv, poolSize, queueSize)
	taskExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	leaser, _ := lease.NewDbLease(
		"timer-lock-"+taskEnv,
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
					logger.Logger.Errorf("timer task grant lease failed with err: %v", err)
					return
				}
				if b {
					logger.Logger.Infof("timer task grant lease success: %v", common.GetInstanceId())
				}
			},
		},
	)
	quit.AddShutdownHook(quit.ShutdownHook(stopFunc), true)
}

func doExecuteTask(runCtx context.Context) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	err := taskmd.IterateExecute(ctx, time.Now().UnixMilli(), taskEnv, func(execute *taskmd.Execute) error {
		rerr := runCtx.Err()
		if rerr == nil {
			return handleExecute(execute)
		}
		return rerr
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func triggerTask(task *taskmd.Task, triggerBy string) {
	taskExecutor.Execute(func() {
		handleTimerTaskAndAppendLog(
			task,
			triggerBy,
			taskmd.ManualTriggerType,
		)
	})
}

func handleTimerTaskAndAppendLog(task *taskmd.Task, triggerBy string, triggerType taskmd.TriggerType) {
	var (
		errLog    string
		isSuccess bool
	)
	err := handleTimerTask(task)
	if err != nil {
		errLog = err.Error()
	} else {
		isSuccess = true
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	taskmd.InsertTaskLog(ctx, taskmd.InsertTaskLogReqDTO{
		TaskId:      task.Id,
		TaskContent: task.Content,
		ErrLog:      errLog,
		TriggerType: triggerType,
		TriggerBy:   triggerBy,
		IsSuccess:   isSuccess,
	})
}

func handleExecute(execute *taskmd.Execute) error {
	return taskExecutor.Execute(func() {
		ctx, closer := xormstore.Context(context.Background())
		task, b, err := taskmd.GetTaskById(ctx, execute.TaskId)
		closer.Close()
		if err == nil && b && resetNextTime(&task, execute.RunVersion) {
			handleTimerTaskAndAppendLog(
				&task,
				taskmd.DefaultTrigger,
				taskmd.AutoTriggerType,
			)
		}
	})
}

func handleTimerTask(task *taskmd.Task) error {
	if task.Content == nil {
		return errors.New("empty task content")
	}
	t := task.Content
	switch t.TaskType {
	case timer.HttpTask:
		if t.HttpTask == nil {
			return errors.New("empty http task")
		}
		return t.HttpTask.DoRequest(httpClient)
	default:
		return fmt.Errorf("unsupported task type: %s", t.TaskType)
	}
}

func resetNextTime(task *taskmd.Task, runVersion int64) bool {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	cron, err := ParseCron(task.CronExp)
	if err != nil {
		logger.Logger.Errorf("parse cron: %s err: %v", task.CronExp, err)
		return false
	}
	now := time.Now()
	next := cron.Next(now)
	if next.After(now) {
		b, err := taskmd.UpdateExecuteNextTime(ctx, task.Id, next.UnixMilli(), runVersion)
		if err != nil {
			logger.Logger.Error(err)
		}
		return err == nil && b
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := taskmd.DisableTask(ctx, task.Id)
		if err2 != nil {
			return err2
		}
		_, err2 = taskmd.DisableExecute(ctx, task.Id)
		return err2
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return false
}
