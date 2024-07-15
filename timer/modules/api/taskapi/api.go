package taskapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/timer/modules/service/tasksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func InitApi() {
	tasksrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/timerTask", apisession.CheckLogin)
		{
			// 创建定时任务
			group.POST("/create", createTask)
			// 定时任务列表
			group.GET("/list", listTask)
			// 启动定时任务
			group.PUT("/enable/:taskId", enableTask)
			// 关闭定时任务
			group.PUT("/disable/:taskId", disableTask)
			// 删除定时任务
			group.DELETE("/delete/:taskId", deleteTask)
			// 触发定时任务
			group.PUT("/trigger/:taskId", triggerTask)
			// 编辑定时任务
			group.POST("/update", updateTask)
		}
		group = e.Group("/api/timerLog", apisession.CheckLogin)
		{
			group.GET("/list", pageLog)
		}
	})
}

func createTask(c *gin.Context) {
	var req CreateTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.CreateTask(c, tasksrv.CreateTaskReqDTO{
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

func enableTask(c *gin.Context) {
	err := tasksrv.Outer.EnableTask(c, tasksrv.EnableTaskReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func disableTask(c *gin.Context) {
	err := tasksrv.Outer.DisableTask(c, tasksrv.DisableTaskReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deleteTask(c *gin.Context) {
	err := tasksrv.Outer.DeleteTask(c, tasksrv.DeleteTaskReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listTask(c *gin.Context) {
	var req ListTaskReqVO
	if util.ShouldBindQuery(&req, c) {
		tasks, total, err := tasksrv.Outer.ListTask(c, tasksrv.ListTaskReqDTO{
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
		data, _ := listutil.Map(tasks, func(t tasksrv.TaskDTO) (TaskVO, error) {
			return TaskVO{
				Id:        t.Id,
				Name:      t.Name,
				CronExp:   t.CronExp,
				Task:      t.Task,
				TeamId:    t.TeamId,
				IsEnabled: t.IsEnabled,
				Env:       t.Env,
				Creator:   t.Creator,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[TaskVO]{
			DataResp: ginutil.DataResp[[]TaskVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func pageLog(c *gin.Context) {
	var req PageLogReqVO
	if util.ShouldBindQuery(&req, c) {
		logs, total, err := tasksrv.Outer.PageTaskLog(c, tasksrv.PageTaskLogReqDTO{
			TaskId:   req.TaskId,
			PageNum:  req.PageNum,
			DateStr:  req.DateStr,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(logs, func(t tasksrv.TaskLogDTO) (TaskLogVO, error) {
			return TaskLogVO{
				Task:        t.Task,
				ErrLog:      t.ErrLog,
				TriggerType: t.TriggerType,
				TriggerBy:   t.TriggerBy,
				IsSuccess:   t.IsSuccess,
				Created:     t.Created.Format("2006-01-02 15:04"),
			}, nil
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
	err := tasksrv.Outer.TriggerTask(c, tasksrv.TriggerTaskReqDTO{
		TaskId:   cast.ToInt64(c.Param("taskId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func updateTask(c *gin.Context) {
	var req UpdateTaskReqVO
	if util.ShouldBindJSON(&req, c) {
		err := tasksrv.Outer.UpdateTask(c, tasksrv.UpdateTaskReqDTO{
			TaskId:   req.TaskId,
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
