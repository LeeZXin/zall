package taskapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/timer/modules/service/tasksrv"
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
		group := e.Group("/api/timerTask", apisession.CheckLogin)
		{
			group.POST("/insert", insertTask)
			group.POST("/list", listTask)
			group.POST("/enable", enableTask)
			group.POST("/disable", disableTask)
			group.POST("/delete", deleteTask)
			group.POST("/listLog", listLog)
			group.POST("/trigger", trigger)
			group.POST("/update", update)
		}
	})
}

func insertTask(c *gin.Context) {
	var req InsertTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.InsertTask(c, tasksrv.InsertTaskReqDTO{
			Name:     req.Name,
			CronExp:  req.CronExp,
			TaskType: req.TaskType,
			HttpTask: req.HttpTask,
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

func enableTask(c *gin.Context) {
	var req EnabledTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.EnableTask(c, tasksrv.EnableTaskReqDTO{
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

func disableTask(c *gin.Context) {
	var req DisableTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.DisableTask(c, tasksrv.DisableTaskReqDTO{
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

func deleteTask(c *gin.Context) {
	var req EnabledTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.DeleteTask(c, tasksrv.DeleteTaskReqDTO{
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

func listTask(c *gin.Context) {
	var req ListTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		tasks, cursor, err := tasksrv.Outer.ListTask(c, tasksrv.ListTaskReqDTO{
			TeamId:   req.TeamId,
			Name:     req.Name,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		resp := ListTaskRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   cursor,
		}
		resp.Data, _ = listutil.Map(tasks, func(t tasksrv.TaskDTO) (TaskVO, error) {
			return TaskVO{
				Id:         t.Id,
				Name:       t.Name,
				CronExp:    t.CronExp,
				TaskType:   t.TaskType,
				HttpTask:   t.HttpTask,
				TeamId:     t.TeamId,
				NextTime:   time.UnixMilli(t.NextTime).Format(timeutil.DefaultTimeFormat),
				TaskStatus: t.TaskStatus.Readable(),
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func listLog(c *gin.Context) {
	var req ListLogReqVO
	if util.ShouldBindJSON(&req, c) {
		logs, cursor, err := tasksrv.Outer.ListTaskLog(c, tasksrv.ListTaskLogReqDTO{
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
		resp.Data, _ = listutil.Map(logs, func(t tasksrv.TaskLogDTO) (TaskLogVO, error) {
			return TaskLogVO{
				TaskType:    t.TaskType,
				HttpTask:    t.HttpTask,
				LogContent:  t.LogContent,
				TriggerType: t.TriggerType.Readable(),
				TriggerBy:   t.TriggerBy,
				TaskStatus:  t.TaskStatus.Readable(),
				Created:     t.Created.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, resp)
	}
}

func trigger(c *gin.Context) {
	var req TriggerTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.TriggerTask(c, tasksrv.TriggerTaskReqDTO{
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

func update(c *gin.Context) {
	var req UpdateTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.UpdateTask(c, tasksrv.UpdateTaskReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			CronExp:  req.CronExp,
			TaskType: req.TaskType,
			HttpTask: req.HttpTask,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
