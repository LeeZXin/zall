package actionsrv

import (
	"context"
	"github.com/LeeZXin/zall/action/modules/model/actionmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"io"
)

const (
	updateAction = iota
	accessAction
	triggerAction
)

type innerImpl struct{}

func (s *innerImpl) ExecuteAction(aid, operator string, triggerType actionmd.TriggerType) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	gitAction, b, err := actionmd.GetActionByAid(ctx, aid)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	if !b {
		logger.Logger.WithContext(ctx).Errorf("%s is not exists", aid)
		return
	}
	var p action.GraphCfg
	// 解析yaml
	err = yaml.Unmarshal([]byte(gitAction.Content), &p)
	if err != nil {
		logger.Logger.Errorf("can not marshal action yaml: %v", err)
		return
	}
	// 转换为action graph
	graph, err := p.ConvertToGraph()
	if err != nil {
		logger.Logger.Errorf("can not convert action graph: %v", err)
		return
	}
	s.execGraph(graph, gitAction.Id, gitAction.Content, operator, triggerType, gitAction.AgentUrl, gitAction.AgentToken)
}

func (s *innerImpl) execGraph(graph *action.Graph, actionId int64, actionYaml, operator string, triggerType actionmd.TriggerType, agentUrl, agentToken string) {
	// 先插入记录
	taskId, err := s.insertTaskRecord(graph, actionId, actionYaml, operator, triggerType, agentUrl, agentToken)
	if err != nil {
		return
	}
	// 执行任务
	s.runGraph(graph, taskId, agentUrl, agentToken)
}

func (s *innerImpl) runGraph(graph *action.Graph, taskId int64, agentUrl, agentToken string) {
	err := graph.Run(action.RunOpts{
		RunWithAgent: true,
		AgentUrl:     agentUrl,
		AgentToken:   agentToken,
		StepOutputFunc: func(stat action.StepOutputStat) {
			defer stat.Output.Close()
			if _, err := s.updateStepStatusWithOldStatus(
				taskId,
				stat.JobName,
				stat.Index,
				actionmd.StepWaitingStatus,
				actionmd.StepRunningStatus,
			); err != nil {
				return
			}
			// 记录日志信息
			logContent, err := io.ReadAll(stat.Output)
			if err != nil {
				logger.Logger.Error(err)
				return
			}
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			_, err = actionmd.UpdateStepLogContent(ctx, taskId, stat.JobName, stat.Index, string(logContent))
			if err != nil {
				logger.Logger.Error(err)
				return
			}
		},
		StepAfterFunc: func(err error, stat action.StepRunStat) {
			// step状态置为success/fail
			if err != nil {
				s.updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					actionmd.StepRunningStatus,
					actionmd.StepFailStatus,
				)
			} else {
				s.updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					actionmd.StepRunningStatus,
					actionmd.StepSuccessStatus,
				)
			}
		},
	})
	if err != nil {
		logger.Logger.Errorf("taskId: %v ran with err: %v", taskId, err)
		// 任务执行失败
		s.updateTaskStatusWithOldStatus(
			taskId,
			actionmd.TaskRunningStatus,
			actionmd.TaskFailStatus,
		)
	} else {
		// 任务执行成功
		s.updateTaskStatusWithOldStatus(
			taskId,
			actionmd.TaskRunningStatus,
			actionmd.TaskSuccessStatus,
		)
	}
}

