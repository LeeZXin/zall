package appapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/appsrv"
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
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/app", apisession.CheckLogin)
		{
			// 创建应用服务
			group.POST("/create", createApp)
			// 编辑应用服务
			group.POST("/update", updateApp)
			// 获取应用服务
			group.GET("/get/:appId", getApp)
			// 删除应用服务
			group.DELETE("/delete/:appId", deleteApp)
			// 应用服务列表
			group.GET("/list/:teamId", listApp)
			// 所有应用服务列表 管理员权限
			group.GET("/listAllByAdmin/:teamId", listAllByAdmin)
			// 所有应用服务列表 超级管理员权限
			group.GET("/listAllBySa", listAllBySa)
			// 迁移团队
			group.PUT("/transferTeam", transferTeam)
		}
	})
}

func createApp(c *gin.Context) {
	var req CreateAppReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.CreateApp(c, appsrv.CreateAppReqDTO{
			AppId:    req.AppId,
			TeamId:   req.TeamId,
			Name:     req.Name,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateApp(c *gin.Context) {
	var req UpdateAppReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.UpdateApp(c, appsrv.UpdateAppReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteApp(c *gin.Context) {
	err := appsrv.DeleteApp(c, appsrv.DeleteAppReqDTO{
		AppId:    c.Param("appId"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func getApp(c *gin.Context) {
	app, err := appsrv.GetAppWithPerm(c, appsrv.GetAppWithPermReqDTO{
		AppId:    c.Param("appId"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[AppWithPermVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: AppWithPermVO{
			AppVO: AppVO{
				AppId: app.AppId,
				Name:  app.Name,
			},
			Perm: app.Perm,
		},
	})
}

func listApp(c *gin.Context) {
	apps, err := appsrv.ListApp(c, appsrv.ListAppReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(apps, func(t appsrv.AppDTO) AppVO {
		return AppVO{
			AppId: t.AppId,
			Name:  t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]AppVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listAllByAdmin(c *gin.Context) {
	apps, err := appsrv.ListAllAppByAdmin(c, appsrv.ListAllAppByAdminReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(apps, func(t appsrv.AppDTO) AppVO {
		return AppVO{
			AppId: t.AppId,
			Name:  t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]AppVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listAllBySa(c *gin.Context) {
	apps, err := appsrv.ListAllAppBySa(c, appsrv.ListAllAppBySaReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(apps, func(t appsrv.AppDTO) AppVO {
		return AppVO{
			AppId: t.AppId,
			Name:  t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]AppVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func transferTeam(c *gin.Context) {
	var req TransferTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.TransferTeam(c, appsrv.TransferTeamReqDTO{
			AppId:    req.AppId,
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
