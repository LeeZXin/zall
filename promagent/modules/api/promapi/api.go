package promapi

import (
	"context"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/service/promsrv"
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
		group := e.Group("/api/promScrapeBySa", apisession.CheckLogin)
		{
			// 创建抓取配置
			group.POST("/create", createScrapeBySa)
			// 编辑抓取配置
			group.POST("/update", updateScrapeBySa)
			// 展示抓取配置列表
			group.GET("/list", listScrapeBySa)
			// 删除抓取配置
			group.DELETE("/delete/:scrapeId", deleteScrapeBySa)
		}
		group = e.Group("/api/promScrapeByTeam", apisession.CheckLogin)
		{
			// 创建抓取配置
			group.POST("/create", createScrapeByTeam)
			// 编辑抓取配置
			group.POST("/update", updateScrapeByTeam)
			// 展示抓取配置列表
			group.GET("/list", listScrapeByTeam)
			// 删除抓取配置
			group.DELETE("/delete/:scrapeId", deleteScrapeByTeam)
		}
	})
}

func createScrapeBySa(c *gin.Context) {
	createScrape(c, promsrv.CreateScrapeBySa)
}

func createScrapeByTeam(c *gin.Context) {
	createScrape(c, promsrv.CreateScrapeByTeam)
}

func createScrape(c *gin.Context, f func(context.Context, promsrv.CreateScrapeReqDTO) error) {
	var req CreateScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := f(c, promsrv.CreateScrapeReqDTO{
			Endpoint:   req.Endpoint,
			AppId:      req.AppId,
			Target:     req.Target,
			TargetType: req.TargetType,
			Env:        req.Env,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteScrapeBySa(c *gin.Context) {
	deleteScrape(c, promsrv.DeleteScrapeBySa)
}

func deleteScrapeByTeam(c *gin.Context) {
	deleteScrape(c, promsrv.DeleteScrapeByTeam)
}

func deleteScrape(c *gin.Context, f func(context.Context, promsrv.DeleteScrapeReqDTO) error) {
	err := f(c, promsrv.DeleteScrapeReqDTO{
		Id:       cast.ToInt64(c.Param("scrapeId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listScrapeBySa(c *gin.Context) {
	var req ListScrapeReqVO
	if util.ShouldBindQuery(&req, c) {
		scrapes, total, err := promsrv.ListScrapeBySa(c, promsrv.ListScrapeReqDTO{
			Endpoint: req.Endpoint,
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(scrapes, func(t promsrv.ScrapeBySaDTO) (ScrapeBySaVO, error) {
			return ScrapeBySaVO{
				Id:         t.Id,
				Endpoint:   t.Endpoint,
				AppId:      t.AppId,
				AppName:    t.AppName,
				TeamId:     t.TeamId,
				TeamName:   t.TeamName,
				Target:     t.Target,
				TargetType: t.TargetType,
				Env:        t.Env,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ScrapeBySaVO]{
			DataResp: ginutil.DataResp[[]ScrapeBySaVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func listScrapeByTeam(c *gin.Context) {
	var req ListScrapeReqVO
	if util.ShouldBindQuery(&req, c) {
		scrapes, total, err := promsrv.ListScrapeByTeam(c, promsrv.ListScrapeReqDTO{
			Endpoint: req.Endpoint,
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(scrapes, func(t promsrv.ScrapeByTeamDTO) (ScrapeByTeamVO, error) {
			return ScrapeByTeamVO{
				Id:         t.Id,
				Endpoint:   t.Endpoint,
				Target:     t.Target,
				TargetType: t.TargetType,
				Env:        t.Env,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ScrapeByTeamVO]{
			DataResp: ginutil.DataResp[[]ScrapeByTeamVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func updateScrapeBySa(c *gin.Context) {
	updateScrape(c, promsrv.UpdateScrapeBySa)
}

func updateScrapeByTeam(c *gin.Context) {
	updateScrape(c, promsrv.UpdateScrapeByTeam)
}

func updateScrape(c *gin.Context, f func(context.Context, promsrv.UpdateScrapeReqDTO) error) {
	var req UpdateScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := f(c, promsrv.UpdateScrapeReqDTO{
			Id:         req.ScrapeId,
			Endpoint:   req.Endpoint,
			Target:     req.Target,
			TargetType: req.TargetType,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
