package deployapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/services/modules/service/deploysrv"
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
		group := e.Group("/api/deployConfig", apisession.CheckLogin)
		{
			group.POST("/list", listConfig)
			group.POST("/update", updateConfig)
			group.POST("/insert", insertConfig)
		}
		group = e.Group("/api/deployPlan")
		{
			group.POST("/insert", insertPlan)
		}

		group = e.Group("/api/deployService")
		{
			group.POST("/list", apisession.CheckLogin)
			group.POST("/deployWithPlan", apisession.CheckLogin)
			group.POST("/deployWithoutPlan", deployWithoutPlan)
			group.POST("/shutdown", apisession.CheckLogin)
			group.POST("/reDeploy", apisession.CheckLogin)
		}

	})
}

func listConfig(c *gin.Context) {
	var req ListConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		configs, err := deploysrv.Outer.ListConfig(c, deploysrv.ListConfigReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(configs, func(t deploysrv.ConfigDTO) (ConfigVO, error) {
			return ConfigVO{
				Id:            t.Id,
				AppId:         t.AppId,
				Name:          t.Name,
				ServiceType:   t.ServiceType.Readable(),
				ProcessConfig: t.ProcessConfig,
				K8sConfig:     t.K8sConfig,
				Created:       t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]ConfigVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func updateConfig(c *gin.Context) {
	var req UpdateConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdateConfig(c, deploysrv.UpdateConfigReqDTO{
			ConfigId:      req.ConfigId,
			Env:           req.Env,
			ProcessConfig: req.ProcessConfig,
			K8sConfig:     req.K8sConfig,
			Name:          req.Name,
			Operator:      apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func insertConfig(c *gin.Context) {
	var req InsertConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.InsertConfig(c, deploysrv.InsertConfigReqDTO{
			AppId:         req.AppId,
			Name:          req.Name,
			ServiceType:   req.ServiceType,
			ProcessConfig: req.ProcessConfig,
			K8sConfig:     req.K8sConfig,
			Env:           req.Env,
			Operator:      apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func insertPlan(c *gin.Context) {
	var req InsertPlanReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.InsertPlan(c, deploysrv.InsertPlanReqDTO{
			Name:     req.Name,
			TeamId:   req.TeamId,
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

func deployWithoutPlan(c *gin.Context) {
	var req DeployServiceWithoutPlanReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Inner.DeployServiceWithoutPlan(c, deploysrv.DeployServiceWithoutPlanReqDTO{
			ConfigId:       req.ConfigId,
			Env:            req.Env,
			ProductVersion: req.ProductVersion,
			Operator:       req.Operator,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
