package alertapi

import (
	"github.com/LeeZXin/zall/alert/modules/service/alertsrv"
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
		group := e.Group("/api/alertConfig", apisession.CheckLogin)
		{
			// 新增配置
			group.POST("/create", createConfig)
			// 编辑配置
			group.POST("/update", updateConfig)
			// 删除配置
			group.DELETE("/delete/:configId", deleteConfig)
			// 启用配置
			group.PUT("/enable/:configId", enableConfig)
			// 关闭配置
			group.PUT("/disable/:configId", disableConfig)
			// 配置列表
			group.GET("/list", listConfig)
		}
	})
}

func createConfig(c *gin.Context) {
	var req CreateConfigReqVO
	if util.ShouldBindJSON(&req, c) {
		err := alertsrv.CreateConfig(c, alertsrv.CreateConfigReqDTO{
			Name:        req.Name,
			Alert:       req.Alert,
			AppId:       req.AppId,
			IntervalSec: req.IntervalSec,
			Env:         req.Env,
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
		err := alertsrv.UpdateConfig(c, alertsrv.UpdateConfigReqDTO{
			Id:          req.Id,
			Name:        req.Name,
			Alert:       req.Alert,
			IntervalSec: req.IntervalSec,
			Operator:    apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func enableConfig(c *gin.Context) {
	err := alertsrv.EnableConfig(c, alertsrv.EnableConfigReqDTO{
		Id:       cast.ToInt64(c.Param("configId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func disableConfig(c *gin.Context) {
	err := alertsrv.DisableConfig(c, alertsrv.DisableConfigReqDTO{
		Id:       cast.ToInt64(c.Param("configId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deleteConfig(c *gin.Context) {
	err := alertsrv.DeleteConfig(c, alertsrv.DeleteConfigReqDTO{
		Id:       cast.ToInt64(c.Param("configId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listConfig(c *gin.Context) {
	var req ListConfigReqVO
	if util.ShouldBindQuery(&req, c) {
		configs, total, err := alertsrv.ListConfig(c, alertsrv.ListConfigReqDTO{
			PageNum:  req.PageNum,
			AppId:    req.AppId,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(configs, func(t alertsrv.ConfigDTO) ConfigVO {
			return ConfigVO{
				Id:          t.Id,
				Name:        t.Name,
				AppId:       t.AppId,
				Content:     t.Content,
				IntervalSec: t.IntervalSec,
				IsEnabled:   t.IsEnabled,
				Creator:     t.Creator,
				Env:         t.Env,
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[ConfigVO]{
			DataResp: ginutil.DataResp[[]ConfigVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}
