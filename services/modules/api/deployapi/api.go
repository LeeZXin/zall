package deployapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/deploy"
	"github.com/LeeZXin/zall/services/modules/service/deploysrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/deploy", apisession.CheckLogin)
		{
			group.POST("/get", getDeploy)
			group.POST("/update", updateDeploy)
		}
	})
}

func getDeploy(c *gin.Context) {
	var req GetDeployReqVO
	if util.ShouldBindJSON(&req, c) {
		config, err := deploysrv.Outer.GetDeploy(c, deploysrv.GetDeployReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[deploy.Config]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     config,
		})
	}
}

func updateDeploy(c *gin.Context) {
	var req UpdateDeployReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdateDeploy(c, deploysrv.UpdateDeployReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Config:   req.Config,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
