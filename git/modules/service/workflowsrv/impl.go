package workflowsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"strconv"
)

const (
	updateWorkflow = iota
	accessWorkflow
	triggerWorkflow
)

type innerImpl struct{}

func (s *innerImpl) FindAndExecute(repoId int64, operator string, triggerType workflowmd.TriggerType, branch string, source workflowmd.SourceType, prId int64) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	workflowList, err := workflowmd.ListWorkflow(ctx, repoId)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	for _, wf := range workflowList {
		if !wf.Source.MatchBranchBySource(source, branch) {
			continue
		}
		s.Execute(&wf, operator, triggerType, branch, prId)
	}
}

func (s *innerImpl) Execute(wf *workflowmd.Workflow, operator string, triggerType workflowmd.TriggerType, branch string, prId int64) error {
	err := zssh.NewAgentCommand(wf.AgentHost, wf.AgentToken, "").
		ExecuteWorkflowAsync(wf.YamlContent, map[string]string{
			action.ActionWfId:        strconv.FormatInt(wf.Id, 10),
			action.ActionOperator:    operator,
			action.ActionTriggerType: strconv.Itoa(int(triggerType)),
			action.ActionBranch:      branch,
			action.ActionPrId:        strconv.FormatInt(prId, 10),
		})
	return err
}

type outerImpl struct{}

func (*outerImpl) CreateWorkflow(ctx context.Context, reqDTO CreateWorkflowReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.WorkflowSrvKeysVO.CreateWorkflow),
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
	err = checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator, updateWorkflow)
	if err != nil {
		return
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.YamlContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidWorkflowContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	err = workflowmd.InsertWorkflow(ctx, workflowmd.InsertWorkflowReqDTO{
		RepoId:      reqDTO.RepoId,
		Name:        reqDTO.Name,
		YamlContent: string(yamlOut),
		AgentHost:   reqDTO.AgentHost,
		AgentToken:  reqDTO.AgentToken,
		Source:      reqDTO.Source,
		Desc:        reqDTO.Desc,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteWorkflow(ctx context.Context, reqDTO DeleteWorkflowReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.WorkflowSrvKeysVO.DeleteWorkflow),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wf, b, err := workflowmd.GetWorkflowById(ctx, reqDTO.WorkflowId)
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
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, updateWorkflow)
	if err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := workflowmd.DeleteWorkflow(ctx, reqDTO.WorkflowId)
		if err2 != nil {
			return err2
		}
		err2 = workflowmd.DeleteTasksByWorkflowId(ctx, reqDTO.WorkflowId)
		if err2 != nil {
			return err2
		}
		return workflowmd.DeleteStepsByWorkflowId(ctx, reqDTO.WorkflowId)
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListWorkflowWithLastTask(ctx context.Context, reqDTO ListWorkflowWithLastTaskReqDTO) ([]WorkflowWithLastTaskDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, err
	}
	ret, err := workflowmd.ListWorkflow(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	taskIdList, _ := listutil.Map(ret, func(t workflowmd.Workflow) (int64, error) {
		return t.LastTaskId, nil
	})
	taskIdList, _ = listutil.Filter(taskIdList, func(i int64) (bool, error) {
		return i > 0, nil
	})
	taskList, err := workflowmd.GetTasksByIdList(ctx, taskIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	taskIdMap, _ := listutil.CollectToMap(taskList, func(t workflowmd.Task) (int64, error) {
		return t.Id, nil
	}, func(t workflowmd.Task) (*TaskDTO, error) {
		dto, _ := task2Dto(t)
		return &dto, nil
	})
	return listutil.Map(ret, func(t workflowmd.Workflow) (WorkflowWithLastTaskDTO, error) {
		return WorkflowWithLastTaskDTO{
			Id:       t.Id,
			Name:     t.Name,
			Desc:     t.Description,
			LastTask: taskIdMap[t.LastTaskId],
		}, nil
	})
}

func (*outerImpl) UpdateWorkflow(ctx context.Context, reqDTO UpdateWorkflowReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.WorkflowSrvKeysVO.UpdateWorkflow),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wf, b, err := workflowmd.GetWorkflowById(ctx, reqDTO.WorkflowId)
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
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, updateWorkflow)
	if err != nil {
		return
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.YamlContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidWorkflowContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	_, err = workflowmd.UpdateWorkflow(ctx, workflowmd.UpdateWorkflowReqDTO{
		Id:         reqDTO.WorkflowId,
		Name:       reqDTO.Name,
		Content:    string(yamlOut),
		AgentHost:  reqDTO.AgentHost,
		AgentToken: reqDTO.AgentToken,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) TriggerWorkflow(ctx context.Context, reqDTO TriggerWorkflowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wf, b, err := workflowmd.GetWorkflowById(ctx, reqDTO.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, triggerWorkflow)
	if err != nil {
		return err
	}
	Inner.Execute(&wf, reqDTO.Operator.Account, workflowmd.ManualTriggerType, reqDTO.Branch, 0)
	return nil
}

func checkPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo, permCode int) error {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	return nil
}

func (*outerImpl) ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]TaskDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wf, b, err := workflowmd.GetWorkflowById(ctx, reqDTO.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, 0, err
	}
	tasks, total, err := workflowmd.ListTaskByWorkflowId(ctx, workflowmd.ListTaskByWorkflowIdReqDTO{
		WorkflowId: reqDTO.WorkflowId,
		PageNum:    reqDTO.PageNum,
		PageSize:   reqDTO.PageSize,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(tasks, task2Dto)
	return data, total, nil
}

func task2Dto(t workflowmd.Task) (TaskDTO, error) {
	return TaskDTO{
		TaskStatus:  t.TaskStatus,
		TriggerType: t.TriggerType,
		Branch:      t.Branch,
		Operator:    t.Operator,
		Created:     t.Created,
		YamlContent: t.YamlContent,
		Id:          t.Id,
	}, nil
}

func (*outerImpl) ListStep(ctx context.Context, reqDTO ListStepReqDTO) ([]StepDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, b, err := workflowmd.GetTaskById(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	wf, b, err := workflowmd.GetWorkflowById(ctx, task.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, err
	}
	steps, err := workflowmd.GetStepByTaskId(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(steps, func(t workflowmd.Step) (StepDTO, error) {
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
