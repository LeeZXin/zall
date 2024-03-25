package appapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/appsrv"
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
		group := e.Group("/api/app", apisession.CheckLogin)
		{
			group.POST("/insert", insertApp)
			group.POST("/update", updateApp)
			group.POST("/delete", deleteApp)
			group.POST("/list", listApp)
			group.POST("/transferTeam", transferTeam)
		}
	})
}

func insertApp(c *gin.Context) {
	var req InsertAppReqVO
	if util.ShouldBindJSON(&req, c) {
		err := appsrv.Outer.InsertApp(c, appsrv.InsertAppReqDTO{
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
	var req ListAppReqVO
	if util.ShouldBindJSON(&req, c) {
		apps, err := appsrv.Outer.ListApp(c, appsrv.ListAppReqDTO{
			AppId:    req.AppId,
			TeamId:   req.TeamId,
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
