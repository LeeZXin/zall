package zalletapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/zalletsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/zalletNode", apisession.CheckLogin)
		{
			// 创建节点
			group.POST("/create", createZalletNode)
			// 编辑节点
			group.POST("/update", updateZalletNode)
			// 删除节点
			group.DELETE("/delete/:nodeId", deleteZalletNode)
			// 节点列表
			group.GET("/list", listZalletNode)
			// 所有节点列表
			group.GET("/listAll", listAllZalletNode)
		}
	})
}

func createZalletNode(c *gin.Context) {
	var req CreateZalletNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := zalletsrv.CreateZalletNode(c, zalletsrv.CreateZalletNodeReqDTO{
			Name:       req.Name,
			AgentHost:  req.AgentHost,
			AgentToken: req.AgentToken,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateZalletNode(c *gin.Context) {
	var req UpdateZalletNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := zalletsrv.UpdateZalletNode(c, zalletsrv.UpdateZalletNodeReqDTO{
			NodeId:     req.NodeId,
			Name:       req.Name,
			AgentHost:  req.AgentHost,
			AgentToken: req.AgentToken,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteZalletNode(c *gin.Context) {
	err := zalletsrv.DeleteZalletNode(c, zalletsrv.DeleteZalletNodeReqDTO{
		Id:       cast.ToInt64(c.Param("nodeId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listZalletNode(c *gin.Context) {
	var req ListZalletNodeReqVO
	if util.ShouldBindQuery(&req, c) {
		nodes, total, err := zalletsrv.ListZalletNode(c, zalletsrv.ListZalletNodeReqDTO{
			Name:     req.Name,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(nodes, func(t zalletsrv.ZalletNodeDTO) ZalletNodeVO {
			return ZalletNodeVO{
				Id:         t.Id,
				Name:       t.Name,
				AgentHost:  t.AgentHost,
				AgentToken: t.AgentToken,
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ZalletNodeVO]{
			DataResp: ginutil.DataResp[[]ZalletNodeVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func listAllZalletNode(c *gin.Context) {
	nodes, err := zalletsrv.ListAllZalletNode(c, zalletsrv.ListAllZalletNodeReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(nodes, func(t zalletsrv.SimpleZalletNodeDTO) SimpleZalletNodeVO {
		return SimpleZalletNodeVO{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleZalletNodeVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
