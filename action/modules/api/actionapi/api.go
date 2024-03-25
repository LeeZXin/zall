package actionapi

import (
	"github.com/LeeZXin/zall/action/modules/model/actionmd"
	"github.com/LeeZXin/zall/action/modules/service/actionsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/action")
		{
			// 创建action
			group.POST("/insert", apisession.CheckLogin, insertAction)
			// 编辑action
			group.POST("/update", apisession.CheckLogin, updateAction)
			// 删除action
			group.POST("/delete", apisession.CheckLogin, deleteAction)
			// 展示action列表
			group.POST("/list", apisession.CheckLogin, listAction)
			// 手动触发action
			group.POST("/trigger", apisession.CheckLogin, triggerAction)
			// webhook使用
			group.Any("/execute/:aid", executeAction)
		}
		group = e.Group("/api/actionTask", apisession.CheckLogin)
		{
			// 获取执行任务列表
			group.POST("/list", listTask)
			// 获取执行任务详情
			group.POST("/steps", getTaskSteps)
		}
	})
}

func executeAction(c *gin.Context) {
	operator := c.Query("operator")
	if operator == "" {
		operator = "webhook"
	}
	triggerType, _ := strconv.ParseInt(c.Query("triggerType"), 10, 64)
	go actionsrv.Inner.ExecuteAction(c.Param("aid"), operator, actionmd.TriggerType(triggerType))
	c.String(http.StatusOK, "received")
}

func insertAction(c *gin.Context) {
	var req InsertActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.InsertAction(c, actionsrv.InsertActionReqDTO{
			Name:          req.Name,
			TeamId:        req.TeamId,
			ActionContent: req.ActionContent,
			AgentUrl:      req.AgentUrl,
			AgentToken:    req.AgentToken,
			Operator:      apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listAction(c *gin.Context) {
	var req ListActionReqVO
	if util.ShouldBindJSON(&req, c) {
		actions, err := actionsrv.Outer.ListAction(c, actionsrv.ListActionReqDTO{
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(actions, func(t actionmd.Action) (ActionVO, error) {
			return ActionVO{
				Id:            t.Id,
				Aid:           t.Aid,
				AgentUrl:      t.AgentUrl,
				AgentToken:    t.AgentToken,
				ActionContent: t.Content,
				Created:       t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]ActionVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func deleteAction(c *gin.Context) {
	var req DeleteActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.DeleteAction(c, actionsrv.DeleteActionReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateAction(c *gin.Context) {
	var req UpdateActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.UpdateAction(c, actionsrv.UpdateActionReqDTO{
			Id:            req.Id,
			Name:          req.Name,
			ActionContent: req.ActionContent,
			AgentUrl:      req.AgentUrl,
			AgentToken:    req.AgentToken,
			Operator:      apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func triggerAction(c *gin.Context) {
	var req TriggerActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.TriggerAction(c, actionsrv.TriggerActionReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listTask(c *gin.Context) {
	var req ListTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		tasks, next, err := actionsrv.Outer.ListTask(c, actionsrv.ListTaskReqDTO{
			ActionId: req.ActionId,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(tasks, func(t actionsrv.TaskDTO) (TaskVO, error) {
			return TaskVO{
				TaskStatus:    t.TaskStatus.Readable(),
				TriggerType:   t.TriggerType.Readable(),
				ActionContent: t.ActionContent,
				Operator:      t.Operator,
				Created:       t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[[]TaskVO]{
			DataResp: ginutil.DataResp[[]TaskVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}

func getTaskSteps(c *gin.Context) {
	var req GetTaskStepsReqVO
	if util.ShouldBindJSON(&req, c) {
		steps, err := actionsrv.Outer.ListStep(c, actionsrv.ListStepReqDTO{
			TaskId:   req.TaskId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(steps, func(t actionsrv.StepDTO) (StepVO, error) {
			return StepVO{
				JobName:    t.JobName,
				StepName:   t.StepName,
				StepIndex:  t.StepIndex,
				LogContent: t.LogContent,
				StepStatus: t.StepStatus.Readable(),
				Created:    t.Created.Format(time.DateTime),
				Updated:    t.Updated.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]StepVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}
