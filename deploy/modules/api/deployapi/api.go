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
		group := e.Group("/api/deployConfig", apisession.CheckLogin)
		{
			// 部署配置列表
			group.GET("/list", listConfig)
			// 编辑部署配置
			group.POST("/update", updateConfig)
			// 新增部署配置
			group.POST("/create", createConfig)
			// 删除配置
			group.DELETE("/delete/:configId", deleteConfig)
		}

		group = e.Group("/api/deployPlan", apisession.CheckLogin)
		{
			// 创建发布计划
			group.POST("/create", createPlan)
			// 关闭发布计划
			group.PUT("/close/:planId", closePlan)
			// 展示发布计划列表
			group.GET("/list", listPlan)
		}

		group = e.Group("/api/deployPlanService", apisession.CheckLogin)
		{
			// 添加发布的服务
			group.POST("/add", addPlanService)
			// 删除未执行服务
			group.DELETE("/deletePending/:serviceId", deletePendingService)
			// 发布计划服务列表
			group.GET("/list/:planId", listPlanService)
			// 启动部署服务
			group.PUT("/start/:serviceId", startPlanService)
		}
	})
}

func listConfig(c *gin.Context) {
	configs, err := deploysrv.Outer.ListConfig(c, deploysrv.ListConfigReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(configs, func(t deploysrv.ConfigDTO) (ConfigVO, error) {
		return ConfigVO{
			Id:      t.Id,
			AppId:   t.AppId,
			Name:    t.Name,
			Content: t.Content,
			Env:     t.Env,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ConfigVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func updateConfig(c *gin.Context) {
	var req UpdateConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdateConfig(c, deploysrv.UpdateConfigReqDTO{
			ConfigId: req.ConfigId,
			Content:  req.Content,
			Name:     req.Name,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func createConfig(c *gin.Context) {
	var req CreateConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreateConfig(c, deploysrv.CreateConfigReqDTO{
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

func deleteConfig(c *gin.Context) {
	err := deploysrv.Outer.DeleteConfig(c, deploysrv.DeleteConfigReqDTO{
		ConfigId: cast.ToInt64(c.Param("configId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func createPlan(c *gin.Context) {
	var req CreatePlanReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreatePlan(c, deploysrv.CreatePlanReqDTO{
			Name:        req.Name,
			TeamId:      req.TeamId,
			Env:         req.Env,
			ExpireHours: req.ExpireHours,
			Operator:    apisession.MustGetLoginUser(c),
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

func addPlanService(c *gin.Context) {
	var req AddDeployPlanServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.AddPlanService(c, deploysrv.AddPlanServiceReqDTO{
			PlanId:             req.PlanId,
			ConfigId:           req.ConfigId,
			LastProductVersion: req.LastProductVersion,
			CurrProductVersion: req.CurrProductVersion,
			Operator:           apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deletePendingService(c *gin.Context) {
	err := deploysrv.Outer.DeletePendingPlanService(c, deploysrv.DeletePendingPlanServiceReqDTO{
		ServiceId: cast.ToInt64(c.Param("serviceId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func startPlanService(c *gin.Context) {
	err := deploysrv.Outer.StartPlanService(c, deploysrv.StartPlanServiceReqDTO{
		ServiceId: cast.ToInt64(c.Param("serviceId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listPlanService(c *gin.Context) {
	items, err := deploysrv.Outer.ListPlanService(c, deploysrv.ListPlanServiceReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(items, func(t deploysrv.PlanServiceDTO) (PlanServiceVO, error) {
		return PlanServiceVO{
			Id:                 t.Id,
			AppId:              t.AppId,
			ConfigId:           t.ConfigId,
			ConfigName:         t.ConfigName,
			CurrProductVersion: t.CurrProductVersion,
			LastProductVersion: t.LastProductVersion,
			ServiceStatus:      t.ServiceStatus,
			Created:            t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]PlanServiceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listPlan(c *gin.Context) {
	var req ListPlanReqVO
	if util.ShouldBindQuery(&req, c) {
		plans, next, err := deploysrv.Outer.ListPlan(c, deploysrv.ListPlanReqDTO{
			TeamId:   req.TeamId,
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
				Id:       t.Id,
				Name:     t.Name,
				IsClosed: t.IsClosed,
				TeamId:   t.TeamId,
				Creator:  t.Creator,
				Expired:  t.Expired.Format(time.DateTime),
				Created:  t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[PlanVO]{
			DataResp: ginutil.DataResp[[]PlanVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}
