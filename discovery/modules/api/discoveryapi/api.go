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
	discoverysrv.Init()
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
		}
		group = e.Group("/api/discoveryService", apisession.CheckLogin)
		{
			// 来源列表
			group.GET("/listSource", listSimpleDiscoverySource)
			// 获取服务列表
			group.GET("/listService/:sourceId", listDiscoveryService)
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
		sources, err := discoverysrv.Outer.ListDiscoverySource(c, discoverysrv.ListDiscoverySourceReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(sources, func(t discoverysrv.DiscoverySourceDTO) (DiscoverySourceVO, error) {
			return DiscoverySourceVO{
				Id:        t.Id,
				Name:      t.Name,
				Endpoints: t.Endpoints,
				Username:  t.Username,
				Password:  t.Password,
				Env:       t.Env,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]DiscoverySourceVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func listSimpleDiscoverySource(c *gin.Context) {
	var req ListDiscoverySourceReqVO
	if util.ShouldBindQuery(&req, c) {
		sources, err := discoverysrv.Outer.ListSimpleDiscoverySource(c, discoverysrv.ListDiscoverySourceReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(sources, func(t discoverysrv.SimpleDiscoverySourceDTO) (SimpleDiscoverySourceVO, error) {
			return SimpleDiscoverySourceVO{
				Id:   t.Id,
				Name: t.Name,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleDiscoverySourceVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func createDiscoverySource(c *gin.Context) {
	var req CreateDiscoverySourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.Outer.CreateDiscoverySource(c, discoverysrv.CreateDiscoverySourceReqDTO{
			AppId:     req.AppId,
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
		err := discoverysrv.Outer.UpdateDiscoverySource(c, discoverysrv.UpdateDiscoverySourceReqDTO{
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
	err := discoverysrv.Outer.DeleteDiscoverySource(c, discoverysrv.DeleteDiscoverySourceReqDTO{
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
	services, err := discoverysrv.Outer.ListDiscoveryService(c, discoverysrv.ListDiscoveryServiceReqDTO{
		SourceId: cast.ToInt64(c.Param("sourceId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(services, func(t discoverysrv.ServiceDTO) (ServiceVO, error) {
		return ServiceVO{
			Server:     t.Server,
			Up:         t.Up,
			InstanceId: t.InstanceId,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func deregisterService(c *gin.Context) {
	var req DeregisterServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := discoverysrv.Outer.DeregisterService(c, discoverysrv.DeregisterServiceReqDTO{
			SourceId:   req.SourceId,
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
		err := discoverysrv.Outer.ReRegisterService(c, discoverysrv.ReRegisterServiceReqDTO{
			SourceId:   req.SourceId,
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
		err := discoverysrv.Outer.DeleteDownService(c, discoverysrv.DeleteDownServiceReqDTO{
			SourceId:   req.SourceId,
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
