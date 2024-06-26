package deployapi

import (
	"github.com/LeeZXin/zall/deploy/modules/service/deploysrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func InitApi() {
	deploysrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/deployPlan", apisession.CheckLogin)
		{
			// 创建发布计划
			group.POST("/create", createPlan)
			// 关闭发布计划
			group.PUT("/close/:planId", closePlan)
			// 展示发布计划列表
			group.GET("/list", listPlan)
			// 开始发布计划
			group.PUT("/start/:planId", startPlan)
			// 服务列表
			group.GET("/listService", listServiceWhenCreatePlan)
		}
		group = e.Group("/api/service", apisession.CheckLogin)
		{
			// 创建服务
			group.POST("/create", createService)
			// 编辑服务
			group.POST("/update", updateService)
			// 删除服务
			group.DELETE("/delete/:serviceId", deleteService)
			// 服务列表
			group.GET("/list", listService)
		}
	})
}

func createPlan(c *gin.Context) {
	var req CreatePlanReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreatePlan(c, deploysrv.CreatePlanReqDTO{
			Name:           req.Name,
			ServiceId:      req.ServiceId,
			ProductVersion: req.ProductVersion,
			Operator:       apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func closePlan(c *gin.Context) {
	err := deploysrv.Outer.ClosePlan(c, deploysrv.ClosePlanReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func startPlan(c *gin.Context) {
	err := deploysrv.Outer.StartPlan(c, deploysrv.StartPlanReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listPlan(c *gin.Context) {
	var req ListPlanReqVO
	if util.ShouldBindQuery(&req, c) {
		plans, total, err := deploysrv.Outer.ListPlan(c, deploysrv.ListPlanReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(plans, func(t deploysrv.PlanDTO) (PlanVO, error) {
			return PlanVO{
				Id:             t.Id,
				ServiceId:      t.ServiceId,
				ServiceName:    t.ServiceName,
				Name:           t.Name,
				ProductVersion: t.ProductVersion,
				PlanStatus:     t.PlanStatus,
				Env:            t.Env,
				Creator:        t.Creator,
				Created:        t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[PlanVO]{
			DataResp: ginutil.DataResp[[]PlanVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func createService(c *gin.Context) {
	var req CreateServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreateService(c, deploysrv.CreateServiceReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Config:   req.Config,
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

func updateService(c *gin.Context) {
	var req UpdateServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdateService(c, deploysrv.UpdateServiceReqDTO{
			ServiceId: req.ServiceId,
			Name:      req.Name,
			Config:    req.Config,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteService(c *gin.Context) {
	err := deploysrv.Outer.DeleteService(c, deploysrv.DeleteServiceReqDTO{
		ServiceId: cast.ToInt64(c.Param("serviceId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listService(c *gin.Context) {
	services, err := deploysrv.Outer.ListService(c, deploysrv.ListServiceReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(services, func(t deploysrv.ServiceDTO) (ServiceVO, error) {
		return ServiceVO{
			Id:          t.Id,
			AppId:       t.AppId,
			Config:      t.Config,
			Env:         t.Env,
			Name:        t.Name,
			ServiceType: t.ServiceType,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listServiceWhenCreatePlan(c *gin.Context) {
	services, err := deploysrv.Outer.ListServiceWhenCreatePlan(c, deploysrv.ListServiceWhenCreatePlanReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(services, func(t deploysrv.SimpleServiceDTO) (SimpleServiceVO, error) {
		return SimpleServiceVO{
			Id:          t.Id,
			Env:         t.Env,
			Name:        t.Name,
			ServiceType: t.ServiceType,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleServiceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