func (s *innerImpl) updateTaskStatusWithOldStatus(taskId int64, oldStatus, newStatus actionmd.TaskStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := actionmd.UpdateTaskStatusWithOldStatus(
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

func (s *innerImpl) updateStepStatusWithOldStatus(taskId int64, jobName string, index int, oldStatus, newStatus actionmd.StepStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := actionmd.UpdateStepStatus(
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

func (*innerImpl) insertTaskRecord(graph *action.Graph, actionId int64, actionYaml, operator string, triggerType actionmd.TriggerType, agentUrl, agentToken string) (int64, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	var taskId int64
	infos := graph.ListJobInfo()
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 插入一条任务记录
		task, err := actionmd.InsertTask(ctx, actionmd.InsertTaskReqDTO{
			ActionId:      actionId,
			TaskStatus:    actionmd.TaskRunningStatus,
			TriggerType:   triggerType,
			ActionContent: actionYaml,
			Operator:      operator,
			AgentUrl:      agentUrl,
			AgentToken:    agentToken,
		})
		if err != nil {
			return err
		}
		taskId = task.Id
		// 插入job记录
		stepsReq := make([]actionmd.InsertStepReqDTO, 0)
		for _, job := range infos {
			for _, step := range job.Steps {
				stepsReq = append(stepsReq, actionmd.InsertStepReqDTO{
					TaskId:     taskId,
					JobName:    job.Name,
					StepName:   step.Name,
					StepIndex:  step.Index,
					StepStatus: actionmd.StepRunningStatus,
				})
			}
		}
		_, err = actionmd.BatchInsertSteps(ctx, stepsReq)
		return err
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return taskId, err
}

type outerImpl struct{}

func (*outerImpl) InsertAction(ctx context.Context, reqDTO InsertActionReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ActionSrvKeysVO.InsertAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err = checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator, updateAction)
	if err != nil {
		return
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.ActionContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidActionContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	err = actionmd.InsertAction(ctx, actionmd.InsertActionReqDTO{
		Aid:        idutil.RandomUuid(),
		TeamId:     reqDTO.TeamId,
		Name:       reqDTO.Name,
		Content:    string(yamlOut),
		AgentUrl:   reqDTO.AgentUrl,
		AgentToken: reqDTO.AgentToken,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteAction(ctx context.Context, reqDTO DeleteActionReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ActionSrvKeysVO.DeleteAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := actionmd.GetActionById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 校验权限
	err = checkPerm(ctx, repoAction.TeamId, reqDTO.Operator, updateAction)
	if err != nil {
		return
	}
	_, err = actionmd.DeleteAction(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListAction(ctx context.Context, reqDTO ListActionReqDTO) ([]actionmd.Action, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkPerm(ctx, reqDTO.TeamId, reqDTO.Operator, accessAction)
	if err != nil {
		return nil, err
	}
	ret, err := actionmd.ListAction(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return ret, nil
}

func (*outerImpl) UpdateAction(ctx context.Context, reqDTO UpdateActionReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.ActionSrvKeysVO.UpdateAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := actionmd.GetActionById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 校验权限
	err = checkPerm(ctx, repoAction.TeamId, reqDTO.Operator, updateAction)
	if err != nil {
		return
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.ActionContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidActionContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	_, err = actionmd.UpdateAction(ctx, actionmd.UpdateActionReqDTO{
		Id:         reqDTO.Id,
		Name:       reqDTO.Name,
		Content:    string(yamlOut),
		AgentUrl:   reqDTO.AgentUrl,
		AgentToken: reqDTO.AgentToken,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) TriggerAction(ctx context.Context, reqDTO TriggerActionReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	gitAction, b, err := actionmd.GetActionById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err = checkPerm(ctx, gitAction.TeamId, reqDTO.Operator, triggerAction)
	if err != nil {
		return err
	}
	go Inner.ExecuteAction(gitAction.Aid, reqDTO.Operator.Account, actionmd.ManualTriggerType)
	return nil
}

func checkPerm(ctx context.Context, teamId int64, operator apisession.UserInfo, permCode int) error {
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case updateAction:
		pass = p.PermDetail.TeamPerm.CanUpdateAction
	case accessAction:
		pass = p.PermDetail.TeamPerm.CanAccessAction
	case triggerAction:
		pass = p.PermDetail.TeamPerm.CanTriggerAction
	}
	if !pass {
		return util.UnauthorizedError()
	}
	return nil
}

func (*outerImpl) ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]TaskDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	gitAction, b, err := actionmd.GetActionById(ctx, reqDTO.ActionId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPerm(ctx, gitAction.TeamId, reqDTO.Operator, accessAction)
	if err != nil {
		return nil, 0, err
	}
	tasks, err := actionmd.GetTask(ctx, actionmd.GetTaskReqDTO{
		ActionId: reqDTO.ActionId,
		Cursor:   reqDTO.Cursor,
		Limit:    reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(tasks) == reqDTO.Limit {
		next = tasks[len(tasks)-1].Id
	}
	data, _ := listutil.Map(tasks, func(t actionmd.Task) (TaskDTO, error) {
		return TaskDTO{
			TaskStatus:    t.TaskStatus,
			TriggerType:   t.TriggerType,
			ActionContent: t.ActionContent,
			Operator:      t.Operator,
			Created:       t.Created,
		}, nil
	})
	return data, next, nil
}

func (*outerImpl) ListStep(ctx context.Context, reqDTO ListStepReqDTO) ([]StepDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, b, err := actionmd.GetTaskById(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	gitAction, b, err := actionmd.GetActionById(ctx, task.ActionId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPerm(ctx, gitAction.TeamId, reqDTO.Operator, accessAction)
	if err != nil {
		return nil, err
	}
	steps, err := actionmd.GetStepByTaskId(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(steps, func(t actionmd.Step) (StepDTO, error) {
		return StepDTO{
			JobName:    t.JobName,
			StepName:   t.StepName,
			StepIndex:  t.StepIndex,
			LogContent: t.LogContent,
			StepStatus: t.StepStatus,
			Created:    t.Created,
			Updated:    t.Updated,
		}, nil
	})
}
