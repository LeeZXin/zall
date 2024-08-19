package notifyapi

import (
	"github.com/LeeZXin/zall/notify/modules/service/notifysrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func InitApi() {
	notifysrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/notifyTpl", apisession.CheckLogin)
		{
			// 创建通知模版
			group.POST("/create", createNotifyTpl)
			// 编辑通知模版
			group.POST("/update", updateNotifyTpl)
			// 通知模版列表
			group.GET("/list", listNotifyTpl)
			// 删除通知模版
			group.DELETE("/delete/:tplId", deleteNotifyTpl)
			// 更换通知模版api key
			group.PUT("/changeApiKey/:tplId", changeNotifyTplApiKey)
			// 所有模板
			group.GET("/listAll/:teamId", listAllTplByTeamId)
		}
		group = e.Group("/api/notify")
		{
			group.POST("/send/:apiKey", sendNotification)
		}
	})
}

func createNotifyTpl(c *gin.Context) {
	var req CreateNotifyTplReqVO
	if util.ShouldBindJSON(&req, c) {
		err := notifysrv.Outer.CreateTpl(c, notifysrv.CreateTplReqDTO{
			Name:     req.Name,
			Cfg:      req.Cfg,
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateNotifyTpl(c *gin.Context) {
	var req UpdateNotifyTplReqVO
	if util.ShouldBindJSON(&req, c) {
		err := notifysrv.Outer.UpdateTpl(c, notifysrv.UpdateTplReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			Cfg:      req.Cfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteNotifyTpl(c *gin.Context) {
	err := notifysrv.Outer.DeleteTpl(c, notifysrv.DeleteTplReqDTO{
		Id:       cast.ToInt64(c.Param("tplId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func changeNotifyTplApiKey(c *gin.Context) {
	err := notifysrv.Outer.ChangeTplApiKey(c, notifysrv.ChangeTplApiKeyReqDTO{
		Id:       cast.ToInt64(c.Param("tplId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listAllTplByTeamId(c *gin.Context) {
	tpls, err := notifysrv.Outer.ListAllTpl(c, notifysrv.ListAllTplReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(tpls, func(t notifysrv.SimpleTplDTO) (SimpleTplVO, error) {
		return SimpleTplVO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleTplVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listNotifyTpl(c *gin.Context) {
	var req ListNotifyTplReqVO
	if util.ShouldBindQuery(&req, c) {
		tpls, total, err := notifysrv.Outer.ListTpl(c, notifysrv.ListTplReqDTO{
			Name:     req.Name,
			PageNum:  req.PageNum,
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(tpls, func(t notifysrv.TplDTO) (TplVO, error) {
			return TplVO{
				Id:        t.Id,
				Name:      t.Name,
				ApiKey:    t.ApiKey,
				NotifyCfg: t.NotifyCfg,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[TplVO]{
			DataResp: ginutil.DataResp[[]TplVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func sendNotification(c *gin.Context) {
	req := make(map[string]string)
	if util.ShouldBindJSON(&req, c) {
		err := notifysrv.Outer.SendNotificationByApiKey(c, notifysrv.SendNotifyByApiKeyReqDTO{
			ApiKey: c.Param("apiKey"),
			Params: req,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
