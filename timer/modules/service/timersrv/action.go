package timersrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zall/timer/modules/model/timermd"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/lease"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

var (
	taskExecutor *executor.Executor
)

func InitTask() {
	poolSize := static.GetInt("timer.poolSize")
	if poolSize <= 0 {
		poolSize = 20
	}
	queueSize := static.GetInt("timer.queueSize")
	if queueSize <= 0 {
		queueSize = 1024
	}
	logger.Logger.Infof("start timer task service with poolSize: %v queueSize: %v", poolSize, queueSize)
	taskExecutor, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	leaser, _ := lease.NewDbLease(
		"timer-lock",
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
					time.Sleep(20 * time.Second)
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
	err := timermd.IterateExecute(ctx, time.Now().UnixMilli(), func(execute *timermd.Execute) error {
		if err := runCtx.Err(); err != nil {
			return err
		}
		return handleExecute(execute)
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func triggerTask(task *timermd.Timer, triggerBy, triggerByName string) error {
	return taskExecutor.Execute(func() {
		handleTimerTaskAndAppendLog(
			task,
			triggerBy,
			triggerByName,
			timer.ManualTriggerType,
		)
	})
}

func handleTimerTaskAndAppendLog(task *timermd.Timer, triggerBy, triggerByName string, triggerType timer.TriggerType) {
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
	timermd.InsertLog(ctx, timermd.InsertLogReqDTO{
		TimerId:     task.Id,
		TaskContent: task.GetContent(),
		ErrLog:      errLog,
		TriggerType: triggerType,
		TriggerBy:   triggerBy,
		IsSuccess:   isSuccess,
	})
	team, b, err := teammd.GetByTeamId(ctx, task.TeamId)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		return
	}
	// 失败通知
	if !isSuccess {
		notifyTimerTaskEvent(
			triggerBy,
			triggerByName,
			team,
			*task,
			triggerType.String(),
			i18n.GetByValue("timerTask.fail"),
			event.TimerTaskFailAction,
		)
	}
}

func handleExecute(execute *timermd.Execute) error {
	return taskExecutor.Execute(func() {
		ctx, closer := xormstore.Context(context.Background())
		task, b, err := timermd.GetTimerById(ctx, execute.TimerId)
		if err == nil && b && resetNextTime(ctx, &task, execute.RunVersion) {
			closer.Close()
			handleTimerTaskAndAppendLog(
				&task,
				timer.DefaultTrigger,
				timer.DefaultTrigger,
				timer.AutoTriggerType,
			)
			return
		}
		closer.Close()
	})
}

func handleTimerTask(task *timermd.Timer) error {
	t := task.GetContent()
	switch t.TaskType {
	case timer.HttpTask:
		if t.HttpTask == nil {
			return errors.New("empty http task")
		}
		return t.HttpTask.DoRequest()
	default:
		return fmt.Errorf("unsupported task type: %s", t.TaskType)
	}
}

func resetNextTime(ctx context.Context, task *timermd.Timer, runVersion int64) bool {
	cron, err := ParseCron(task.CronExp)
	if err != nil {
		logger.Logger.Errorf("parse cron: %s err: %v", task.CronExp, err)
		return false
	}
	now := time.Now()
	next := cron.Next(now)
	if next.After(now) {
		b, err := timermd.UpdateExecuteNextTime(ctx, task.Id, next.UnixMilli(), runVersion)
		if err != nil {
			logger.Logger.Error(err)
		}
		return err == nil && b
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := timermd.DisableTimer(ctx, task.Id)
		if err2 != nil {
			return err2
		}
		_, err2 = timermd.DisableExecute(ctx, task.Id)
		return err2
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return false
}
