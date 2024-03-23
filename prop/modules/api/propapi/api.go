package propapi

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/prop/modules/service/propsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/http/httptask"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
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
			group.POST("/getAuth", getAuth)
			group.POST("/grantAuth", grantAuth)
		}
		group = e.Group("/api/prop/content", apisession.CheckLogin)
		{
			group.POST("/insert", insertContent)
			group.POST("/list", listContent)
			group.POST("/update", updateContent)
			group.POST("/delete", deleteContent)
			group.POST("/deploy", deployContent)
			group.POST("/listHistory", listHistory)
			group.POST("/listDeploy", listDeploy)
		}
	})
	httptask.AppendHttpTask("checkPropDbEtcdConsistent", func(_ []byte, _ url.Values) {
		envs, b := cfgsrv.Inner.GetEnvCfg(context.Background())
		if b {
			for _, env := range envs {
				propsrv.Inner.CheckConsistent(env)
			}
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
			Env:       req.Env,
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
			Env:       req.Env,
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
			Env:      req.Env,
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
	var req ListEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		nodes, err := propsrv.Outer.ListEtcdNode(c, propsrv.ListEtcdNodeReqDTO{
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(nodes, func(t propsrv.EtcdNodeDTO) (EtcdNodeVO, error) {
			return EtcdNodeVO{
				NodeId:    t.NodeId,
				Endpoints: t.Endpoints,
				Username:  t.Username,
				Password:  t.Password,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]EtcdNodeVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func listSimpleEtcdNode(c *gin.Context) {
	var req ListSimpleEtcdNodeReqVO
	if util.ShouldBindJSON(&req, c) {
		nodes, err := propsrv.Outer.ListSimpleEtcdNode(c, req.Env)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     nodes,
		})
	}
}

func insertContent(c *gin.Context) {
	var req InsertContentReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.InsertPropContent(c, propsrv.InsertPropContentReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Content:  req.Content,
			Env:      req.Env,
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
			Env:      req.Env,
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
			Env:      req.Env,
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
			Env:          req.Env,
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
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(contents, func(t propsrv.PropContentDTO) (PropContentVO, error) {
			return PropContentVO{
				Id:    t.Id,
				AppId: t.AppId,
				Name:  t.Name,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]PropContentVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
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
			Env:       req.Env,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(histories, func(t propsrv.HistoryDTO) (HistoryVO, error) {
			return HistoryVO{
				ContentId: t.ContentId,
				Content:   t.Content,
				Version:   t.Version,
				Created:   t.Created.Format(time.DateTime),
				Creator:   t.Creator,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[[]HistoryVO]{
			DataResp: ginutil.DataResp[[]HistoryVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: cursor,
		})
	}
}

func grantAuth(c *gin.Context) {
	var req GrantAuthReqVO
	if util.ShouldBindJSON(&req, c) {
		err := propsrv.Outer.GrantAuth(c, propsrv.GrantAuthReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getAuth(c *gin.Context) {
	var req GetAuthReqVO
	if util.ShouldBindJSON(&req, c) {
		username, password, err := propsrv.Outer.GetAuth(c, propsrv.GetAuthReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, GetAuthRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Username: username,
			Password: password,
		})
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
			Env:       req.Env,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(deploys, func(t propsrv.DeployDTO) (DeployVO, error) {
			return DeployVO{
				ContentId: t.ContentId,
				Content:   t.Content,
				Version:   t.Version,
				NodeId:    t.NodeId,
				Created:   t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[[]DeployVO]{
			DataResp: ginutil.DataResp[[]DeployVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: cursor,
		})
	}
}
