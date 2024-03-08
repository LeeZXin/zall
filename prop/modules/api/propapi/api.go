package propapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/prop/modules/service/propsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/timeutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/prop/etcdNode", apisession.CheckLogin)
		{
			group.POST("/insert", insertEtcdNode)
			group.GET("/list", listEtcdNode)
			group.GET("/listSimple", listSimpleEtcdNode)
			group.POST("/update", updateEtcdNode)
			group.POST("/delete", deleteEtcdNode)
		}
		group = e.Group("/api/prop/content", apisession.CheckLogin)
		{
			group.Any("/grantAuth", grantAuth)
			group.POST("/insert", insertContent)
			group.POST("/list", listContent)
			group.POST("/update", updateContent)
			group.POST("/delete", deleteContent)
			group.POST("/deploy", deployContent)
			group.POST("/listHistory", listHistory)
			group.POST("/listDeploy", listDeploy)
		}
	})
}

func insertEtcdNode(c *gin.Context) {
	var req InsertEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.InsertEtcdNode(c, propsrv.InsertEtcdNodeReqDTO{
			NodeId:    req.NodeId,
			Endpoints: req.Endpoints,
			Username:  req.Username,
			Password:  req.Password,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateEtcdNode(c *gin.Context) {
	var req UpdateEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.UpdateEtcdNode(c, propsrv.UpdateEtcdNodeReqDTO{
			NodeId:    req.NodeId,
			Endpoints: req.Endpoints,
			Username:  req.Username,
			Password:  req.Password,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteEtcdNode(c *gin.Context) {
	var req DeleteEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.DeleteEtcdNode(c, propsrv.DeleteEtcdNodeReqDTO{
			NodeId:   req.NodeId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listEtcdNode(c *gin.Context) {
	nodes, err := propsrv.Outer.ListEtcdNode(c, propsrv.ListEtcdNodeReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	resp := ListEtcdNodeRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
	}
	resp.Data, _ = listutil.Map(nodes, func(t propsrv.EtcdNodeDTO) (EtcdNodeVO, error) {
		return EtcdNodeVO{
			NodeId:    t.NodeId,
			Endpoints: t.Endpoints,
			Username:  t.Username,
			Password:  t.Password,
		}, nil
	})
	c.JSON(http.StatusOK, resp)
}

func listSimpleEtcdNode(c *gin.Context) {
	nodes, err := propsrv.Outer.ListSimpleEtcdNode(c)
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ListSimpleEtcdNodeRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     nodes,
	})
}

func insertContent(c *gin.Context) {
	var req InsertContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.InsertPropContent(c, propsrv.InsertPropContentReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Content:  req.Content,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateContent(c *gin.Context) {
	var req UpdateContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.UpdatePropContent(c, propsrv.UpdatePropContentReqDTO{
			Id:       req.Id,
			Content:  req.Content,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteContent(c *gin.Context) {
	var req DeleteContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.DeletePropContent(c, propsrv.DeletePropContentReqDTO{
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

func deployContent(c *gin.Context) {
	var req DeployContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.DeployPropContent(c, propsrv.DeployPropContentReqDTO{
			Id:           req.Id,
			Version:      req.Version,
			EtcdNodeList: req.EtcdNodeList,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listContent(c *gin.Context) {
	var req ListContentReqVO
	if util.ShouldBindJSON(&req, c) {
		contents, err := propsrv.Outer.ListPropContent(c, propsrv.ListPropContentReqDTO{
			AppId:    req.AppId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListContentRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		resp.Data, _ = listutil.Map(contents, func(t propsrv.PropContentDTO) (PropContentVO, error) {
			return PropContentVO{
				Id:    t.Id,
				AppId: t.AppId,
				Name:  t.Name,
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func listHistory(c *gin.Context) {
	var req ListHistoryReqVO
	if util.ShouldBindJSON(&req, c) {
		histories, cursor, err := propsrv.Outer.ListHistory(c, propsrv.ListHistoryReqDTO{
			ContentId: req.ContentId,
			Version:   req.Version,
			Cursor:    req.Cursor,
			Limit:     req.Limit,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListHistoryRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   cursor,
		}
		resp.Data, _ = listutil.Map(histories, func(t propsrv.HistoryDTO) (HistoryVO, error) {
			return HistoryVO{
				ContentId: t.ContentId,
				Content:   t.Content,
				Version:   t.Version,
				Created:   t.Created.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func grantAuth(c *gin.Context) {
	var req GrantAuthReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.GrantAuth(c, propsrv.GrantAuthReqDTO{
			AppId:    req.AppId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listDeploy(c *gin.Context) {
	var req ListDeployReqVO
	if util.ShouldBindJSON(&req, c) {
		deploys, cursor, err := propsrv.Outer.ListDeploy(c, propsrv.ListDeployReqDTO{
			ContentId: req.ContentId,
			Version:   req.Version,
			Cursor:    req.Cursor,
			Limit:     req.Limit,
			NodeId:    req.NodeId,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListDeployRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   cursor,
		}
		resp.Data, _ = listutil.Map(deploys, func(t propsrv.DeployDTO) (DeployVO, error) {
			return DeployVO{
				ContentId: t.ContentId,
				Content:   t.Content,
				Version:   t.Version,
				NodeId:    t.NodeId,
				Created:   t.Created.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}
