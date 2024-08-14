package promapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/service/promsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func InitApi() {
	promsrv.Init()
	cfgsrv.InitInner()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/promScrape", apisession.CheckLogin)
		{
			// 创建抓取配置
			group.POST("/create", createScrape)
			// 编辑抓取配置
			group.POST("/update", updateScrape)
			// 展示抓取配置列表
			group.GET("/list", listScrape)
			// 删除抓取配置
			group.DELETE("/delete/:scrapeId", deleteScrape)
		}
	})
}

func createScrape(c *gin.Context) {
	var req CreateScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := promsrv.Outer.CreateScrape(c, promsrv.CreateScrapeReqDTO{
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

func deleteScrape(c *gin.Context) {
	err := promsrv.Outer.DeleteScrape(c, promsrv.DeleteScrapeReqDTO{
		ScrapeId: cast.ToInt64(c.Param("scrapeId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listScrape(c *gin.Context) {
	scraps, err := promsrv.Outer.ListScrape(c, promsrv.ListScrapeReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(scraps, func(t promsrv.ScrapeDTO) (ScrapeVO, error) {
		return ScrapeVO{
			Id:         t.Id,
			Endpoint:   t.Endpoint,
			AppId:      t.AppId,
			Target:     t.Target,
			TargetType: t.TargetType,
			Created:    t.Created.Format(time.DateTime),
			Env:        t.Env,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ScrapeVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func updateScrape(c *gin.Context) {
	var req UpdateScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := promsrv.Outer.UpdateScrape(c, promsrv.UpdateScrapeReqDTO{
			ScrapeId:   req.ScrapeId,
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
