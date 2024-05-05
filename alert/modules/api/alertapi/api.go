package alertapi

import (
	"github.com/LeeZXin/zall/alert/modules/service/alertsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/alertConfig", apisession.CheckLogin)
		{
			group.POST("/insert", insertConfig)
			group.POST("/update", updateConfig)
			group.POST("/delete", deleteConfig)
			group.POST("/list", listConfig)
		}
	})
}

type demoReq struct {
	Content string `json:"content"`
}

func demo(c *gin.Context) {
	var req demoReq
	if util.ShouldBindJSON(&req, c) {
		logger.Logger.WithContext(c).Error(c.Request.Header, req.Content)
		c.String(http.StatusOK, "")
	}
}

func insertConfig(c *gin.Context) {
	var req InsertConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := alertsrv.Outer.InsertConfig(c, alertsrv.InsertConfigReqDTO{
			Name:        req.Name,
			Alert:       req.Alert,
			AppId:       req.AppId,
			IntervalSec: req.IntervalSec,
			SilenceSec:  req.SilenceSec,
			Enabled:     req.Enabled,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateConfig(c *gin.Context) {
	var req UpdateConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := alertsrv.Outer.UpdateConfig(c, alertsrv.UpdateConfigReqDTO{
			Id:          req.Id,
			Name:        req.Name,
			Alert:       req.Alert,
			IntervalSec: req.IntervalSec,
			SilenceSec:  req.SilenceSec,
			Enabled:     req.Enabled,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteConfig(c *gin.Context) {
	var req DeleteConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := alertsrv.Outer.DeleteConfig(c, alertsrv.DeleteConfigReqDTO{
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

func listConfig(c *gin.Context) {
	var req ListConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		configs, next, err := alertsrv.Outer.ListConfig(c, alertsrv.ListConfigReqDTO{
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			AppId:    req.AppId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(configs, func(t alertsrv.ConfigDTO) (ConfigVO, error) {
			return ConfigVO{
				Id:          t.Id,
				Name:        t.Name,
				AppId:       t.AppId,
				Content:     t.Content,
				IntervalSec: t.IntervalSec,
				SilenceSec:  t.SilenceSec,
				Enabled:     t.Enabled,
				Created:     t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[ConfigVO]{
			DataResp: ginutil.DataResp[[]ConfigVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}
