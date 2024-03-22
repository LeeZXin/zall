package gitactionapi

import (
	"github.com/LeeZXin/zall/git/modules/model/gitactionmd"
	"github.com/LeeZXin/zall/git/modules/service/gitactionsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/gitAction")
		{
			// 创建action
			group.POST("/insert", insertAction)
			// 编辑action
			group.POST("/update", updateAction)
			// 删除action
			group.POST("/delete", deleteAction)
			// 展示action列表
			group.POST("/list", listAction)
			// 手动触发action
			group.POST("/trigger", triggerAction)
		}
		group = e.Group("/api/gitActionNode", apisession.CheckLogin)
		{
			group.POST("/insert", insertNode)
			group.POST("/update", updateNode)
			group.Any("/list", listNode)
			group.POST("/delete", deleteNode)
		}
	})
}

func insertAction(c *gin.Context) {
	var req InsertActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := gitactionsrv.Outer.InsertAction(c, gitactionsrv.InsertActionReqDTO{
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

func listAction(c *gin.Context) {
	var req ListActionReqVO
	if util.ShouldBindJSON(&req, c) {
		actions, err := gitactionsrv.Outer.ListAction(c, gitactionsrv.ListActionReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListActionRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		resp.Data, _ = listutil.Map(actions, func(t gitactionmd.Action) (ActionVO, error) {
			return ActionVO{
				Id:            t.Id,
				ActionContent: t.Content,
				Created:       t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func deleteAction(c *gin.Context) {
	var req DeleteActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := gitactionsrv.Outer.DeleteAction(c, gitactionsrv.DeleteActionReqDTO{
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
		err := gitactionsrv.Outer.UpdateAction(c, gitactionsrv.UpdateActionReqDTO{
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
		err := gitactionsrv.Outer.TriggerAction(c, gitactionsrv.TriggerActionReqDTO{
			Id:       req.Id,
			Args:     req.Args,
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
		err := gitactionsrv.Outer.InsertNode(c, gitactionsrv.InsertNodeReqDTO{
			Name:     req.Name,
			HttpHost: req.HttpHost,
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
		err := gitactionsrv.Outer.UpdateNode(c, gitactionsrv.UpdateNodeReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			HttpHost: req.HttpHost,
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
		err := gitactionsrv.Outer.DeleteNode(c, gitactionsrv.DeleteNodeReqDTO{
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
	nodes, err := gitactionsrv.Outer.ListNode(c, gitactionsrv.ListNodeReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	ret := ListGitNodeRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
	}
	ret.Data, _ = listutil.Map(nodes, func(t gitactionsrv.NodeDTO) (NodeVO, error) {
		return NodeVO{
			Id:       t.Id,
			Name:     t.Name,
			HttpHost: t.HttpHost,
		}, nil
	})
	c.JSON(http.StatusOK, ret)
}
