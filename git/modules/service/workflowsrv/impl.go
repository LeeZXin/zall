package workflowsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/eventbus"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
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
	accessVars
	updateVars
)

type innerImpl struct{}

func (s *innerImpl) CheckWorkflowToken(ctx context.Context, repoId int64, token string) (usermd.UserInfo, bool) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wfToken, b, err := workflowmd.GetTokenByRepoIdAndContent(ctx, repoId, token)
	// 数据库错误 或 不存在 或 已过期
	if err != nil || !b || wfToken.IsExpired() {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		return usermd.UserInfo{}, false
	}
	// 获取操作用户信息
	user, b, err := usermd.GetByAccount(ctx, wfToken.Operator)
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		return usermd.UserInfo{}, false
	}
	return user.ToUserInfo(), true
}

// TaskCallback 工作流回调
func (s *innerImpl) TaskCallback(taskId string, task sshagent.TaskStatusCallbackReq) {
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
	case sshagent.SuccessStatus:
		oldStatus = workflowmd.TaskRunningStatus
		finalStatus = workflowmd.TaskSuccessStatus
	case sshagent.FailStatus:
		oldStatus = workflowmd.TaskRunningStatus
		finalStatus = workflowmd.TaskFailStatus
	case sshagent.TimeoutStatus:
		oldStatus = workflowmd.TaskRunningStatus
		finalStatus = workflowmd.TaskTimeoutStatus
	case sshagent.RunningStatus:
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
	// 如果是终态 删除token
	if finalStatus.IsEndType() {
		err = workflowmd.DeleteTokenByTaskId(ctx, taskmd.Id)
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func (s *innerImpl) FindAndExecute(reqDTO FindAndExecuteWorkflowReqDTO) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	workflowList, err := workflowmd.ListWorkflowByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	for _, wf := range workflowList {
		if !wf.Source.MatchBranchBySource(reqDTO.Source, reqDTO.Branch) {
			continue
		}
		s.Execute(wf, ExecuteWorkflowReqDTO{
			RepoPath:    reqDTO.RepoPath,
			Operator:    reqDTO.Operator,
			TriggerType: reqDTO.TriggerType,
			Branch:      reqDTO.Branch,
			PrId:        reqDTO.PrId,
		})
	}
}

func (s *innerImpl) Execute(wf workflowmd.Workflow, reqDTO ExecuteWorkflowReqDTO) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	varsList, err := workflowmd.ListVarsByRepoId(ctx, wf.RepoId)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	now := time.Now()
	bizId := now.Format("2006010215") + idutil.RandomUuid()
	var task workflowmd.Task
	gitToken := idutil.RandomUuid()
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		task, err2 = workflowmd.InsertTask(ctx, workflowmd.InsertTaskReqDTO{
			WorkflowId:  wf.Id,
			TaskStatus:  workflowmd.TaskQueueStatus,
			TriggerType: reqDTO.TriggerType,
			YamlContent: wf.YamlContent,
			Operator:    reqDTO.Operator,
			Branch:      reqDTO.Branch,
			PrId:        reqDTO.PrId,
			AgentHost:   wf.AgentHost,
			AgentToken:  wf.AgentToken,
			BizId:       bizId,
		})
		if err2 != nil {
			return err2
		}
		return workflowmd.InsertToken(ctx, workflowmd.InsertTokenReqDTO{
			RepoId:   wf.RepoId,
			TaskId:   task.Id,
			Content:  gitToken,
			Expired:  time.Now().Add(24 * time.Hour),
			Operator: reqDTO.Operator,
		})
	})
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	url := static.GetString("workflow.callback.url")
	if url == "" {
		url = fmt.Sprintf("http://%s:%d/api/v1/workflow/internal/taskCallBack", common.GetLocalIP(), common.HttpServerPort())
		logger.Logger.Infof("callback url: %s", url)
	}
	envs := make(map[string]string, len(varsList)+9)
	{
		envs["GIT_BRANCH"] = reqDTO.Branch
		envs["GIT_PR_ID"] = strconv.FormatInt(reqDTO.PrId, 10)
		envs["GIT_OPERATOR"] = reqDTO.Operator
		envs["GIT_REPO_ID"] = strconv.FormatInt(wf.RepoId, 10)
		envs["GIT_REPO_PATH"] = reqDTO.RepoPath
		envs["GIT_TOKEN"] = gitToken
		envs["GIT_TRIGGER_TYPE"] = strconv.Itoa(int(reqDTO.TriggerType))
		envs[action.EnvCallBackUrl] = url
		envs[action.EnvCallBackToken] = static.GetString("workflow.callback.token")
		for _, vars := range varsList {
			envs[vars.Name] = vars.Content
		}
	}
	err = sshagent.NewAgentCommand(wf.AgentHost, wf.AgentToken).
		ExecuteWorkflow(wf.YamlContent, bizId, envs)
	ctx2, closer2 := xormstore.Context(context.Background())
	defer closer2.Close()
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		workflowmd.DeleteTaskById(ctx2, task.Id)
		workflowmd.DeleteTokenByTaskId(ctx2, task.Id)
	} else {
		workflowmd.UpdateLastTaskIdByWorkflowId(ctx2, wf.Id, task.Id)
	}
	return err
}

