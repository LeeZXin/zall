package tasksrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

var (
	taskExecutor   *executor.Executor
	heartbeatTask  *taskutil.PeriodicalTask
	executeTask    *taskutil.PeriodicalTask
	compensateTask *taskutil.PeriodicalTask

	taskEnv string
)

func InitTask() {
	taskEnv = static.GetString("timer.env")
	if taskEnv == "" {
		logger.Logger.Fatal("timer task started with empty env")
	}
	logger.Logger.Infof("start timer task service with env: %s", taskEnv)
	taskExecutor, _ = executor.NewExecutor(20, 1024, time.Minute, executor.CallerRunsStrategy)
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
		taskmd.DeleteInstance(ctx, common.GetInstanceId(), taskEnv)
	}, true)
	time.Sleep(time.Second)
	// 执行任务
	executeTask, _ = taskutil.NewPeriodicalTask(10*time.Second, doExecuteTask)
	executeTask.Start()
	quit.AddShutdownHook(executeTask.Stop, true)
	// 异常检查任务
	compensateTask, _ = taskutil.NewPeriodicalTask(5*time.Minute, doCompensateTask)
	compensateTask.Start()
	quit.AddShutdownHook(compensateTask.Stop, true)
}

func getInstanceIndex(ctx context.Context) (int64, int64) {
	instances, err := taskmd.GetValidInstances(ctx, time.Now().Add(-10*time.Second).UnixMilli(), taskEnv)
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

func doHeartbeat() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := taskmd.UpdateHeartbeatTime(ctx, common.GetInstanceId(), time.Now().UnixMilli(), taskEnv)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		err = taskmd.InsertInstance(ctx, common.GetInstanceId(), taskEnv)
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
		err := taskmd.IterateTask(ctx, time.Now().UnixMilli(), taskmd.Pending, taskEnv, func(task *taskmd.Task) error {
			if task.Id%total == index {
				handleTask(task)
			}
			return nil
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func doCompensateTask() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	total, index := getInstanceIndex(ctx)
	if index > 0 {
		// 过去五分钟内未执行完成的任务重置下次执行时间
		err := taskmd.IterateTask(ctx, time.Now().Add(-5*time.Minute).UnixMilli(), taskmd.Running, taskEnv, func(task *taskmd.Task) error {
			if task.Id%total == index {
				resetNextTime(task, taskmd.Failed)
			}
			return nil
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func triggerTask(task *taskmd.Task, triggerBy, env string) {
	taskExecutor.Execute(func() {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		logContent, status := handleTaskContent(task)
		err := taskmd.InsertTaskLog(ctx, taskmd.InsertTaskLogReqDTO{
			TaskId:      task.Id,
			TaskContent: task.Content,
			LogContent:  logContent,
			TriggerType: taskmd.ManualTriggerType,
			TriggerBy:   triggerBy,
			TaskStatus:  status,
			Env:         env,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}

func handleTask(task *taskmd.Task) {
	taskExecutor.Execute(func() {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		b, err := taskmd.UpdateTaskStatus(ctx, taskmd.UpdateTaskStatusReqDTO{
			TaskId:    task.Id,
			NewStatus: taskmd.Running,
			Version:   task.Version,
			Env:       taskEnv,
		})
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		if !b {
			return
		}
		logContent, status := handleTaskContent(task)
		err = taskmd.InsertTaskLog(ctx, taskmd.InsertTaskLogReqDTO{
			TaskId:      task.Id,
			TaskContent: task.Content,
			LogContent:  logContent,
			TriggerType: taskmd.AutoTriggerType,
			TriggerBy:   taskmd.DefaultTrigger,
			TaskStatus:  status,
			Env:         taskEnv,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
		// 重新计算下次执行时间
		task.Version += 1
		resetNextTime(task, status)
	})

}

func handleTaskContent(task *taskmd.Task) (string, taskmd.TaskStatus) {
	var obj TaskObj
	err := json.Unmarshal([]byte(task.Content), &obj)
	if err != nil {
		return fmt.Sprintf("invalid task: %d, content: %v ", task.Id, task.Content), taskmd.Failed
	}
	log := util.NewSimpleLogger()
	switch obj.TaskType {
	case HttpTaskType:
		if handleHttpTask(obj.Content, log) {
			return log.ToString(), taskmd.Successful
		}
		return log.ToString(), taskmd.Failed
	default:
		return fmt.Sprintf("unsupported task type: %s", obj.TaskType), taskmd.Failed
	}
}

func resetNextTime(task *taskmd.Task, runStatus taskmd.TaskStatus) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	cron, err := parseCron(task.CronExp)
	if err != nil {
		logger.Logger.Errorf("parse cron: %s err: %v", task.CronExp, err)
		return
	}
	now := time.Now()
	next := cron.Next(now)
	if next.After(now) {
		_, err = taskmd.UpdateTaskNextTimeAndStatus(ctx, taskmd.UpdateTaskNextTimeAndStatusReqDTO{
			TaskId:   task.Id,
			Status:   taskmd.Pending,
			NextTime: next.UnixMilli(),
			Version:  task.Version,
			Env:      taskEnv,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	} else {
		_, err = taskmd.UpdateTaskStatus(ctx, taskmd.UpdateTaskStatusReqDTO{
			TaskId:    task.Id,
			NewStatus: runStatus,
			Version:   task.Version,
			Env:       taskEnv,
		})
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}
