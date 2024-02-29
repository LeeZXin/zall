package actionsrv

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/LeeZXin/zall/git/modules/model/actiontaskmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
	"time"
)

var (
	instanceId   = static.GetString("actions.instance.id")
	instanceName = static.GetString("actions.instance.name")
	poolSize     = static.GetInt("actions.poolSize")
	queueSize    = static.GetInt("actions.queueSize")
	host         = static.GetString("actions.host")

	runner *executor.Executor

	heartbeatTask *taskutil.PeriodicalTask

	instanceHeartbeatInterval = 5 * time.Second
	validHeartbeatInterval    = 10 * time.Second

	// 获取挂起任务并尝试放入协程池
	suspendTask *taskutil.PeriodicalTask
)

func InitSrv() {
	if !actiontaskmd.IsInstanceIdValid(instanceId) {
		logger.Logger.Fatal("actions.instance.id is invalid")
	}
	if host == "" {
		logger.Logger.Fatal("actions.host not set")
	}
	if instanceName == "" {
		instanceName = "default"
	}
	// 检查异常任务 找出该节点未执行成功的任务并置失败
	checkFail()
	// 协程池默认值
	if poolSize <= 0 {
		poolSize = 10
	}
	if queueSize <= 0 {
		queueSize = 1024
	}
	// 上报心跳
	heartbeatTask, _ = taskutil.NewPeriodicalTask(instanceHeartbeatInterval, doHeartbeat)
	heartbeatTask.Start()
	quit.AddShutdownHook(heartbeatTask.Stop, true)
	// 每十分钟检查挂起任务
	suspendTask, _ = taskutil.NewPeriodicalTask(10*time.Minute, checkSuspend)
	quit.AddShutdownHook(suspendTask.Stop, true)
	// 删除实例
	quit.AddShutdownHook(deleteInstance, true)
	// 协程池
	runner, _ = executor.NewExecutor(poolSize, queueSize, time.Minute, executor.AbortStrategy)
	quit.AddShutdownHook(runner.Shutdown)
}

func TriggerGitAction(ctx context.Context, reqDTO action.Webhook) {
	// 检查参数 暂时只支持分支push
	if !strings.HasPrefix(reqDTO.Ref, git.BranchPrefix) {
		return
	}
	if reqDTO.YamlContent != "" {
		runWithYamlContent(ctx, reqDTO, reqDTO.YamlContent, 0)
	} else {
		runWithoutYamlContent(ctx, reqDTO, 0)
	}
}

func runWithYamlContent(ctx context.Context, reqDTO action.Webhook, yamlContent string, taskId int64) {
	var p action.GraphCfg
	// 解析yaml
	err := yaml.Unmarshal([]byte(yamlContent), &p)
	if err != nil {
		return
	}
	// 转换为action graph
	graph, err := p.ConvertToGraph()
	if err != nil {
		return
	}
	// 获取是否配置push
	refs, b := graph.GetSupportedRefs(action.PushAction)
	if !b {
		return
	}
	// 配置里包含push 分支
	if util.WildcardMatchBranches(refs.Branches, util.BaseRefName(reqDTO.Ref)) {
		execGraph(ctx, graph, reqDTO, taskId)
	}
}

func runWithoutYamlContent(ctx context.Context, reqDTO action.Webhook, taskId int64) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	actions, err := repomd.ListAction(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	for _, n := range actions {
		runWithYamlContent(ctx, reqDTO, n.Content, taskId)
	}
}