type outerImpl struct{}

func newOuterService() OuterService {
	psub.Subscribe(eventbus.PullRequestEventTopic, func(data any) {
		event, ok := data.(eventbus.PullRequestEvent)
		if ok && event.Action == string(webhook.PrMergeAction) {
			Inner.FindAndExecute(FindAndExecuteWorkflowReqDTO{
				RepoId:      event.RepoId,
				RepoPath:    event.RepoPath,
				Operator:    event.Account,
				TriggerType: workflowmd.HookTriggerType,
				Branch:      event.Ref,
				Source:      workflowmd.PullRequestTriggerSource,
				PrId:        event.PrId,
			})
		}
	})
	return new(outerImpl)
}

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
	_, err = checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator, updateWorkflow)
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
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, updateWorkflow)
	if err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := workflowmd.DeleteWorkflow(ctx, reqDTO.WorkflowId)
		if err2 != nil {
			return err2
		}
		return workflowmd.DeleteTasksByWorkflowId(ctx, reqDTO.WorkflowId)
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
	_, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, err
	}
	ret, err := workflowmd.ListWorkflowByRepoId(ctx, reqDTO.RepoId)
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
	}, func(t workflowmd.Task) (*TaskWithoutYamlContentDTO, error) {
		task := task2WithoutYamlContentDto(t)
		return &task, nil
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
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, updateWorkflow)
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
	repo, err := checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, triggerWorkflow)
	if err != nil {
		return err
	}
	err = Inner.Execute(wf, ExecuteWorkflowReqDTO{
		RepoPath:    repo.Path,
		Operator:    reqDTO.Operator.Account,
		TriggerType: workflowmd.ManualTriggerType,
		Branch:      reqDTO.Branch,
		PrId:        0,
	})
	if err != nil {
		if strings.Contains(err.Error(), "out of capacity") {
			return util.NewBizErr(apicode.OutOfWorkflowCapacityErrCode, i18n.SystemTooManyOperation)
		}
		return util.NewBizErr(apicode.ProxyAbnormalErrCode, i18n.SystemProxyAbnormal)
	}
	return nil
}

func checkPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo, permCode int) (repomd.Repo, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, util.InternalError(err)
	}
	if !b {
		return repo, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return repo, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, util.InternalError(err)
	}
	if !b {
		return repo, util.UnauthorizedError()
	}
	if p.IsAdmin {
		return repo, nil
	}
	var pass bool
	switch permCode {
	case accessWorkflow:
		pass = p.PermDetail.GetRepoPerm(repoId).CanManageWorkflow ||
			p.PermDetail.GetRepoPerm(repoId).CanTriggerWorkflow
	case updateWorkflow:
		pass = p.PermDetail.GetRepoPerm(repoId).CanManageWorkflow
	case triggerWorkflow:
		pass = p.PermDetail.GetRepoPerm(repoId).CanTriggerWorkflow
	default:
		return repo, util.UnauthorizedError()
	}
	if pass {
		return repo, nil
	}
	return repo, util.UnauthorizedError()
}

