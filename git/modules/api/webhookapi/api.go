package webhookapi

import (
	"github.com/LeeZXin/zall/git/modules/service/webhooksrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/webhook", apisession.CheckLogin)
		{
			// 展示webhook列表
			group.GET("/list/:repoId", listWebhook)
			// 新增
			group.POST("/create", createWebhook)
			// 编辑
			group.POST("/update", updateWebhook)
			// 删除
			group.DELETE("/delete/:webhookId", deleteWebhook)
			// ping
			group.PUT("/ping/:webhookId", pingWebhook)
		}
		group = e.Group("/api/webhook")
		{
			group.POST("/demo", demoWebhook)
		}
	})
}

func demoWebhook(c *gin.Context) {
	defer c.Request.Body.Close()
	logger.Logger.WithContext(c).Info("event-type: ", c.GetHeader(webhook.XEventType))
	signature := c.GetHeader(webhook.XSignature)
	logger.Logger.WithContext(c).Info("signature: ", signature)
	all, _ := io.ReadAll(c.Request.Body)
	mySecret := "hello world"
	mySignature := webhook.CreateSignature(all, mySecret)
	logger.Logger.WithContext(c).Info("mySignature: ", mySignature)
	logger.Logger.WithContext(c).Infof("mySignature is equal with signature = %v", mySignature == signature)
	logger.Logger.WithContext(c).Info("request: ", string(all))
	c.String(http.StatusOK, "ok")
}

func createWebhook(c *gin.Context) {
	var req CreateWebhookReqVO
	if util.ShouldBindJSON(&req, c) {
		err := webhooksrv.CreateWebhook(c, webhooksrv.CreateWebhookReqDTO{
			RepoId:   req.RepoId,
			HookUrl:  req.HookUrl,
			Secret:   req.Secret,
			Events:   req.Events,
			Operator: apisession.MustGetLoginUser(c),
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
		err := webhooksrv.UpdateWebhook(c, webhooksrv.UpdateWebhookReqDTO{
			Id:       req.WebhookId,
			HookUrl:  req.HookUrl,
			Secret:   req.Secret,
			Events:   req.Events,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteWebhook(c *gin.Context) {
	err := webhooksrv.DeleteWebhook(c, webhooksrv.DeleteWebhookReqDTO{
		Id:       cast.ToInt64(c.Param("webhookId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func pingWebhook(c *gin.Context) {
	err := webhooksrv.PingWebhook(c, webhooksrv.PingWebhookReqDTO{
		WebhookId: cast.ToInt64(c.Param("webhookId")),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listWebhook(c *gin.Context) {
	hookList, err := webhooksrv.ListWebhook(c, webhooksrv.ListWebhookReqDTO{
		RepoId:   cast.ToInt64(c.Param("repoId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(hookList, func(t webhooksrv.WebhookDTO) (WebhookVO, error) {
		return WebhookVO{
			Id:      t.Id,
			RepoId:  t.RepoId,
			HookUrl: t.HookUrl,
			Secret:  t.Secret,
			Events:  t.Events,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]WebhookVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
