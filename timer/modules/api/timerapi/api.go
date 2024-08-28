package timerapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/timer/modules/service/timersrv"
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
		group := e.Group("/api/timer", apisession.CheckLogin)
		{
			// 创建定时任务
			group.POST("/create", createTimer)
			// 定时任务列表
			group.GET("/list", listTimer)
			// 启动定时任务
			group.PUT("/enable/:timerId", enableTimer)
			// 关闭定时任务
			group.PUT("/disable/:timerId", disableTimer)
			// 删除定时任务
			group.DELETE("/delete/:timerId", deleteTimer)
			// 触发定时任务
			group.PUT("/trigger/:timerId", triggerTask)
			// 编辑定时任务
			group.POST("/update", updateTimer)
		}
		group = e.Group("/api/timerLog", apisession.CheckLogin)
		{
			// 执行日志列表
			group.GET("/list", listLog)
		}
	})
}

func createTimer(c *gin.Context) {
	var req CreateTimerReqVO
	if util.ShouldBindJSON(&req, c) {
		err := timersrv.CreateTimer(c, timersrv.CreateTimerReqDTO{
			Name:     req.Name,
			CronExp:  req.CronExp,
			TeamId:   req.TeamId,
			Task:     req.Task,
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

func enableTimer(c *gin.Context) {
	err := timersrv.EnableTimer(c, timersrv.EnableTimerReqDTO{
		Id:       cast.ToInt64(c.Param("timerId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func disableTimer(c *gin.Context) {
	err := timersrv.DisableTimer(c, timersrv.DisableTimerReqDTO{
		Id:       cast.ToInt64(c.Param("timerId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deleteTimer(c *gin.Context) {
	err := timersrv.DeleteTimer(c, timersrv.DeleteTimerReqDTO{
		Id:       cast.ToInt64(c.Param("timerId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listTimer(c *gin.Context) {
	var req ListTimerReqVO
	if util.ShouldBindQuery(&req, c) {
		tasks, total, err := timersrv.ListTimer(c, timersrv.ListTimerReqDTO{
			TeamId:   req.TeamId,
			Name:     req.Name,
			PageNum:  req.PageNum,
			Env:      req.Env,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(tasks, func(t timersrv.TimerDTO) TimerVO {
			return TimerVO{
				Id:        t.Id,
				Name:      t.Name,
				CronExp:   t.CronExp,
				Task:      t.Task,
				TeamId:    t.TeamId,
				IsEnabled: t.IsEnabled,
				Env:       t.Env,
				Creator:   t.Creator,
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[TimerVO]{
			DataResp: ginutil.DataResp[[]TimerVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func listLog(c *gin.Context) {
	var req ListLogReqVO
	if util.ShouldBindQuery(&req, c) {
		logs, total, err := timersrv.ListLog(c, timersrv.ListLogReqDTO{
			Id:       req.Id,
			PageNum:  req.PageNum,
			Month:    req.Month,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(logs, func(t timersrv.LogDTO) TaskLogVO {
			return TaskLogVO{
				Task:        t.Task,
				ErrLog:      t.ErrLog,
				TriggerType: t.TriggerType,
				TriggerBy:   t.TriggerBy,
				IsSuccess:   t.IsSuccess,
				Created:     t.Created.Format("2006-01-02 15:04"),
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[TaskLogVO]{
			DataResp: ginutil.DataResp[[]TaskLogVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func triggerTask(c *gin.Context) {
	err := timersrv.TriggerTask(c, timersrv.TriggerTaskReqDTO{
		Id:       cast.ToInt64(c.Param("timerId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func updateTimer(c *gin.Context) {
	var req UpdateTimerReqVO
	if util.ShouldBindJSON(&req, c) {
		err := timersrv.UpdateTimer(c, timersrv.UpdateTimerReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			CronExp:  req.CronExp,
			Task:     req.Task,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
