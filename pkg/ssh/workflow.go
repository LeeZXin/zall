package ssh

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"io"
	"path/filepath"
	"strconv"
)

type executeWorkflowReq struct {
	PrId        int64
	WfId        int64
	Operator    string
	TriggerType workflowmd.TriggerType
	Branch      string
	YamlContent string
}

func (r *executeWorkflowReq) ToMap() map[string]string {
	return map[string]string{
		"git.prId":     strconv.FormatInt(r.PrId, 10),
		"git.operator": r.Operator,
		"git.branch":   r.Branch,
	}
}

func (r *executeWorkflowReq) IsValid() bool {
	return r.WfId > 0 && len(r.Operator) > 0 && r.TriggerType.IsValid() && len(r.Branch) > 0
}

func executeWorkflow(req executeWorkflowReq, gm *graphMap) error {
	ctx, closer := xormstore.Context(context.Background())
	var p action.GraphCfg
	// 解析yaml
	err := yaml.Unmarshal([]byte(req.YamlContent), &p)
	if err != nil {
		return err
	}
	graph, err := p.ConvertToGraph()
	if err != nil {
		return err
	}
	// 先插入记录
	taskId, err := insertTaskRecord(graph, req)
	if err != nil {
		return err
	}
	// 插入map
	if gm.PutIfAbsent(taskId, graph) {
		if req.PrId > 0 {
			// by case 让合并请求和工作流关联
			pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
				{
					PrId:    req.PrId,
					Action:  pullrequestmd.NewWorkflowAction(taskId),
					Account: req.Operator,
				},
			})
			closer.Close()
		} else {
			closer.Close()
		}
		// 执行任务
		runGraph(graph, taskId, req.ToMap())
		// 删除map
		gm.Remove(taskId)
	} else {
		closer.Close()
	}
	return nil
}

func runGraph(graph *action.Graph, taskId int64, args map[string]string) {
	err := graph.Run(action.RunOpts{
		Workdir: filepath.Join(pwdDir, "workflow", strconv.FormatInt(taskId, 10)),
		Args:    args,
		StepOutputFunc: func(stat action.StepOutputStat) {
			defer stat.Output.Close()
			if _, err := updateStepStatusWithOldStatus(
				taskId,
				stat.JobName,
				stat.Index,
				workflowmd.StepWaitingStatus,
				workflowmd.StepRunningStatus,
			); err != nil {
				return
			}
			// 记录日志信息
			logContent, err := io.ReadAll(stat.Output)
			if err != nil {
				return
			}
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			workflowmd.UpdateStepLogContent(ctx, taskId, stat.JobName, stat.Index, string(logContent))
		},
		StepAfterFunc: func(err error, stat action.StepRunStat) {
			// step状态置为success/fail
			if err != nil {
				updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					workflowmd.StepRunningStatus,
					workflowmd.StepFailStatus,
				)
			} else {
				updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					workflowmd.StepRunningStatus,
					workflowmd.StepSuccessStatus,
				)
			}
		},
	})
	if err != nil {
		// 任务执行失败
		updateTaskStatusWithOldStatus(
			taskId,
			workflowmd.TaskRunningStatus,
			workflowmd.TaskFailStatus,
		)
	} else {
		// 任务执行成功
		updateTaskStatusWithOldStatus(
			taskId,
			workflowmd.TaskRunningStatus,
			workflowmd.TaskSuccessStatus,
		)
	}
}

func updateTaskStatusWithOldStatus(taskId int64, oldStatus, newStatus workflowmd.TaskStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := workflowmd.UpdateTaskStatusWithOldStatus(
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

func updateStepStatusWithOldStatus(taskId int64, jobName string, index int, oldStatus, newStatus workflowmd.StepStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := workflowmd.UpdateStepStatus(
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

func insertTaskRecord(graph *action.Graph, req executeWorkflowReq) (int64, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	var taskId int64
	infos := graph.ListJobInfo()
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 插入一条任务记录
		task, err := workflowmd.InsertTask(ctx, workflowmd.InsertTaskReqDTO{
			WorkflowId:  req.WfId,
			TaskStatus:  workflowmd.TaskRunningStatus,
			TriggerType: req.TriggerType,
			Operator:    req.Operator,
			Branch:      req.Branch,
			YamlContent: req.YamlContent,
			PrId:        req.PrId,
		})
		if err != nil {
			return err
		}
		taskId = task.Id
		// 插入job记录
		stepsReq := make([]workflowmd.InsertStepReqDTO, 0)
		for _, job := range infos {
			for _, step := range job.Steps {
				stepsReq = append(stepsReq, workflowmd.InsertStepReqDTO{
					WorkflowId: req.WfId,
					TaskId:     taskId,
					JobName:    job.Name,
					StepName:   step.Name,
					StepIndex:  step.Index,
					StepStatus: workflowmd.StepRunningStatus,
				})
			}
		}
		_, err = workflowmd.BatchInsertSteps(ctx, stepsReq)
		if err != nil {
			return err
		}
		// 更新最新任务id
		_, err = workflowmd.UpdateLastTaskIdByWorkflowId(ctx, req.WfId, taskId)
		return err
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return taskId, err
}
