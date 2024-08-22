package githook

import (
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/githook"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/rpcheader"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/v1/git/hook", checkHookToken)
		{
			group.POST("/pre-receive", preReceive)
			group.POST("/post-receive", postReceive)
		}
	})
}

func checkHookToken(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if authorization != git.HookToken() {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.UnauthorizedCode.Int(),
			Message: i18n.GetByKey(i18n.SystemUnauthorized),
		})
		c.Abort()
	} else {
		c.Next()
	}
}

func preReceive(c *gin.Context) {
	var req githook.Opts
	if util.ShouldBindJSON(&req, c) {
		err := doPreReceive(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func postReceive(c *gin.Context) {
	var req githook.Opts
	if ginutil.ShouldBind(&req, c) {
		newCtx, _ := rpcheader.NewCtxFromOldCtx(c)
		go doPostReceive(newCtx, req)
		util.DefaultOkResponse(c)
	}
}