func (*outerImpl) ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]TaskWithoutYamlContentDTO, int64, error) {
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
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
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
	data, _ := listutil.Map(tasks, func(t workflowmd.Task) (TaskWithoutYamlContentDTO, error) {
		return task2WithoutYamlContentDto(t), nil
	})
	return data, total, nil
}

// ListTaskByPrId 合并请求相关工作流任务列表
func (*outerImpl) ListTaskByPrId(ctx context.Context, reqDTO ListTaskByPrIdReqDTO) ([]WorkflowTaskDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	_, err = checkPermByRepoId(ctx, pr.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, err
	}
	tasks, err := workflowmd.ListTaskByPrId(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	wfIdList, _ := listutil.Map(tasks, func(t workflowmd.Task) (int64, error) {
		return t.WorkflowId, nil
	})
	nameMap, err := workflowmd.BatchGetWorkflowNameById(ctx, wfIdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tasks, func(t workflowmd.Task) (WorkflowTaskDTO, error) {
		return WorkflowTaskDTO{
			Name:                      nameMap[t.WorkflowId],
			TaskWithoutYamlContentDTO: task2WithoutYamlContentDto(t),
		}, nil
	})
}

func task2WithoutYamlContentDto(t workflowmd.Task) TaskWithoutYamlContentDTO {
	return TaskWithoutYamlContentDTO{
		TaskStatus:  t.TaskStatus,
		TriggerType: t.TriggerType,
		Branch:      t.Branch,
		Operator:    t.Operator,
		Created:     t.Created,
		Id:          t.Id,
		PrId:        t.PrId,
		Duration:    t.Duration,
		WorkflowId:  t.WorkflowId,
	}
}

func task2Dto(t workflowmd.Task) (TaskDTO, error) {
	return TaskDTO{
		TaskWithoutYamlContentDTO: task2WithoutYamlContentDto(t),
		YamlContent:               t.YamlContent,
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
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, updateWorkflow)
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
	if !b ||
		(task.TaskStatus != workflowmd.TaskRunningStatus &&
			task.TaskStatus != workflowmd.TaskQueueStatus) {
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
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, triggerWorkflow)
	if err != nil {
		return err
	}
	err = sshagent.NewAgentCommand(task.AgentHost, task.AgentToken).KillWorkflow(task.BizId)
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
	workflowmd.DeleteTokenByTaskId(ctx, task.Id)
	return nil
}

// GetTaskStatus 获取任务状态
func (*outerImpl) GetTaskStatus(ctx context.Context, reqDTO GetTaskStatusReqDTO) (sshagent.TaskStatus, error) {
	if err := reqDTO.IsValid(); err != nil {
		return sshagent.TaskStatus{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, b, err := workflowmd.GetTaskById(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return sshagent.TaskStatus{}, util.InternalError(err)
	}
	if !b {
		return sshagent.TaskStatus{}, util.InvalidArgsError()
	}
	wf, b, err := workflowmd.GetWorkflowById(ctx, task.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return sshagent.TaskStatus{}, util.InternalError(err)
	}
	if !b {
		return sshagent.TaskStatus{}, util.ThereHasBugErr()
	}
	// 校验权限
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return sshagent.TaskStatus{}, err
	}
	ret, err := sshagent.
		NewAgentCommand(task.AgentHost, task.AgentToken).
		GetWorkflowTaskStatus(task.BizId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return sshagent.TaskStatus{}, util.NewBizErr(apicode.ProxyAbnormalErrCode, i18n.SystemProxyAbnormal)
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
	case sshagent.SuccessStatus:
		return workflowmd.TaskSuccessStatus
	case sshagent.QueueStatus:
		return workflowmd.TaskQueueStatus
	case sshagent.CancelStatus:
		return workflowmd.TaskCancelStatus
	case sshagent.TimeoutStatus:
		return workflowmd.TaskTimeoutStatus
	case sshagent.FailStatus:
		return workflowmd.TaskFailStatus
	case sshagent.RunningStatus:
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
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return nil, err
	}
	logContent, err := sshagent.
		NewAgentCommand(task.AgentHost, task.AgentToken).
		GetLogContent(task.BizId, reqDTO.JobName, reqDTO.StepIndex)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return []string{}, nil
	}
	return strings.Split(logContent, "\n"), nil
}

// GetTaskDetail 获取任务详情
func (*outerImpl) GetTaskDetail(ctx context.Context, reqDTO GetTaskDetailReqDTO) (TaskDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TaskDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, b, err := workflowmd.GetTaskById(ctx, reqDTO.TaskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TaskDTO{}, util.InternalError(err)
	}
	if !b {
		return TaskDTO{}, util.InvalidArgsError()
	}
	wf, b, err := workflowmd.GetWorkflowById(ctx, task.WorkflowId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TaskDTO{}, util.InternalError(err)
	}
	if !b {
		return TaskDTO{}, util.ThereHasBugErr()
	}
	// 校验权限
	_, err = checkPermByRepoId(ctx, wf.RepoId, reqDTO.Operator, accessWorkflow)
	if err != nil {
		return TaskDTO{}, err
	}
	return task2Dto(task)
}

// ListVars 展示变量列表
func (*outerImpl) ListVars(ctx context.Context, reqDTO ListVarsReqDTO) ([]VarsWithoutContentDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator, accessVars)
	if err != nil {
		return nil, err
	}
	varsList, err := workflowmd.ListVarsByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(varsList, func(t workflowmd.Vars) (VarsWithoutContentDTO, error) {
		return VarsWithoutContentDTO{
			VarsId: t.Id,
			Name:   t.Name,
		}, nil
	})
}

