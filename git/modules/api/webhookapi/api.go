package webhookapi

import (
	"github.com/LeeZXin/zall/git/modules/service/webhooksrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/webhook", apisession.CheckLogin)
		{
			group.GET("/list", listWebhook)
			group.POST("/insert", insertWebhook)
			group.POST("/update", updateWebhook)
			group.POST("/delete", deleteWebhook)
		}
	})
}

func insertWebhook(c *gin.Context) {
	var req InsertWebhookReqVO
	if util.ShouldBindJSON(&req, c) {
		err := webhooksrv.Outer.InsertWebhook(c, webhooksrv.InsertWebhookReqDTO{
			RepoId:      req.RepoId,
			HookUrl:     req.HookUrl,
			HttpHeaders: req.Headers,
			HookType:    req.HookType,
			WildBranch:  req.WildBranch,
			WildTag:     req.WildTag,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateWebhook(c *gin.Context) {
	var req UpdateWebhookReqVO
	if util.ShouldBindJSON(&req, c) {
		err := webhooksrv.Outer.UpdateWebhook(c, webhooksrv.UpdateWebhookReqDTO{
			Id:          req.Id,
			HookUrl:     req.HookUrl,
			HttpHeaders: req.HttpHeaders,
			HookType:    req.HookType,
			WildBranch:  req.WildBranch,
			WildTag:     req.WildTag,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteWebhook(c *gin.Context) {
	var req DeleteWebhookReqVO
	if util.ShouldBindJSON(&req, c) {
		err := webhooksrv.Outer.DeleteWebhook(c, webhooksrv.DeleteWebhookReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listWebhook(c *gin.Context) {
	var req ListWebhookReqVO
	if util.ShouldBindJSON(&req, c) {
		hookList, err := webhooksrv.Outer.ListWebhook(c, webhooksrv.ListWebhookReqDTO{
			RepoId:   req.RepoId,
			HookType: req.HookType,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(hookList, func(t webhooksrv.WebhookDTO) (WebhookVO, error) {
			return WebhookVO{
				Id:          t.Id,
				RepoId:      t.RepoId,
				HookUrl:     t.HookUrl,
				HttpHeaders: t.HttpHeaders,
				HookType:    t.HookType.Readable(),
				WildBranch:  t.WildBranch,
				WildTag:     t.WildTag,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]WebhookVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}
