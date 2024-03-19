package gitactionsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/gitactionmd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"gopkg.in/yaml.v3"
	"io"
)

const (
	updateAction = iota
	accessAction
	triggerAction
)

type innerImpl struct{}

func (s *innerImpl) ExecuteAction(ctx context.Context, hook action.Hook) {
	var p action.GraphCfg
	// 解析yaml
	err := yaml.Unmarshal([]byte(hook.ActionYaml), &p)
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
	s.execGraph(ctx, graph, &hook)
}

func (s *innerImpl) execGraph(ctx context.Context, graph *action.Graph, hook *action.Hook) {
	var (
		err    error
		closer xormutil.Closer
	)
	ctx, closer = xormstore.Context(ctx)
	defer closer.Close()
	// 先插入记录
	taskId, err := s.insertTaskRecord(ctx, graph, hook)
	if err != nil {
		return
	}
	// 执行任务
	s.runGraph(graph, taskId, hook)
}

func (s *innerImpl) runGraph(graph *action.Graph, taskId int64, hook *action.Hook) {
	err := graph.Run(action.RunOpts{
		Args: hook.Args,
		StepOutputFunc: func(stat action.StepOutputStat) {
			defer stat.Output.Close()
			if _, err := s.updateStepStatusWithOldStatus(
				taskId,
				stat.JobName,
				stat.Index,
				gitactionmd.StepWaitingStatus,
				gitactionmd.StepRunningStatus,
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
			_, err = gitactionmd.UpdateStepLogContent(ctx, taskId, stat.JobName, stat.Index, string(logContent))
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
					gitactionmd.StepRunningStatus,
					gitactionmd.StepFailStatus,
				)
			} else {
				s.updateStepStatusWithOldStatus(
					taskId,
					stat.JobName,
					stat.Index,
					gitactionmd.StepRunningStatus,
					gitactionmd.StepSuccessStatus,
				)
			}
		},
	})
	if err != nil {
		// 任务执行失败
		s.updateTaskStatusWithOldStatus(
			taskId,
			gitactionmd.TaskRunningStatus,
			gitactionmd.TaskFailStatus,
		)
	} else {
		// 任务执行成功
		s.updateTaskStatusWithOldStatus(
			taskId,
			gitactionmd.TaskRunningStatus,
			gitactionmd.TaskSuccessStatus,
		)
	}
}

