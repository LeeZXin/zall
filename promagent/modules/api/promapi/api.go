package promapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/service/promsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/promScrape", apisession.CheckLogin)
		{
			group.POST("/insert", insertScrape)
			group.POST("/update", updateScrape)
			group.POST("/list", listScrape)
			group.POST("/delete", deleteScrape)
		}
	})
}

func insertScrape(c *gin.Context) {
	var req InsertScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := promsrv.Outer.InsertScrape(c, promsrv.InsertScrapeReqDTO{
			ServerUrl:  req.ServerUrl,
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
	var req DeleteScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := promsrv.Outer.DeleteScrape(c, promsrv.DeleteScrapeReqDTO{
			Id:       req.Id,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listScrape(c *gin.Context) {
	var req ListScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		scraps, err := promsrv.Outer.ListScrape(c, promsrv.ListScrapeReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(scraps, func(t promsrv.ScrapeDTO) (ScrapeVO, error) {
			return ScrapeVO{
				Id:         t.Id,
				ServerUrl:  t.ServerUrl,
				AppId:      t.AppId,
				Target:     t.Target,
				TargetType: t.TargetType.Readable(),
				Created:    t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]ScrapeVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func updateScrape(c *gin.Context) {
	var req UpdateScrapeReqVO
	if util.ShouldBindJSON(&req, c) {
		err := promsrv.Outer.UpdateScrape(c, promsrv.UpdateScrapeReqDTO{
			Id:         req.Id,
			ServerUrl:  req.ServerUrl,
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
