package statusapi

import (
	"github.com/LeeZXin/zall/deploy/modules/service/statussrv"
	"github.com/LeeZXin/zall/pkg/status"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	apiKey string
)

func InitApi() {
	statussrv.Init()
	apiKey = static.GetString("zallet.status.api.key")
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/service/v1/status", checkAuthorization)
		{
			group.GET("/actions", getActions)
			group.GET("/list", listService)
			group.PUT("/kill", killService)
			group.PUT("/restart", restartService)
		}
	})
}

func killService(c *gin.Context) {
	err := statussrv.Outer.KillService(c, c.Query("serviceId"))
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func restartService(c *gin.Context) {
	err := statussrv.Outer.RestartService(c, c.Query("serviceId"))
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func listService(c *gin.Context) {
	var req status.ListServiceReq
	if shouldBindQuery(&req, c) {
		srvs, err := statussrv.Outer.ListService(c, req)
		if err != nil {
			c.String(http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.JSON(http.StatusOK, srvs)
	}
}

func getActions(c *gin.Context) {
	c.JSON(http.StatusOK, statussrv.Outer.GetActions(c))
}

func shouldBindQuery(obj any, c *gin.Context) bool {
	err := ginutil.BindQuery(c, obj)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return false
	}
	return true
}

func checkAuthorization(c *gin.Context) {
	if apiKey != c.GetHeader("Authorization") {
		c.String(http.StatusUnauthorized, "invalid api key")
		c.Abort()
		return
	}
	c.Next()
}
