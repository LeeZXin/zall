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
	cfgsrv.Inner.InitSysCfg()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sysCfg", apisession.CheckLogin)
		{
			group.GET("/get", getSysCfg)
			group.POST("/update", updateSysCfg)
		}
		group = e.Group("/api/gitCfg", apisession.CheckLogin)
		{
			group.GET("/get", getGitCfg)
			group.POST("/update", updateGitCfg)
		}
		group = e.Group("/api/envCfg", apisession.CheckLogin)
		{
			group.GET("/get", getEnvCfg)
			group.POST("/update", updateEnvCfg)
		}
	})
}

func getSysCfg(c *gin.Context) {
	cfg, err := cfgsrv.Outer.GetSysCfg(c, cfgsrv.GetSysCfgReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, GetSysCfgRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Cfg:      cfg,
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
	c.JSON(http.StatusOK, GetGitCfgRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Cfg:      cfg,
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
	c.JSON(http.StatusOK, GetEnvCfgRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Cfg:      cfg,
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