// CreateVars 新增密钥
func (*outerImpl) CreateVars(ctx context.Context, reqDTO CreateVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator, updateVars)
	if err != nil {
		return err
	}
	exists, err := workflowmd.ExistsVars(ctx, workflowmd.ExistsVarsReqDTO{
		RepoId: reqDTO.RepoId,
		Name:   reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if exists {
		return util.AlreadyExistsError()
	}
	err = workflowmd.InsertVars(ctx, workflowmd.InsertVarsReqDTO{
		RepoId:  reqDTO.RepoId,
		Name:    reqDTO.Name,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateVars 编辑密钥
func (*outerImpl) UpdateVars(ctx context.Context, reqDTO UpdateVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	vars, b, err := workflowmd.GetVarsById(ctx, reqDTO.VarsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	_, err = checkPermByRepoId(ctx, vars.RepoId, reqDTO.Operator, updateVars)
	if err != nil {
		return err
	}
	_, err = workflowmd.UpdateVars(ctx, workflowmd.UpdateVarsReqDTO{
		Id:      reqDTO.VarsId,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteVars 删除变量
func (*outerImpl) DeleteVars(ctx context.Context, reqDTO DeleteVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	vars, b, err := workflowmd.GetVarsById(ctx, reqDTO.VarsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	_, err = checkPermByRepoId(ctx, vars.RepoId, reqDTO.Operator, updateVars)
	if err != nil {
		return err
	}
	_, err = workflowmd.DeleteVars(ctx, reqDTO.VarsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

/*
GetVarsContent 获取密钥内容
只有密钥编辑权限才可以获取内容
*/
func (*outerImpl) GetVarsContent(ctx context.Context, reqDTO GetVarsContentReqDTO) (VarsDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return VarsDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	vars, b, err := workflowmd.GetVarsById(ctx, reqDTO.VarsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return VarsDTO{}, util.InternalError(err)
	}
	if !b {
		return VarsDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	_, err = checkPermByRepoId(ctx, vars.RepoId, reqDTO.Operator, updateVars)
	if err != nil {
		return VarsDTO{}, err
	}
	return VarsDTO{
		VarsWithoutContentDTO: VarsWithoutContentDTO{
			VarsId: vars.Id,
			Name:   vars.Name,
		},
		Content: vars.Content,
	}, nil
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