func (s *innerImpl) updateTaskStatusWithOldStatus(taskId int64, oldStatus, newStatus gitactionmd.TaskStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := gitactionmd.UpdateTaskStatus(
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

func (s *innerImpl) updateStepStatusWithOldStatus(taskId int64, jobName string, index int, oldStatus, newStatus gitactionmd.StepStatus) (bool, error) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	b, err := gitactionmd.UpdateStepStatus(
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

func (*innerImpl) insertTaskRecord(ctx context.Context, graph *action.Graph, hook *action.Hook) (int64, error) {
	var taskId int64
	infos := graph.ListJobInfo()
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 插入一条任务记录
		task, err := gitactionmd.InsertTask(ctx, gitactionmd.InsertTaskReqDTO{
			ActionId:      hook.ActionId,
			TaskName:      graph.Name,
			TriggerType:   hook.TriggerType,
			TaskStatus:    gitactionmd.TaskRunningStatus,
			ActionContent: hook.ActionYaml,
		})
		if err != nil {
			return err
		}
		taskId = task.Id
		// 插入job记录
		stepsReq := make([]gitactionmd.InsertStepReqDTO, 0)
		for _, job := range infos {
			for _, step := range job.Steps {
				stepsReq = append(stepsReq, gitactionmd.InsertStepReqDTO{
					TaskId:     taskId,
					JobName:    job.Name,
					StepName:   step.Name,
					StepIndex:  step.Index,
					StepStatus: gitactionmd.StepRunningStatus,
				})
			}
		}
		_, err = gitactionmd.BatchInsertSteps(ctx, stepsReq)
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
			OpDesc:     i18n.GetByKey(i18n.GitActionSrvKeysVO.InsertAction),
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
	err = getPerm(ctx, reqDTO.Id, reqDTO.Operator, updateAction)
	if err != nil {
		return
	}
	// 检查nodeId
	_, b, err := gitactionmd.GetNodeById(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		return util.InvalidArgsError()
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.ActionContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidActionContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	err = gitactionmd.InsertAction(ctx, gitactionmd.InsertActionReqDTO{
		RepoId:     reqDTO.Id,
		NodeId:     reqDTO.NodeId,
		Content:    string(yamlOut),
		PushBranch: reqDTO.PushBranch,
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
			OpDesc:     i18n.GetByKey(i18n.GitActionSrvKeysVO.DeleteAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := gitactionmd.GetActionById(ctx, reqDTO.Id)
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
	err = getPerm(ctx, repoAction.Id, reqDTO.Operator, updateAction)
	if err != nil {
		return
	}
	_, err = gitactionmd.DeleteAction(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListAction(ctx context.Context, reqDTO ListActionReqDTO) ([]gitactionmd.Action, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := getPerm(ctx, reqDTO.Id, reqDTO.Operator, accessAction)
	if err != nil {
		return nil, err
	}
	ret, err := gitactionmd.ListAction(ctx, reqDTO.Id)
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
			OpDesc:     i18n.GetByKey(i18n.GitActionSrvKeysVO.UpdateAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := gitactionmd.GetActionById(ctx, reqDTO.Id)
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
	err = getPerm(ctx, repoAction.Id, reqDTO.Operator, updateAction)
	if err != nil {
		return
	}
	// 检查nodeId
	_, b, err = gitactionmd.GetNodeById(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		return util.InvalidArgsError()
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.ActionContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidActionContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	_, err = gitactionmd.UpdateAction(ctx, gitactionmd.UpdateActionReqDTO{
		Id:         reqDTO.Id,
		Content:    string(yamlOut),
		NodeId:     reqDTO.NodeId,
		PushBranch: reqDTO.PushBranch,
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
	repoAction, b, err := gitactionmd.GetActionById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err = getPerm(ctx, repoAction.Id, reqDTO.Operator, triggerAction)
	if err != nil {
		return err
	}
	node, b, err := gitactionmd.GetNodeById(ctx, repoAction.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	action.ManualTriggerAction(repoAction.Content, node.HttpHost, reqDTO.Args, reqDTO.Id)
	return nil
}

func getPerm(ctx context.Context, repoId int64, operator apisession.UserInfo, permCode int) error {
	repo, b := reposrv.Inner.GetByRepoId(ctx, repoId)
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case updateAction:
		pass = p.PermDetail.GetRepoPerm(repoId).CanUpdateAction
	case accessAction:
		pass = p.PermDetail.GetRepoPerm(repoId).CanAccessAction
	case triggerAction:
		pass = p.PermDetail.GetRepoPerm(repoId).CanTriggerAction
	}
	if !pass {
		return util.UnauthorizedError()
	}
	return nil
}

func (*outerImpl) InsertNode(ctx context.Context, reqDTO InsertNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GitActionSrvKeysVO.InsertNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	err = gitactionmd.InsertNode(ctx, gitactionmd.InsertNodeReqDTO{
		Name:     reqDTO.Name,
		HttpHost: reqDTO.HttpHost,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteNode(ctx context.Context, reqDTO DeleteNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GitActionSrvKeysVO.DeleteNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, err = gitactionmd.DeleteNode(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateNode(ctx context.Context, reqDTO UpdateNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GitActionSrvKeysVO.UpdateNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = gitactionmd.UpdateNode(ctx, gitactionmd.UpdateNodeReqDTO{
		Id:       reqDTO.Id,
		Name:     reqDTO.Name,
		HttpHost: reqDTO.HttpHost,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListNode(ctx context.Context, reqDTO ListNodeReqDTO) ([]NodeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	all, err := gitactionmd.GetAllNodes(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(all, func(t gitactionmd.Node) (NodeDTO, error) {
		return NodeDTO{
			Id:       t.Id,
			Name:     t.Name,
			HttpHost: t.HttpHost,
		}, nil
	})
}
