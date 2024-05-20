package workflowapi

import (
	"github.com/LeeZXin/zall/git/modules/service/workflowsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/workflow"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/workflow", apisession.CheckLogin)
		{
			// 创建工作流
			group.POST("/create", createWorkflow)
			// 编辑工作流
			group.POST("/update", updateWorkflow)
			// 删除工作流
			group.DELETE("/delete/:workflowId", deleteWorkflow)
			// 展示工作流列表
			group.GET("/list/:repoId", listWorkflow)
			// 手动触发工作流
			group.PUT("/trigger/:workflowId", triggerWorkflow)
			// 工作流详情
			group.GET("/detail/:workflowId", getWorkflowDetail)
		}
		group = e.Group("/api/workflowTask", apisession.CheckLogin)
		{
			// 停止工作流
			group.PUT("/kill/:taskId", killWorkflowTask)
			// 获取任务详情
			group.GET("/detail/:taskId", getWorkflowTask)
			// 获取执行任务列表
			group.GET("/list/:workflowId", listTask)
			// 获取执行任务详情
			group.GET("/steps/:taskId", getTaskSteps)
		}
		e.POST("/api/v1/workflow/internal/taskCallBack", internalTaskCallback)
	})
}

func internalTaskCallback(c *gin.Context) {
	if c.GetHeader("Authorization") != static.GetString("workflow.callback.token") {
		c.String(http.StatusForbidden, "invalid token")
		return
	}
	var req workflow.TaskStatus
	if ginutil.ShouldBind(&req, c) {
		workflowsrv.Inner.TaskCallback(c.Query("taskId"), req)
		c.String(http.StatusOK, "")
	}
}

func killWorkflowTask(c *gin.Context) {
	err := workflowsrv.Outer.KillWorkflowTask(c, workflowsrv.KillWorkflowTaskReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func getWorkflowDetail(c *gin.Context) {
	detail, err := workflowsrv.Outer.GetWorkflowDetail(c, workflowsrv.GetWorkflowDetailReqDTO{
		WorkflowId: cast.ToInt64(c.Param("workflowId")),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[WorkflowVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: WorkflowVO{
			Id:          detail.Id,
			Name:        detail.Name,
			Desc:        detail.Desc,
			RepoId:      detail.RepoId,
			YamlContent: detail.YamlContent,
			Source:      detail.Source,
			AgentHost:   detail.AgentHost,
			AgentToken:  detail.AgentToken,
		},
	})
}

func createWorkflow(c *gin.Context) {
	var req CreateWorkflowReqVO
	if util.ShouldBindJSON(&req, c) {
		err := workflowsrv.Outer.CreateWorkflow(c, workflowsrv.CreateWorkflowReqDTO{
			Name:        req.Name,
			RepoId:      req.RepoId,
			YamlContent: req.YamlContent,
			AgentHost:   req.AgentHost,
			AgentToken:  req.AgentToken,
			Source:      req.Source,
			Desc:        req.Desc,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listWorkflow(c *gin.Context) {
	workflows, err := workflowsrv.Outer.ListWorkflowWithLastTask(c, workflowsrv.ListWorkflowWithLastTaskReqDTO{
		RepoId:   cast.ToInt64(c.Param("repoId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(workflows, func(t workflowsrv.WorkflowWithLastTaskDTO) (WorkflowWithLastTaskVO, error) {
		ret := WorkflowWithLastTaskVO{
			Id:   t.Id,
			Name: t.Name,
			Desc: t.Desc,
		}
		if t.LastTask != nil {
			vo, _ := task2Vo(*t.LastTask)
			ret.LastTask = &vo
		}
		return ret, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]WorkflowWithLastTaskVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func deleteWorkflow(c *gin.Context) {
	err := workflowsrv.Outer.DeleteWorkflow(c, workflowsrv.DeleteWorkflowReqDTO{
		WorkflowId: cast.ToInt64(c.Param("workflowId")),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func updateWorkflow(c *gin.Context) {
	var req UpdateWorkflowReqVO
	if util.ShouldBindJSON(&req, c) {
		err := workflowsrv.Outer.UpdateWorkflow(c, workflowsrv.UpdateWorkflowReqDTO{
			WorkflowId:  req.WorkflowId,
			Name:        req.Name,
			YamlContent: req.YamlContent,
			AgentHost:   req.AgentHost,
			AgentToken:  req.AgentToken,
			Source:      req.Source,
			Desc:        req.Desc,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func triggerWorkflow(c *gin.Context) {
	err := workflowsrv.Outer.TriggerWorkflow(c, workflowsrv.TriggerWorkflowReqDTO{
		WorkflowId: cast.ToInt64(c.Param("workflowId")),
		Branch:     c.Query("branch"),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listTask(c *gin.Context) {
	var req ginutil.Page2Req
	if util.ShouldBindQuery(&req, c) {
		tasks, total, err := workflowsrv.Outer.ListTask(c, workflowsrv.ListTaskReqDTO{
			WorkflowId: cast.ToInt64(c.Param("workflowId")),
			Page2Req:   req,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(tasks, task2Vo)
		c.JSON(http.StatusOK, ginutil.Page2Resp[TaskVO]{
			DataResp: ginutil.DataResp[[]TaskVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func task2Vo(t workflowsrv.TaskDTO) (TaskVO, error) {
	return TaskVO{
		Id:          t.Id,
		TaskStatus:  t.TaskStatus,
		TriggerType: t.TriggerType,
		Operator:    t.Operator,
		YamlContent: t.YamlContent,
		Created:     t.Created.Format(time.DateTime),
		Branch:      t.Branch,
		PrId:        t.PrId,
		Duration:    t.Duration,
	}, nil
}

func getTaskSteps(c *gin.Context) {
	steps, err := workflowsrv.Outer.ListStep(c, workflowsrv.ListStepReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(steps, step2Vo)
	c.JSON(http.StatusOK, ginutil.DataResp[[]StepVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func step2Vo(t workflowsrv.StepDTO) (StepVO, error) {
	return StepVO{
		JobName:    t.JobName,
		StepName:   t.StepName,
		StepIndex:  t.StepIndex,
		LogContent: t.LogContent,
		StepStatus: t.StepStatus,
		Created:    t.Created.Format(time.DateTime),
		Duration:   t.Duration,
	}, nil
}

func getWorkflowTask(c *gin.Context) {
	detail, err := workflowsrv.Outer.GetTaskDetail(c, workflowsrv.GetTaskDetailReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	taskVo, _ := task2Vo(detail.TaskDTO)
	stepVos, _ := listutil.Map(detail.Steps, step2Vo)
	c.JSON(http.StatusOK, ginutil.DataResp[TaskWithStepsVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: TaskWithStepsVO{
			TaskVO: taskVo,
			Steps:  stepVos,
		},
	})
}