// SelectAndIncrJobCountInstances 选择job-count最小的节点 并自增 for update加锁
func SelectAndIncrJobCountInstances(ctx context.Context, instanceId string) (actiontaskmd.Instance, bool, error) {
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	least := time.Now().Add(-validHeartbeatInterval).UnixMilli()
	var (
		instance      actiontaskmd.Instance
		isNotFoundErr bool
	)
	err := mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		var (
			err error
			b   bool
		)
		instance, b, err = actiontaskmd.SelectLeastJobCountInstancesForUpdate(ctx, least, instanceId)
		if err != nil {
			return err
		}
		if !b {
			isNotFoundErr = true
			return errors.New("not found")
		}
		_, err = actiontaskmd.IncrJobCount(ctx, instance.InstanceId)
		if err != nil {
			return err
		}
		return nil
	})
	if !isNotFoundErr && err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return actiontaskmd.Instance{}, false, util.InternalError(err)
	}
	if instance.InstanceId == "" {
		return actiontaskmd.Instance{}, false, nil
	}
	return instance, true, nil
}

func execGraph(ctx context.Context, graph *action.Graph, reqDTO action.Webhook, taskId int64) {
	if taskId == 0 {
		var (
			err    error
			closer xormutil.Closer
		)
		ctx, closer = mysqlstore.Context(ctx)
		defer closer.Close()
		// 先插入记录
		taskId, err = insertTaskRecord(ctx, graph, reqDTO)
		if err != nil {
			return
		}
	}
	// 执行任务
	if err := runner.Execute(func() {
		runGraph(graph, taskId, reqDTO)
	}); err != nil {
		// 超出协程池范围
		updateTaskStatusWithOldStatus(
			taskId,
			actiontaskmd.TaskWaitingStatus,
			actiontaskmd.TaskSuspendStatus,
		)
	}
}

func runGraph(graph *action.Graph, taskId int64, reqDTO action.Webhook) {
	err := graph.Run(action.RunOpts{
		Args: map[string]string{
			"GIT_BRANCH": util.BaseRefName(reqDTO.Ref),
		},
		BeforeStartFunc: func(stat action.GraphRunStat) error {
			// job状态置为running
			updated, err := updateTaskStatusWithOldStatus(
				taskId,
				actiontaskmd.TaskWaitingStatus,
				actiontaskmd.TaskRunningStatus,
			)
			if err != nil {
				return err
			}
			if !updated {
				return errors.New("failed")
			}
			return nil
		},
		StepOutputFunc: func(stat action.StepOutputStat) {
			defer stat.Output.Close()
			if _, err := updateStepStatusWithOldStatus(
				taskId,
				stat.JobName,
				stat.Index,
				actiontaskmd.StepWaitingStatus,
				actiontaskmd.StepRunningStatus,
			); err != nil {
				return
			}
			// 记录日志信息
			logContent, err := io.ReadAll(stat.Output)
			if err != nil {
				logger.Logger.Error(err)
				return
			}
			ctx, closer := mysqlstore.Context(context.Background())
			defer closer.Close()
			_, err = actiontaskmd.UpdateStepLogContent(ctx, taskId, stat.JobName, stat.Index, string(logContent))
			if err != nil {
				logger.Logger.Error(err)
				return
			}
		},
		StepAfterFunc: func(err error, stat action.StepRunStat) {
			// step状态置为success/fail
			if err != nil {
				updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					actiontaskmd.StepRunningStatus,
					actiontaskmd.StepFailStatus,
				)
			} else {
				updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					actiontaskmd.StepRunningStatus,
					actiontaskmd.StepSuccessStatus,
				)
			}
		},
	})
	if err != nil {
		// 任务执行失败
		updateTaskStatusWithOldStatus(
			taskId,
			actiontaskmd.TaskRunningStatus,
			actiontaskmd.TaskFailStatus,
		)
	} else {
		// 任务执行成功
		updateTaskStatusWithOldStatus(
			taskId,
			actiontaskmd.TaskRunningStatus,
			actiontaskmd.TaskSuccessStatus,
		)
	}
}

func updateTaskStatusWithOldStatus(taskId int64, oldStatus, newStatus actiontaskmd.TaskStatus) (bool, error) {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	b, err := actiontaskmd.UpdateTaskStatus(
		ctx,
		taskId,
		oldStatus,
		newStatus,
	)
	if err != nil {
		logger.Logger.Error(err)
	}
	return b, err
}

