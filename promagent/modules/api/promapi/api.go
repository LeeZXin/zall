package promapi

import (
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
		err := promsrv.CreateScrape(c, promsrv.CreateScrapeReqDTO{
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
	err := promsrv.DeleteScrape(c, promsrv.DeleteScrapeReqDTO{
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
	var req ListScrapeReqVO
	if util.ShouldBindQuery(&req, c) {
		scraps, total, err := promsrv.ListScrape(c, promsrv.ListScrapeReqDTO{
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
		data, _ := listutil.Map(scraps, func(t promsrv.ScrapeDTO) (ScrapeVO, error) {
			return ScrapeVO{
				Id:         t.Id,
				Endpoint:   t.Endpoint,
				AppId:      t.AppId,
				AppName:    t.AppName,
				Target:     t.Target,
				TargetType: t.TargetType,
				Created:    t.Created.Format(time.DateTime),
				Env:        t.Env,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ScrapeVO]{
			DataResp: ginutil.DataResp[[]ScrapeVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}

}

func updateScrape(c *gin.Context) {
	var req UpdateScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := promsrv.UpdateScrape(c, promsrv.UpdateScrapeReqDTO{
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
