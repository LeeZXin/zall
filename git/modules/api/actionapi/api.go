package actionapi

import (
	"github.com/LeeZXin/zall/git/modules/model/actionmd"
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
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
			group.POST("/list", listTask)
			group.POST("/steps", getTaskSteps)
		}
		group = e.Group("/api/actionNode", apisession.CheckLogin)
		{
			group.POST("/insert", insertNode)
			group.POST("/update", updateNode)
			group.Any("/list", listNode)
			group.Any("/all", allNode)
			group.POST("/delete", deleteNode)
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
			NodeId:        req.NodeId,
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
				ActionContent: t.Content,
				Created:       t.Created.Format(time.DateTime),
				NodeId:        t.NodeId,
				Aid:           t.Aid,
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
			ActionContent: req.ActionContent,
			NodeId:        req.NodeId,
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

func insertNode(c *gin.Context) {
	var req InsertNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.InsertNode(c, actionsrv.InsertNodeReqDTO{
			Name:     req.Name,
			HttpHost: req.HttpHost,
			Token:    req.Token,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateNode(c *gin.Context) {
	var req UpdateNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.UpdateNode(c, actionsrv.UpdateNodeReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			HttpHost: req.HttpHost,
			Token:    req.Token,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteNode(c *gin.Context) {
	var req DeleteNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := actionsrv.Outer.DeleteNode(c, actionsrv.DeleteNodeReqDTO{
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

func listNode(c *gin.Context) {
	nodes, err := actionsrv.Outer.ListNode(c, actionsrv.ListNodeReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(nodes, func(t actionsrv.NodeDTO) (NodeVO, error) {
		return NodeVO{
			Id:       t.Id,
			Name:     t.Name,
			HttpHost: t.HttpHost,
			Token:    t.Token,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]NodeVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func allNode(c *gin.Context) {
	nodes, err := actionsrv.Outer.AllNode(c, actionsrv.AllNodeReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(nodes, func(t actionsrv.SimpleNodeDTO) (SimpleNodeVO, error) {
		return SimpleNodeVO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleNodeVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
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
