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
	appsrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/app", apisession.CheckLogin)
		{
			// 创建应用服务
			group.POST("/create", createApp)
			group.POST("/update", updateApp)
			group.POST("/delete", deleteApp)
			// 应用服务列表
			group.GET("/list/:teamId", listApp)
			// 所有应用服务列表 管理员权限
			group.GET("/listAllByAdmin/:teamId", listAllByAdmin)
			group.POST("/transferTeam", transferTeam)
		}
	})
}

func createApp(c *gin.Context) {
	var req CreateAppReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.Outer.CreateApp(c, appsrv.CreateAppReqDTO{
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
		err := appsrv.Outer.UpdateApp(c, appsrv.UpdateAppReqDTO{
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
	var req DeleteAppReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.Outer.DeleteApp(c, appsrv.DeleteAppReqDTO{
			AppId:    req.AppId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listApp(c *gin.Context) {
	apps, err := appsrv.Outer.ListApp(c, appsrv.ListAppReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(apps, func(t appsrv.AppDTO) (AppVO, error) {
		return AppVO{
			AppId: t.AppId,
			Name:  t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]AppVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listAllByAdmin(c *gin.Context) {
	apps, err := appsrv.Outer.ListAllAppByAdmin(c, appsrv.ListAppReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(apps, func(t appsrv.AppDTO) (AppVO, error) {
		return AppVO{
			AppId: t.AppId,
			Name:  t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]AppVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func transferTeam(c *gin.Context) {
	var req TransferTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.Outer.TransferTeam(c, appsrv.TransferTeamReqDTO{
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
