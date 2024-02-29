package apisession

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context) {
	cfg, b := cfgsrv.Inner.GetGitCfg(c)
	if !b || c.GetHeader("Authorization") != cfg.RepoToken {
		c.String(http.StatusUnauthorized, "repo token failed")
		c.Abort()
		return
	}
	c.Next()
}
