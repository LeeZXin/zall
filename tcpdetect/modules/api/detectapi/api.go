package detectapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/tcpdetect/modules/service/detectsrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/timeutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/tcpDetect", apisession.CheckLogin)
		{
			group.POST("insert", insertDetect)
			group.POST("delete", deleteDetect)
			group.POST("update", updateDetect)
			group.POST("list", listDetect)
			group.POST("listLog", listLog)
			group.POST("enable", enableDetect)
			group.POST("/disable", disableDetect)
		}
	})
}

func insertDetect(c *gin.Context) {
	var req InsertDetectReqVO
	if util.ShouldBindJSON(&req, c) {
		err := detectsrv.Outer.InsertDetect(c, detectsrv.InsertDetectReqDTO{
			Ip:       req.Ip,
			Port:     req.Port,
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

func updateDetect(c *gin.Context) {
	var req UpdateDetectReqVO
	if util.ShouldBindJSON(&req, c) {
		err := detectsrv.Outer.UpdateDetect(c, detectsrv.UpdateDetectReqDTO{
			Id:       req.Id,
			Ip:       req.Ip,
			Port:     req.Port,
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

func deleteDetect(c *gin.Context) {
	var req DeleteDetectReqVO
	if util.ShouldBindJSON(&req, c) {
		err := detectsrv.Outer.DeleteDetect(c, detectsrv.DeleteDetectReqDTO{
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

func listDetect(c *gin.Context) {
	var req ListDetectReqVO
	if util.ShouldBindJSON(&req, c) {
		detects, cursor, err := detectsrv.Outer.ListDetect(c, detectsrv.ListDetectReqDTO{
			Name:     req.Name,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListDetectRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   cursor,
		}
		validHeartbeatTime := time.Now().Add(-20 * time.Second).UnixMilli()
		resp.Data, _ = listutil.Map(detects, func(t detectsrv.DetectDTO) (DetectVO, error) {
			return DetectVO{
				Id:            t.Id,
				Ip:            t.Ip,
				Port:          t.Port,
				Name:          t.Name,
				HeartbeatTime: time.UnixMilli(t.HeartbeatTime).Format(timeutil.DefaultTimeFormat),
				Valid:         t.HeartbeatTime > validHeartbeatTime,
				Enabled:       t.Enabled,
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func listLog(c *gin.Context) {
	var req ListLogReqVO
	if util.ShouldBindJSON(&req, c) {
		logs, cursor, err := detectsrv.Outer.ListLog(c, detectsrv.ListLogReqDTO{
			Id:       req.Id,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListLogRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   cursor,
		}
		resp.Data, _ = listutil.Map(logs, func(t detectsrv.LogDTO) (LogVO, error) {
			return LogVO{
				Ip:      t.Ip,
				Port:    t.Port,
				Valid:   t.Valid,
				Created: t.Created.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func enableDetect(c *gin.Context) {
	var req EnableDetectReqVO
	if util.ShouldBindJSON(&req, c) {
		err := detectsrv.Outer.EnabledDetect(c, detectsrv.EnableDetectReqDTO{
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

func disableDetect(c *gin.Context) {
	var req DisableDetectReqVO
	if util.ShouldBindJSON(&req, c) {
		err := detectsrv.Outer.DisableDetect(c, detectsrv.DisableDetectReqDTO{
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
