package actionapi

import (
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	authToken = static.GetString("actions.token")
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/actions", checkToken)
		{
			group.POST("/git", gitAction)
		}
	})
}

func gitAction(c *gin.Context) {
	var req action.Webhook
	if ginutil.ShouldBind(&req, c) {
		go actionsrv.TriggerGitAction(c, req)
		c.String(http.StatusOK, "received")
	}
}

func checkToken(c *gin.Context) {
	if c.GetHeader("Authorization") != authToken {
		c.String(http.StatusForbidden, "invalid token")
		c.Abort()
		return
	}
	c.Next()
}
