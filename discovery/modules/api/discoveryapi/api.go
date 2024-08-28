package discoveryapi

import (
	"github.com/LeeZXin/zall/discovery/modules/service/discoverysrv"
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
		group := e.Group("/api/discoverySource", apisession.CheckLogin)
		{
			// 列表
			group.GET("/list", listDiscoverySource)
			// 创建
			group.POST("/create", createDiscoverySource)
			// 编辑
			group.POST("/update", updateDiscoverySource)
			// 删除
			group.DELETE("/delete/:sourceId", deleteDiscoverySource)
			// 所有来源列表
			group.GET("/listAll/:env", listAllDiscoverySource)
			// 绑定来源列表
			group.GET("/listBind", listBindDiscoverySource)
			// 绑定来源
			group.POST("/bindApp", bindAppAndDiscoverySource)
		}
		group = e.Group("/api/discoveryService", apisession.CheckLogin)
		{
			// 获取服务列表
			group.GET("/listService/:bindId", listDiscoveryService)
			// 下线服务
			group.POST("/deregister", deregisterService)
			// 上线服务
			group.POST("/reRegister", reRegisterService)
			// 删除下线服务
			group.DELETE("/deleteDownService", deleteDownService)
		}
	})
}

func listDiscoverySource(c *gin.Context) {
	var req ListDiscoverySourceReqVO
	if util.ShouldBindQuery(&req, c) {
		sources, err := discoverysrv.ListDiscoverySource(c, discoverysrv.ListDiscoverySourceReqDTO{
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(sources, func(t discoverysrv.DiscoverySourceDTO) DiscoverySourceVO {
			return DiscoverySourceVO{
				Id:        t.Id,
				Name:      t.Name,
				Endpoints: t.Endpoints,
				Username:  t.Username,
				Password:  t.Password,
				Env:       t.Env,
			}
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]DiscoverySourceVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func listBindDiscoverySource(c *gin.Context) {
	var req ListBindDiscoverySourceReqVO
	if util.ShouldBindQuery(&req, c) {
		sources, err := discoverysrv.ListBindDiscoverySource(c, discoverysrv.ListBindDiscoverySourceReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(sources, func(t discoverysrv.SimpleBindDiscoverySourceDTO) SimpleBindDiscoverySourceVO {
			return SimpleBindDiscoverySourceVO{
				Id:     t.Id,
				Name:   t.Name,
				BindId: t.BindId,
				Env:    t.Env,
			}
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleBindDiscoverySourceVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func createDiscoverySource(c *gin.Context) {
	var req CreateDiscoverySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.CreateDiscoverySource(c, discoverysrv.CreateDiscoverySourceReqDTO{
			Name:      req.Name,
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

func updateDiscoverySource(c *gin.Context) {
	var req UpdateDiscoverySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.UpdateDiscoverySource(c, discoverysrv.UpdateDiscoverySourceReqDTO{
			SourceId:  req.SourceId,
			Name:      req.Name,
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

func deleteDiscoverySource(c *gin.Context) {
	err := discoverysrv.DeleteDiscoverySource(c, discoverysrv.DeleteDiscoverySourceReqDTO{
		SourceId: cast.ToInt64(c.Param("sourceId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listDiscoveryService(c *gin.Context) {
	services, err := discoverysrv.ListDiscoveryService(c, discoverysrv.ListDiscoveryServiceReqDTO{
		BindId:   cast.ToInt64(c.Param("bindId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(services, func(t discoverysrv.ServiceDTO) ServiceVO {
		return ServiceVO{
			Server:     t.Server,
			Up:         t.Up,
			InstanceId: t.InstanceId,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func deregisterService(c *gin.Context) {
	var req DeregisterServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.DeregisterService(c, discoverysrv.DeregisterServiceReqDTO{
			BindId:     req.BindId,
			InstanceId: req.InstanceId,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func reRegisterService(c *gin.Context) {
	var req ReRegisterServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.ReRegisterService(c, discoverysrv.ReRegisterServiceReqDTO{
			BindId:     req.BindId,
			InstanceId: req.InstanceId,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteDownService(c *gin.Context) {
	var req DeleteDownServiceReqVO
	if util.ShouldBindQuery(&req, c) {
		err := discoverysrv.DeleteDownService(c, discoverysrv.DeleteDownServiceReqDTO{
			BindId:     req.BindId,
			InstanceId: req.InstanceId,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func bindAppAndDiscoverySource(c *gin.Context) {
	var req BindAppAndDiscoverySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.BindAppAndDiscoverySource(c, discoverysrv.BindAppAndDiscoverySourceReqDTO{
			AppId:        req.AppId,
			SourceIdList: req.SourceIdList,
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

func listAllDiscoverySource(c *gin.Context) {
	nodes, err := discoverysrv.ListAllDiscoverySource(c, discoverysrv.ListAllDiscoverySourceReqDTO{
		Env:      c.Param("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(nodes, func(t discoverysrv.SimpleDiscoverySourceDTO) SimpleDiscoverySourceVO {
		return SimpleDiscoverySourceVO{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleDiscoverySourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
