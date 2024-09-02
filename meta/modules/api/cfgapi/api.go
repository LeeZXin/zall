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
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sysCfg")
		{
			// 获取系统配置 无需鉴权
			group.GET("/get", getSysCfg)
			// 编辑系统配置
			group.POST("/update", apisession.CheckLogin, updateSysCfg)
		}
		group = e.Group("/api/gitCfg", apisession.CheckLogin)
		{
			// 获取git配置
			group.GET("/get", getGitCfg)
			// 编辑git配置
			group.POST("/update", updateGitCfg)
		}
		group = e.Group("/api/envCfg", apisession.CheckLogin)
		{
			// 获取环境列表
			group.GET("/get", getEnvCfg)
			// 编辑环境列表
			group.POST("/update", updateEnvCfg)
		}
		group = e.Group("/api/gitRepoServerCfg", apisession.CheckLogin)
		{
			// 获取git仓库服务url
			group.GET("/get", getGitRepoServerCfg)
			// 编辑git仓库服务url
			group.POST("/update", updateGitRepoServerCfg)
		}
		group = e.Group("/api/loginCfg")
		{
			group.GET("/get", getLoginCfg)
			// 超级管理员获取登录配置
			group.GET("/getBySa", apisession.CheckLogin, getLoginCfgBySa)
			// 编辑登录配置
			group.POST("/update", apisession.CheckLogin, updateLoginCfg)
		}
	})
}

func getSysCfg(c *gin.Context) {
	cfg, err := cfgsrv.GetSysCfg(c)
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
		err := cfgsrv.UpdateSysCfg(c, cfgsrv.UpdateSysCfgReqDTO{
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
	cfg, err := cfgsrv.GetGitCfg(c, cfgsrv.GetGitCfgReqDTO{
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
		err := cfgsrv.UpdateGitCfg(c, cfgsrv.UpdateGitCfgReqDTO{
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
	cfg, err := cfgsrv.GetEnvCfg(c, cfgsrv.GetEnvCfgReqDTO{
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
		err := cfgsrv.UpdateEnvCfg(c, cfgsrv.UpdateEnvCfgReqDTO{
			EnvCfg:   req.EnvCfg,
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
	cfg, err := cfgsrv.GetGitRepoServerCfg(c, cfgsrv.GetGitRepoServerUrlReqDTO{
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

func getLoginCfgBySa(c *gin.Context) {
	cfg, err := cfgsrv.GetLoginCfgBySa(c, cfgsrv.GetLoginCfgBySaReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[cfgsrv.LoginCfg]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     cfg,
	})
}

func getLoginCfg(c *gin.Context) {
	cfg, err := cfgsrv.GetLoginCfg(c)
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[cfgsrv.LoginCfg]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     cfg,
	})
}

func updateGitRepoServerCfg(c *gin.Context) {
	var req UpdateGitRepoServerCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.UpdateGitRepoServerCfg(c, cfgsrv.UpdateGitRepoServerCfgReqDTO{
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

func updateLoginCfg(c *gin.Context) {
	var req UpdateLoginCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.UpdateLoginCfg(c, cfgsrv.UpdateLoginCfgReqDTO{
			LoginCfg: req.LoginCfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
