package gitnodeapi

import (
	"github.com/LeeZXin/zall/git/modules/service/gitnodesrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/gitNode", apisession.CheckLogin)
		{
			group.POST("/insert", insertNode)
			group.POST("/update", updateNode)
			group.Any("/list", listNode)
			group.POST("/delete", deleteNode)
		}
	})
}

func insertNode(c *gin.Context) {
	var req InsertNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := gitnodesrv.Outer.InsertNode(c, gitnodesrv.InsertNodeReqDTO{
			Name:     req.Name,
			HttpHost: req.HttpHost,
			SshHost:  req.SshHost,
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
		err := gitnodesrv.Outer.UpdateNode(c, gitnodesrv.UpdateNodeReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			HttpHost: req.HttpHost,
			SshHost:  req.SshHost,
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
		err := gitnodesrv.Outer.DeleteNode(c, gitnodesrv.DeleteNodeReqDTO{
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
	nodes, err := gitnodesrv.Outer.ListNode(c, gitnodesrv.ListNodeReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(nodes, func(t gitnodesrv.NodeDTO) (NodeVO, error) {
		return NodeVO{
			Id:       t.Id,
			Name:     t.Name,
			HttpHost: t.HttpHost,
			SshHost:  t.SshHost,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]NodeVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
