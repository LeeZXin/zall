package actionagentapi

import (
	"github.com/LeeZXin/zall/git/modules/service/gitactionsrv"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/actionAgent", checkToken)
		{
			group.POST("/execute", executeAction)
		}
	})
}

func executeAction(c *gin.Context) {
	var req action.Hook
	if ginutil.ShouldBind(&req, c) {
		go gitactionsrv.Inner.ExecuteAction(c, req)
		c.String(http.StatusOK, "received")
	}
}

func checkToken(c *gin.Context) {
	cfg, b := cfgsrv.Inner.GetGitCfg(c)
	if !b {
		c.String(http.StatusInternalServerError, "can not get git config")
		c.Abort()
		return
	}
	if c.GetHeader("Authorization") != cfg.ActionToken {
		c.String(http.StatusForbidden, "invalid token")
		c.Abort()
		return
	}
	c.Next()
}