func updateStepStatusWithOldStatus(taskId int64, jobName string, index int, oldStatus, newStatus actiontaskmd.StepStatus) (bool, error) {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	b, err := actiontaskmd.UpdateStepStatus(
		ctx,
		taskId,
		jobName,
		index,
		oldStatus,
		newStatus,
	)
	if err != nil {
		logger.Logger.Error(err)
	}
	return b, err
}

func insertTaskRecord(ctx context.Context, graph *action.Graph, reqDTO action.Webhook) (int64, error) {
	var taskId int64
	infos := graph.ListJobInfo()
	err := mysqlstore.WithTx(ctx, func(ctx context.Context) error {
		// 插入一条任务记录
		task, err := actiontaskmd.InsertTask(ctx, actiontaskmd.InsertTaskReqDTO{
			RepoId:      reqDTO.RepoId,
			TaskName:    graph.Name,
			InstanceId:  instanceId,
			TriggerType: actiontaskmd.TriggerType(reqDTO.TriggerType),
			TaskType:    actiontaskmd.GitTaskType,
			TaskStatus:  actiontaskmd.TaskWaitingStatus,
			Hook:        reqDTO,
		})
		if err != nil {
			return err
		}
		taskId = task.Id
		// 插入job记录
		stepsReq := make([]actiontaskmd.InsertStepReqDTO, 0)
		for _, job := range infos {
			for _, step := range job.Steps {
				stepsReq = append(stepsReq, actiontaskmd.InsertStepReqDTO{
					TaskId:     taskId,
					JobName:    job.Name,
					StepName:   step.Name,
					StepIndex:  step.Index,
					StepStatus: actiontaskmd.StepWaitingStatus,
				})
			}
		}
		_, err = actiontaskmd.BatchInsertSteps(ctx, stepsReq)
		return err
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return taskId, err
}

// doHeartbeat 执行心跳
func doHeartbeat() {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	updated, err := actiontaskmd.UpdateInstanceHeartbeat(ctx, instanceId)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !updated {
		err = actiontaskmd.InsertInstance(ctx, instanceId, instanceName, host)
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

// deleteInstance 删除实例
func deleteInstance() {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	err := actiontaskmd.DeleteInstance(ctx, instanceId)
	if err != nil {
		logger.Logger.Error(err)
	}
}

func checkFail() {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	rows, err := actiontaskmd.UpdateInstanceTaskFailStatus(ctx, instanceId)
	if err != nil {
		logger.Logger.Error(err)
	} else if rows > 0 {
		logger.Logger.Infof("auto fail task node: %s taskCount: %d", instanceId, rows)
	}
}

func checkSuspend() {
	ctx, closer := mysqlstore.Context(context.Background())
	defer closer.Close()
	// 遍历挂起任务 并尝试加入协程池
	err := actiontaskmd.IterateTask(ctx, instanceId, actiontaskmd.TaskSuspendStatus, func(task *actiontaskmd.Task) error {
		updated, err := updateTaskStatusWithOldStatus(
			task.Id,
			actiontaskmd.TaskSuspendStatus,
			actiontaskmd.TaskWaitingStatus,
		)
		if err != nil {
			return err
		}
		if updated {
			var hook action.Webhook
			err = json.Unmarshal([]byte(task.HookContent), &hook)
			if err != nil {
				// 转化json失败 一般不会走到这 除非手动修改数据库
				_, err = updateTaskStatusWithOldStatus(
					task.Id,
					actiontaskmd.TaskWaitingStatus,
					actiontaskmd.TaskErrStatus,
				)
				return err
			}
			// 尝试执行任务
			if hook.YamlContent != "" {
				runWithYamlContent(ctx, hook, hook.YamlContent, task.Id)
			} else {
				runWithoutYamlContent(ctx, hook, task.Id)
			}
		}
		return nil
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}
