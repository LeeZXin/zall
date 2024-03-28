package deployapi

import (
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/services/modules/service/deploysrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	deployToken string
)

func InitDeploy() {
	deployToken = static.GetString("deploy.token")
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/innerDeploy", checkToken)
		{
			group.POST("/deployWithoutPlan", deployWithoutPlan)
		}
	})
}

func InitApi() {
	InitDeploy()
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

		group = e.Group("/api/deployService", apisession.CheckLogin)
		{
			group.POST("/list", listService)
			group.POST("/deployWithPlan")
			group.POST("/stop", stopService)
			group.POST("/reDeploy", reDeployService)
			group.POST("/restart", restartService)
		}

	})
}

func checkToken(c *gin.Context) {
	if c.Query("t") != deployToken {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.UnauthorizedCode.Int(),
			Message: "invalid token",
		})
		c.Abort()
	}
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
			AppId:          req.AppId,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func reDeployService(c *gin.Context) {
	var req ReDeployServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.ReDeployService(c, deploysrv.ReDeployServiceReqDTO{
			ConfigId: req.ConfigId,
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

func stopService(c *gin.Context) {
	var req StopServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.StopService(c, deploysrv.StopServiceReqDTO{
			ConfigId: req.ConfigId,
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

func restartService(c *gin.Context) {
	var req RestartServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.RestartService(c, deploysrv.RestartServiceReqDTO{
			ConfigId: req.ConfigId,
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

func listService(c *gin.Context) {
	var req ListServiceReqVO
	if util.ShouldBindJSON(&req, c) {
		services, err := deploysrv.Outer.ListService(c, deploysrv.ListServiceReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(services, func(t deploysrv.ServiceDTO) (ServiceVO, error) {
			return ServiceVO{
				CurrProductVersion: t.CurrProductVersion,
				LastProductVersion: t.LastProductVersion,
				ServiceType:        t.ServiceType.Readable(),
				ProcessConfig:      t.ProcessConfig,
				K8sConfig:          t.K8sConfig,
				ActiveStatus:       t.ActiveStatus.Readable(),
				StartTime:          time.UnixMilli(t.StartTime).Format(time.DateTime),
				ProbeTime:          time.UnixMilli(t.ProbeTime).Format(time.DateTime),
				Created:            t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}
