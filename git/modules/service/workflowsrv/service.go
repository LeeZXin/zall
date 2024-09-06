package workflowsrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/model/zalletmd"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
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
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.GitWorkflowTopic, func(data any) {
			req, ok := data.(event.GitWorkflowEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.GitWorkflowEventCreateAction:
						return events.GitWorkflow.Create
					case event.GitWorkflowEventUpdateAction:
						return events.GitWorkflow.Update
					case event.GitWorkflowEventDeleteAction:
						return events.GitWorkflow.Delete
					case event.GitWorkflowEventTriggerAction:
						return events.GitWorkflow.Trigger
					case event.GitWorkflowEventKillAction:
						return events.GitWorkflow.Kill
					default:
						return false
					}
				})
			}
		})
		psub.Subscribe(event.GitWorkflowVarsTopic, func(data any) {
			req, ok := data.(event.GitWorkflowVarsEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.GitWorkflowVarsEventCreateAction:
						return events.GitWorkflowVars.Create
					case event.GitWorkflowVarsEventUpdateAction:
						return events.GitWorkflowVars.Update
					case event.GitWorkflowVarsEventDeleteAction:
						return events.GitWorkflowVars.Delete
					default:
						return false
					}
				})
			}
		})
	})
}

func CheckWorkflowToken(ctx context.Context, repoId int64, token string) (usermd.UserInfo, bool) {
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
func TaskCallback(taskId string, req sshagent.TaskStatusCallbackReq) {
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
		oldStatus sshagent.Status
	)
	switch req.Status {
	case sshagent.SuccessStatus:
		oldStatus = sshagent.RunningStatus
	case sshagent.FailStatus:
		oldStatus = sshagent.RunningStatus
	case sshagent.TimeoutStatus:
		oldStatus = sshagent.RunningStatus
	case sshagent.RunningStatus:
		oldStatus = sshagent.QueueStatus
	case sshagent.CancelStatus:
		oldStatus = sshagent.RunningStatus
	default:
		return
	}
	duration, _ := time.ParseDuration(strconv.FormatInt(req.Duration, 10) + "ms")
	if req.Task != nil {
		_, err = workflowmd.UpdateTaskStatusAndDurationAndStatusLog(ctx,
			taskmd.Id,
			oldStatus,
			req.Status,
			duration,
			*req.Task,
		)
	} else {
		_, err = workflowmd.UpdateTaskStatusAndDuration(ctx,
			taskmd.Id,
			oldStatus,
			req.Status,
			duration,
		)
	}
	if err != nil {
		logger.Logger.Error(err)
	}
	// 如果是终态 删除token
	if req.Status.IsFinalType() {
		err = workflowmd.DeleteTokenByTaskId(ctx, taskmd.Id)
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

func FindAndExecute(reqDTO FindAndExecuteWorkflowReqDTO) {
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
		Execute(wf, ExecuteWorkflowReqDTO{
			RepoPath:    reqDTO.RepoPath,
			Operator:    reqDTO.Operator,
			TriggerType: reqDTO.TriggerType,
			Branch:      reqDTO.Branch,
			PrId:        reqDTO.PrId,
		})
	}
}

func Execute(wf workflowmd.Workflow, reqDTO ExecuteWorkflowReqDTO) error {
	ctx, closer := xormstore.Context(context.Background())
	varsList, err := workflowmd.ListVarsByRepoId(ctx, wf.RepoId, nil)
	if err != nil {
		logger.Logger.Error(err)
		closer.Close()
		return err
	}
	// 获取agent
	agent, b, err := zalletmd.GetZalletNodeByNodeId(ctx, wf.AgentId)
	if err != nil {
		logger.Logger.Error(err)
		closer.Close()
		return err
	}
	if !b {
		closer.Close()
		return fmt.Errorf("zallet agent id: %d not found", wf.AgentId)
	}
	now := time.Now()
	bizId := now.Format("2006010215") + idutil.RandomUuid()
	var task workflowmd.Task
	gitToken := idutil.RandomUuid()
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		task, err2 = workflowmd.InsertTask(ctx, workflowmd.InsertTaskReqDTO{
			RepoId:       wf.RepoId,
			WorkflowId:   wf.Id,
			WorkflowName: wf.Name,
			TaskStatus:   sshagent.QueueStatus,
			TriggerType:  reqDTO.TriggerType,
			YamlContent:  wf.YamlContent,
			Operator:     reqDTO.Operator,
			Branch:       reqDTO.Branch,
			PrId:         reqDTO.PrId,
			AgentHost:    agent.AgentHost,
			AgentToken:   agent.AgentToken,
			BizId:        bizId,
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
		closer.Close()
		return err
	}
	closer.Close()
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
	err = sshagent.NewAgentCommand(agent.AgentHost, agent.AgentToken).
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

// CreateWorkflow 创建工作流
func CreateWorkflow(ctx context.Context, reqDTO CreateWorkflowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, team, err := checkManageWorkflowPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.YamlContent), &graph)
	if err != nil || graph.IsValid() != nil {
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidWorkflowContent)
	}
	// 校验agentId
	b, err := zalletmd.ExistZalletNodeByNodeId(ctx, reqDTO.AgentId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	wf, err := workflowmd.InsertWorkflow(ctx, workflowmd.InsertWorkflowReqDTO{
		RepoId:      reqDTO.RepoId,
		Name:        reqDTO.Name,
		YamlContent: reqDTO.YamlContent,
		AgentId:     reqDTO.AgentId,
		Source:      reqDTO.Source,
		Desc:        reqDTO.Desc,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyWorkflowEvent(wf, repo, team, reqDTO.Operator, event.GitWorkflowEventCreateAction)
	return nil
}

// DeleteWorkflow 删除工作流
func DeleteWorkflow(ctx context.Context, reqDTO DeleteWorkflowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	wf, repo, team, err := checkManageWorkflowPermByWorkflowId(ctx, reqDTO.WorkflowId, reqDTO.Operator)
	if err != nil {
		return err
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
		return util.InternalError(err)
	}
	notifyWorkflowEvent(wf, repo, team, reqDTO.Operator, event.GitWorkflowEventDeleteAction)
	return nil
}

// ListWorkflowWithLastTask 工作流列表 + 最近执行任务
func ListWorkflowWithLastTask(ctx context.Context, reqDTO ListWorkflowWithLastTaskReqDTO) ([]WorkflowWithLastTaskDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkAccessRepoPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	ret, err := workflowmd.ListWorkflowByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	taskIdList := listutil.MapNe(ret, func(t workflowmd.Workflow) int64 {
		return t.LastTaskId
	})
	taskIdList = listutil.FilterNe(taskIdList, func(i int64) bool {
		return i > 0
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

// UpdateWorkflow 编辑工作流
func UpdateWorkflow(ctx context.Context, reqDTO UpdateWorkflowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	wf, repo, team, err := checkManageWorkflowPermByWorkflowId(ctx, reqDTO.WorkflowId, reqDTO.Operator)
	if err != nil {
		return err
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.YamlContent), &graph)
	if err != nil || graph.IsValid() != nil {
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidWorkflowContent)
	}
	// 校验agentId
	b, err := zalletmd.ExistZalletNodeByNodeId(ctx, reqDTO.AgentId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = workflowmd.UpdateWorkflow(ctx, workflowmd.UpdateWorkflowReqDTO{
		Id:      reqDTO.WorkflowId,
		Name:    reqDTO.Name,
		Content: reqDTO.YamlContent,
		AgentId: reqDTO.AgentId,
		Desc:    reqDTO.Desc,
		Source:  reqDTO.Source,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)

	}
	notifyWorkflowEvent(wf, repo, team, reqDTO.Operator, event.GitWorkflowEventUpdateAction)
	return nil
}

// TriggerWorkflow 手动触发工作流
func TriggerWorkflow(ctx context.Context, reqDTO TriggerWorkflowReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	wf, repo, team, err := checkTriggerWorkflowPermByWorkflowId(ctx, reqDTO.WorkflowId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = Execute(wf, ExecuteWorkflowReqDTO{
		RepoPath:    repo.Path,
		Operator:    reqDTO.Operator.Account,
		TriggerType: workflowmd.ManualTriggerType,
		Branch:      reqDTO.Branch,
		PrId:        0,
	})
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "out of capacity") {
			return util.NewBizErr(apicode.OutOfWorkflowCapacityErrCode, i18n.SystemTooManyOperation)
		}
		return util.NewBizErr(apicode.ProxyAbnormalErrCode, i18n.SystemProxyAbnormal)
	}
	notifyWorkflowEvent(wf, repo, team, reqDTO.Operator, event.GitWorkflowEventTriggerAction)
	return nil
}

func checkManageWorkflowPermByWorkflowId(ctx context.Context, wfId int64, operator apisession.UserInfo) (workflowmd.Workflow, repomd.Repo, teammd.Team, error) {
	wf, b, err := workflowmd.GetWorkflowById(ctx, wfId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflowmd.Workflow{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return workflowmd.Workflow{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkManageWorkflowPermByRepoId(ctx, wf.RepoId, operator)
	return wf, repo, team, err
}

func checkManageWorkflowPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return repo, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, team, util.InternalError(err)
	}
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanManageWorkflow) {
		return repo, team, nil
	}
	return repo, team, util.UnauthorizedError()
}

func checkManageWorkflowPermByVarId(ctx context.Context, varsId int64, operator apisession.UserInfo) (workflowmd.Vars, repomd.Repo, teammd.Team, error) {
	vars, b, err := workflowmd.GetVarsById(ctx, varsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflowmd.Vars{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return workflowmd.Vars{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkManageWorkflowPermByRepoId(ctx, vars.RepoId, operator)
	return vars, repo, team, err
}

func checkAccessRepoPermByPrId(ctx context.Context, prId int64, operator apisession.UserInfo) (repomd.Repo, error) {
	pr, b, err := pullrequestmd.GetPullRequestById(ctx, prId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, util.InvalidArgsError()
	}
	return checkAccessRepoPermByRepoId(ctx, pr.RepoId, operator)
}

func checkAccessRepoPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, error) {
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
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanAccessRepo) {
		return repo, nil
	}
	return repo, util.UnauthorizedError()
}

func checkAccessRepoPermByWorkflowId(ctx context.Context, wfId int64, operator apisession.UserInfo) (workflowmd.Workflow, repomd.Repo, error) {
	wf, b, err := workflowmd.GetWorkflowById(ctx, wfId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return wf, repomd.Repo{}, util.InternalError(err)
	}
	if !b {
		return wf, repomd.Repo{}, util.InvalidArgsError()
	}
	repo, err := checkAccessRepoPermByRepoId(ctx, wf.RepoId, operator)
	return wf, repo, err
}

func checkAccessRepoPermByTaskId(ctx context.Context, taskId int64, operator apisession.UserInfo) (workflowmd.Task, workflowmd.Workflow, repomd.Repo, error) {
	task, b, err := workflowmd.GetTaskById(ctx, taskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflowmd.Task{}, workflowmd.Workflow{}, repomd.Repo{}, util.InternalError(err)
	}
	if !b {
		return workflowmd.Task{}, workflowmd.Workflow{}, repomd.Repo{}, util.InvalidArgsError()
	}
	wf, repo, err := checkAccessRepoPermByWorkflowId(ctx, task.WorkflowId, operator)
	return task, wf, repo, err
}

func checkTriggerWorkflowPermByWorkflowId(ctx context.Context, wfId int64, operator apisession.UserInfo) (workflowmd.Workflow, repomd.Repo, teammd.Team, error) {
	wf, b, err := workflowmd.GetWorkflowById(ctx, wfId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return workflowmd.Workflow{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return workflowmd.Workflow{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkTriggerWorkflowPermByRepoId(ctx, wf.RepoId, operator)
	return wf, repo, team, err
}

func checkTriggerWorkflowPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return repo, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, team, util.InternalError(err)
	}
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanTriggerWorkflow) {
		return repo, team, nil
	}
	return repo, team, util.UnauthorizedError()
}

// ListTask 工作流任务列表
func ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]TaskWithoutYamlContentDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkAccessRepoPermByWorkflowId(ctx, reqDTO.WorkflowId, reqDTO.Operator)
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
	data := listutil.MapNe(tasks, task2WithoutYamlContentDto)
	return data, total, nil
}

// ListTaskByPrId 合并请求相关工作流任务列表
func ListTaskByPrId(ctx context.Context, reqDTO ListTaskByPrIdReqDTO) ([]WorkflowTaskDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkAccessRepoPermByPrId(ctx, reqDTO.PrId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	tasks, err := workflowmd.ListTaskByPrId(ctx, reqDTO.PrId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tasks, func(t workflowmd.Task) (WorkflowTaskDTO, error) {
		return WorkflowTaskDTO{
			Name:                      t.WorkflowName,
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
func GetWorkflowDetail(ctx context.Context, reqDTO GetWorkflowDetailReqDTO) (WorkflowDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return WorkflowDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	wf, _, err := checkAccessRepoPermByWorkflowId(ctx, reqDTO.WorkflowId, reqDTO.Operator)
	if err != nil {
		return WorkflowDTO{}, err
	}
	return workflow2Dto(wf), nil
}

// KillWorkflowTask 停止工作流
func KillWorkflowTask(ctx context.Context, reqDTO KillWorkflowTaskReqDTO) error {
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
	if b && task.TaskStatus == sshagent.CancelStatus {
		return nil
	}
	if !b ||
		(task.TaskStatus != sshagent.RunningStatus &&
			task.TaskStatus != sshagent.QueueStatus) {
		return util.InvalidArgsError()
	}
	// 校验权限
	wf, repo, team, err := checkTriggerWorkflowPermByWorkflowId(ctx, task.WorkflowId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = sshagent.NewAgentCommand(task.AgentHost, task.AgentToken).KillWorkflow(task.BizId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	notifyWorkflowEvent(wf, repo, team, reqDTO.Operator, event.GitWorkflowEventKillAction)
	return nil
}

// GetTaskStatus 获取任务状态
func GetTaskStatus(ctx context.Context, reqDTO GetTaskStatusReqDTO) (sshagent.TaskStatus, error) {
	if err := reqDTO.IsValid(); err != nil {
		return sshagent.TaskStatus{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	task, _, _, err := checkAccessRepoPermByTaskId(ctx, reqDTO.TaskId, reqDTO.Operator)
	if err != nil {
		return sshagent.TaskStatus{}, err
	}
	var ret sshagent.TaskStatus
	if task.StatusLog != nil && task.StatusLog.Data.Status != "" {
		ret = task.StatusLog.Data
	} else {
		ret, err = sshagent.
			NewAgentCommand(task.AgentHost, task.AgentToken).
			GetWorkflowTaskStatus(task.BizId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return sshagent.TaskStatus{}, util.NewBizErr(apicode.ProxyAbnormalErrCode, i18n.SystemProxyAbnormal)
		}
		switch ret.Status {
		case sshagent.CancelStatus, sshagent.SuccessStatus, sshagent.FailStatus:
			duration, _ := time.ParseDuration(fmt.Sprintf("%dms", ret.Duration))
			workflowmd.UpdateTaskStatusAndDurationAndStatusLog(
				ctx,
				task.Id,
				ret.Status,
				task.TaskStatus,
				duration,
				ret,
			)
		}
	}
	return ret, nil
}

// GetLogContent 获取日志内容
func GetLogContent(ctx context.Context, reqDTO GetLogContentReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, _, _, err := checkAccessRepoPermByTaskId(ctx, reqDTO.TaskId, reqDTO.Operator)
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
func GetTaskDetail(ctx context.Context, reqDTO GetTaskDetailReqDTO) (TaskDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TaskDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, _, _, err := checkAccessRepoPermByTaskId(ctx, reqDTO.TaskId, reqDTO.Operator)
	if err != nil {
		return TaskDTO{}, err
	}
	return task2Dto(task)
}

// ListVars 展示变量列表
func ListVars(ctx context.Context, reqDTO ListVarsReqDTO) ([]VarsWithoutContentDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkManageWorkflowPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	varsList, err := workflowmd.ListVarsByRepoId(ctx, reqDTO.RepoId, []string{"id", "name"})
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
func CreateVars(ctx context.Context, reqDTO CreateVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	repo, team, err := checkManageWorkflowPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
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
	vars, err := workflowmd.InsertVars(ctx, workflowmd.InsertVarsReqDTO{
		RepoId:  reqDTO.RepoId,
		Name:    reqDTO.Name,
		Content: reqDTO.Content,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyWorkflowVarsEvent(vars, repo, team, reqDTO.Operator, event.GitWorkflowVarsEventCreateAction)
	return nil
}

// UpdateVars 编辑密钥
func UpdateVars(ctx context.Context, reqDTO UpdateVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	vars, repo, team, err := checkManageWorkflowPermByVarId(ctx, reqDTO.VarsId, reqDTO.Operator)
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
	notifyWorkflowVarsEvent(vars, repo, team, reqDTO.Operator, event.GitWorkflowVarsEventUpdateAction)
	return nil
}

// DeleteVars 删除变量
func DeleteVars(ctx context.Context, reqDTO DeleteVarsReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	vars, repo, team, err := checkManageWorkflowPermByVarId(ctx, reqDTO.VarsId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = workflowmd.DeleteVars(ctx, reqDTO.VarsId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyWorkflowVarsEvent(vars, repo, team, reqDTO.Operator, event.GitWorkflowVarsEventDeleteAction)
	return nil
}

/*
GetVarsContent 获取密钥内容
只有密钥编辑权限才可以获取内容
*/
func GetVarsContent(ctx context.Context, reqDTO GetVarsContentReqDTO) (VarsDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return VarsDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	vars, _, _, err := checkManageWorkflowPermByVarId(ctx, reqDTO.VarsId, reqDTO.Operator)
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
		AgentId:     wf.AgentId,
	}
	if wf.Source != nil {
		ret.Source = *wf.Source
	} else {
		ret.Source = workflowmd.Source{}
	}
	return ret
}

func notifyWorkflowEvent(workflow workflowmd.Workflow, repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, action event.GitWorkflowEventAction) {
	initPsub()
	psub.Publish(event.GitWorkflowTopic, event.GitWorkflowEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoId:   repo.Id,
			RepoPath: repo.Name,
			RepoName: repo.Path,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Name,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action:       action,
		WorkflowId:   workflow.Id,
		WorkflowName: workflow.Name,
	})
}

func notifyWorkflowVarsEvent(vars workflowmd.Vars, repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, action event.GitWorkflowVarsEventAction) {
	initPsub()
	psub.Publish(event.GitWorkflowVarsTopic, event.GitWorkflowVarsEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoId:   repo.Id,
			RepoPath: repo.Name,
			RepoName: repo.Path,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Name,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action:   action,
		VarsId:   vars.Id,
		VarsName: vars.Name,
	})
}
