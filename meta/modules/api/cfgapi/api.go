package cfgapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	{
		// 初始化全局配置
		cfgsrv.Inner.InitSysCfg()
	}
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sysCfg")
		{
			group.GET("/get", getSysCfg)
			group.POST("/update", apisession.CheckLogin, updateSysCfg)
		}
		group = e.Group("/api/gitCfg", apisession.CheckLogin)
		{
			group.GET("/get", getGitCfg)
			group.POST("/update", updateGitCfg)
		}
		group = e.Group("/api/envCfg", apisession.CheckLogin)
		{
			// 获取环境列表
			group.GET("/get", getEnvCfg)
			group.POST("/update", updateEnvCfg)
		}
		group = e.Group("/api/gitRepoServerCfg", apisession.CheckLogin)
		{
			// 获取git仓库服务url
			group.GET("/get", getGitRepoServerCfg)
			// 编辑git仓库服务url
			group.POST("/update", updateGitRepoServerCfg)
		}
		group = e.Group("/api/zonesCfg", apisession.CheckLogin)
		{
			// 获取单元调用列表
			group.GET("/get", getZonesCfg)
			// 编辑git仓库服务url
			group.POST("/update", updateZonesCfg)

		}
	})
}

func getZonesCfg(c *gin.Context) {
	zones, err := cfgsrv.Outer.GetZonesCfg(c, cfgsrv.GetZonesCfgReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     zones,
	})
}

func updateZonesCfg(c *gin.Context) {
	var req UpdateZonesCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.Outer.UpdateZonesCfg(c, cfgsrv.UpdateZonesCfgReqDTO{
			ZonesCfg: req.ZonesCfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getSysCfg(c *gin.Context) {
	cfg, err := cfgsrv.Outer.GetSysCfg(c)
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[cfgsrv.SysCfg]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     cfg,
	})
}

func updateSysCfg(c *gin.Context) {
	var req UpdateSysCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.Outer.UpdateSysCfg(c, cfgsrv.UpdateSysCfgReqDTO{
			SysCfg:   req.SysCfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getGitCfg(c *gin.Context) {
	cfg, err := cfgsrv.Outer.GetGitCfg(c, cfgsrv.GetGitCfgReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[cfgsrv.GitCfg]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     cfg,
	})
}

func updateGitCfg(c *gin.Context) {
	var req UpdateGitCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.Outer.UpdateGitCfg(c, cfgsrv.UpdateGitCfgReqDTO{
			GitCfg:   req.GitCfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getEnvCfg(c *gin.Context) {
	cfg, err := cfgsrv.Outer.GetEnvCfg(c, cfgsrv.GetEnvCfgReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     cfg,
	})
}

func updateEnvCfg(c *gin.Context) {
	var req UpdateEnvCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.Outer.UpdateEnvCfg(c, cfgsrv.UpdateEnvCfgReqDTO{
			Envs:     req.Envs,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getGitRepoServerCfg(c *gin.Context) {
	cfg, err := cfgsrv.Outer.GetGitRepoServerCfg(c, cfgsrv.GetGitRepoServerUrlReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[cfgsrv.GitRepoServerCfg]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     cfg,
	})
}

func updateGitRepoServerCfg(c *gin.Context) {
	var req UpdateGitRepoServerCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.Outer.UpdateGitRepoServerCfg(c, cfgsrv.UpdateGitRepoServerCfgReqDTO{
			GitRepoServerCfg: req.GitRepoServerCfg,
			Operator:         apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
