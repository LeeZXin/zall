package workflowsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/workflow"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"time"
)

const (
	updateWorkflow = iota
	accessWorkflow
	triggerWorkflow
)

type innerImpl struct{}

// TaskCallback 工作流回调
func (s *innerImpl) TaskCallback(taskId string, task workflow.TaskStatusCallbackReq) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	taskmd, b, err := workflowmd.GetTaskByBizId(ctx, taskId)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		return
	}
	var (
		oldStatus, finalStatus workflowmd.TaskStatus
	)
	switch task.Status {
	case workflow.SuccessStatus:
		oldStatus = workflowmd.TaskRunningStatus
		finalStatus = workflowmd.TaskSuccessStatus
	case workflow.FailStatus:
		oldStatus = workflowmd.TaskRunningStatus
		finalStatus = workflowmd.TaskFailStatus
	case workflow.TimeoutStatus:
		oldStatus = workflowmd.TaskRunningStatus
		finalStatus = workflowmd.TaskTimeoutStatus
	case workflow.RunningStatus:
		oldStatus = workflowmd.TaskQueueStatus
		finalStatus = workflowmd.TaskRunningStatus
	default:
		return
	}
	duration, _ := time.ParseDuration(strconv.FormatInt(task.Duration, 10) + "ms")
	_, err = workflowmd.UpdateTaskStatusAndDuration(ctx,
		taskmd.Id,
		oldStatus,
		finalStatus,
		duration,
	)
	if err != nil {
		logger.Logger.Error(err)
	}
}

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
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	now := time.Now()
	bizId := now.Format("2006010215") + idutil.RandomUuid()
	task, err := workflowmd.InsertTask(ctx, workflowmd.InsertTaskReqDTO{
		WorkflowId:  wf.Id,
		TaskStatus:  workflowmd.TaskQueueStatus,
		TriggerType: triggerType,
		YamlContent: wf.YamlContent,
		Operator:    operator,
		Branch:      branch,
		PrId:        prId,
		AgentHost:   wf.AgentHost,
		AgentToken:  wf.AgentToken,
		BizId:       bizId,
	})
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	url := static.GetString("workflow.callback.url")
	if url == "" {
		url = fmt.Sprintf("http://%s%d/api/v1/workflow/internal/taskCallBack", common.GetLocalIP(), common.HttpServerPort())
	}
	err = workflow.NewAgentCommand(wf.AgentHost, wf.AgentToken, "").
		ExecuteWorkflow(wf.YamlContent, bizId, map[string]string{
			action.EnvCallBackUrl:   url,
			action.EnvCallBackToken: static.GetString("workflow.callback.token"),
		})
	if err != nil {
		if strings.Contains(err.Error(), "out of capacity") {
			workflowmd.DeleteTaskById(ctx, task.Id)
		}
	} else {
		workflowmd.UpdateLastTaskIdByWorkflowId(ctx, wf.Id, task.Id)
		// 把合并请求和工作流关联起来
		if prId > 0 {
			pullrequestmd.BatchInsertTimeline(ctx, []pullrequestmd.InsertTimelineReqDTO{
				{
					PrId:    prId,
					Action:  pullrequestmd.NewWorkflowAction(task.Id),
					Account: operator,
				},
			})
		}
	}
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
	err = workflowmd.InsertWorkflow(ctx, workflowmd.InsertWorkflowReqDTO{
		RepoId:      reqDTO.RepoId,
		Name:        reqDTO.Name,
		YamlContent: reqDTO.YamlContent,
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
	_, err = workflowmd.UpdateWorkflow(ctx, workflowmd.UpdateWorkflowReqDTO{
		Id:         reqDTO.WorkflowId,
		Name:       reqDTO.Name,
		Content:    reqDTO.YamlContent,
		AgentHost:  reqDTO.AgentHost,
		AgentToken: reqDTO.AgentToken,
		Desc:       reqDTO.Desc,
		Source:     reqDTO.Source,
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
	err = Inner.Execute(&wf, reqDTO.Operator.Account, workflowmd.ManualTriggerType, reqDTO.Branch, 0)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		if strings.Contains(err.Error(), "out of capacity") {
			return util.NewBizErr(apicode.OutOfWorkflowCapacityErrCode, i18n.SystemTooManyOperation)
		}
		return util.InternalError(err)
	}
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
	var pass bool
	switch permCode {
	case accessWorkflow:
		pass = p.PermDetail.GetRepoPerm(repoId).CanAccessWorkflow
	case updateWorkflow:
		pass = p.PermDetail.GetRepoPerm(repoId).CanUpdateWorkflow
	case triggerWorkflow:
		pass = p.PermDetail.GetRepoPerm(repoId).CanTriggerWorkflow
	default:
		return util.UnauthorizedError()
	}
	if pass {
		return nil
	}
	return util.UnauthorizedError()
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
		Id:          t.Id,
		PrId:        t.PrId,
		YamlContent: t.YamlContent,
		Duration:    t.Duration,
	}, nil
}

// GetWorkflowDetail 获取工作流详情
func (*outerImpl) GetWorkflowDetail(ctx context.Context, reqDTO GetWorkflowDetailReqDTO) (WorkflowDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return WorkflowDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wf, b, err := workflowmd.GetWorkflowById(ctx, reqDTO.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return WorkflowDTO{}, util.InternalError(err)
	}
	if !b {
		return WorkflowDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, updateWorkflow)
	if err != nil {
		return WorkflowDTO{}, err
	}
	return workflow2Dto(wf), nil
}

// KillWorkflowTask 停止工作流
func (*outerImpl) KillWorkflowTask(ctx context.Context, reqDTO KillWorkflowTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, b, err := workflowmd.GetTaskById(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b && task.TaskStatus == workflowmd.TaskCancelStatus {
		return nil
	}
	if !b || (task.TaskStatus != workflowmd.TaskRunningStatus && task.TaskStatus != workflowmd.TaskQueueStatus) {
		return util.InvalidArgsError()
	}
	wf, b, err := workflowmd.GetWorkflowById(ctx, task.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, triggerWorkflow)
	if err != nil {
		return err
	}
	err = workflow.NewAgentCommand(task.AgentHost, task.AgentToken, "").KillWorkflow(task.BizId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	_, err = workflowmd.UpdateTaskStatusAndDuration(ctx,
		reqDTO.TaskId,
		task.TaskStatus,
		workflowmd.TaskCancelStatus,
		time.Since(task.Created),
	)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetTaskStatus 获取任务状态
func (*outerImpl) GetTaskStatus(ctx context.Context, reqDTO GetTaskStatusReqDTO) (workflow.TaskStatus, error) {
	if err := reqDTO.IsValid(); err != nil {
		return workflow.TaskStatus{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, b, err := workflowmd.GetTaskById(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflow.TaskStatus{}, util.InternalError(err)
	}
	if !b {
		return workflow.TaskStatus{}, util.InvalidArgsError()
	}
	wf, b, err := workflowmd.GetWorkflowById(ctx, task.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflow.TaskStatus{}, util.InternalError(err)
	}
	if !b {
		return workflow.TaskStatus{}, util.ThereHasBugErr()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return workflow.TaskStatus{}, err
	}
	ret, err := workflow.
		NewAgentCommand(task.AgentHost, task.AgentToken, "").
		GetWorkflowTaskStatus(task.BizId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflow.TaskStatus{}, util.InternalError(err)
	}
	// 检查结果和数据库是否一致
	{
		remoteStatus := mapTaskStatus(ret.Status)
		if remoteStatus > -1 && remoteStatus != task.TaskStatus {
			duration, _ := time.ParseDuration(fmt.Sprintf("%dms", ret.Duration))
			workflowmd.UpdateTaskStatusAndDuration(
				ctx, task.Id, task.TaskStatus, remoteStatus, duration)
		}
	}
	return ret, nil
}

func mapTaskStatus(status string) workflowmd.TaskStatus {
	switch status {
	case workflow.SuccessStatus:
		return workflowmd.TaskSuccessStatus
	case workflow.QueueStatus:
		return workflowmd.TaskQueueStatus
	case workflow.CancelStatus:
		return workflowmd.TaskCancelStatus
	case workflow.TimeoutStatus:
		return workflowmd.TaskTimeoutStatus
	case workflow.FailStatus:
		return workflowmd.TaskFailStatus
	case workflow.RunningStatus:
		return workflowmd.TaskRunningStatus
	default:
		return -1
	}
}

// GetLogContent 获取日志内容
func (*outerImpl) GetLogContent(ctx context.Context, reqDTO GetLogContentReqDTO) ([]string, error) {
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
		return nil, util.ThereHasBugErr()
	}
	// 校验权限
	err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, err
	}
	logContent, err := workflow.
		NewAgentCommand(task.AgentHost, task.AgentToken, "").
		GetLogContent(task.BizId, reqDTO.JobName, reqDTO.StepIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return []string{}, nil
	}
	return strings.Split(logContent, "\n"), nil
}

func workflow2Dto(wf workflowmd.Workflow) WorkflowDTO {
	ret := WorkflowDTO{
		Id:          wf.Id,
		Name:        wf.Name,
		Desc:        wf.Description,
		RepoId:      wf.RepoId,
		YamlContent: wf.YamlContent,
		AgentHost:   wf.AgentHost,
		AgentToken:  wf.AgentToken,
	}
	if wf.Source != nil {
		ret.Source = *wf.Source
	} else {
		ret.Source = workflowmd.Source{}
	}
	return ret
}
