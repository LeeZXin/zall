package deployapi

import (
	"github.com/LeeZXin/zall/deploy/modules/service/deploysrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/status"
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
	deploysrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/deployPlan", apisession.CheckLogin)
		{
			// 创建发布计划
			group.POST("/create", createPlan)
			// 关闭发布计划
			group.PUT("/close/:planId", closePlan)
			// 展示发布计划列表
			group.GET("/list", listPlan)
			// 开始发布计划
			group.PUT("/start/:planId", startPlan)
			// 服务列表
			group.GET("/listPipeline", listPipelineWhenCreatePlan)
			// 计划详情
			group.GET("/detail/:planId", getPlanDetail)
			// 流水线详情
			group.GET("/listStages/:planId", listStages)
		}
		group = e.Group("/api/deployStage", apisession.CheckLogin)
		{
			// 确认交互阶段
			group.POST("/confirm", confirmInteractStage)
			// 重新执行某个agent
			group.PUT("/redoAgent/:stageId", redoAgentStage)
			// 强行重新执行某个stage
			group.POST("/forceRedoStage", forceRedoStage)
			// 中止某个阶段
			group.PUT("/kill/:planId/:index", killStage)
		}
		group = e.Group("/api/pipeline", apisession.CheckLogin)
		{
			// 创建服务
			group.POST("/create", createPipeline)
			// 编辑服务
			group.POST("/update", updatePipeline)
			// 删除服务
			group.DELETE("/delete/:pipelineId", deletePipeline)
			// 服务列表
			group.GET("/list", listPipeline)
		}
		group = e.Group("/api/pipelineVars", apisession.CheckLogin)
		{
			// 变量列表
			group.GET("/list", listPipelineVars)
			// 创建变量
			group.POST("/create", createPipelineVars)
			// 编辑变量
			group.POST("/update", updatePipelineVars)
			// 删除变量
			group.DELETE("/delete/:varsId", deletePipelineVars)
			// 变量内容
			group.GET("/content/:varsId", getPipelineVars)
		}
		group = e.Group("/api/serviceSource", apisession.CheckLogin)
		{
			// 创建服务来源
			group.POST("/create", createServiceSource)
			// 编辑服务来源
			group.POST("/update", updateServiceSource)
			// 删除服务来源
			group.DELETE("/delete/:sourceId", deleteServiceSource)
			// 服务来源列表
			group.GET("/list", listServiceSource)
		}
		group = e.Group("/api/service", apisession.CheckLogin)
		{
			// 展示服务来源
			group.GET("/listSource", listStatusSource)
			// 展示服务状态
			group.GET("/listStatus/:sourceId", listServiceStatus)
			// 获取操作列表
			group.GET("/listActions/:sourceId", listServiceStatusActions)
			// 操作服务
			group.PUT("/doAction", doServiceStatusAction)

		}
	})
}

func doServiceStatusAction(c *gin.Context) {
	var req DoServiceStatusActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.DoStatusAction(c, deploysrv.DoStatusActionReqDTO{
			SourceId:  req.SourceId,
			ServiceId: req.ServiceId,
			Action:    req.Action,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}

}

func listServiceStatus(c *gin.Context) {
	services, err := deploysrv.Outer.ListServiceStatus(c, deploysrv.ListServiceStatusReqDTO{
		SourceId: cast.ToInt64(c.Param("sourceId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]status.Service]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     services,
	})
}

func listServiceStatusActions(c *gin.Context) {
	actions, err := deploysrv.Outer.ListStatusActions(c, deploysrv.ListStatusActionReqDTO{
		SourceId: cast.ToInt64(c.Param("sourceId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     actions,
	})
}

func listPipelineVars(c *gin.Context) {
	varsList, err := deploysrv.Outer.ListPipelineVars(c, deploysrv.ListPipelineVarsReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(varsList, func(t deploysrv.PipelineVarsWithoutContentDTO) (PipelineVarsWithoutContentVO, error) {
		return PipelineVarsWithoutContentVO{
			Id:    t.Id,
			Name:  t.Name,
			AppId: t.AppId,
			Env:   t.Env,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]PipelineVarsWithoutContentVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func createPipelineVars(c *gin.Context) {
	var req CreatePipelineVarsReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreatePipelineVars(c, deploysrv.CreatePipelineVarsReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Name:     req.Name,
			Content:  req.Content,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updatePipelineVars(c *gin.Context) {
	var req UpdatePipelineVarsReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdatePipelineVars(c, deploysrv.UpdatePipelineVarsReqDTO{
			Id:       req.Id,
			Content:  req.Content,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getPipelineVars(c *gin.Context) {
	vars, err := deploysrv.Outer.GetPipelineVarsContent(c, deploysrv.GetPipelineVarsContentReqDTO{
		Id:       cast.ToInt64(c.Param("varsId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[PipelineVarsVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: PipelineVarsVO{
			Id:      vars.Id,
			Name:    vars.Name,
			AppId:   vars.AppId,
			Env:     vars.Env,
			Content: vars.Content,
		},
	})
}

func deletePipelineVars(c *gin.Context) {
	err := deploysrv.Outer.DeletePipelineVars(c, deploysrv.DeletePipelineVarsReqDTO{
		Id:       cast.ToInt64(c.Param("varsId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listServiceSource(c *gin.Context) {
	sources, err := deploysrv.Outer.ListServiceSource(c, deploysrv.ListServiceSourceReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(sources, func(t deploysrv.ServiceSourceDTO) (ServiceSourceVO, error) {
		return ServiceSourceVO{
			Id:      t.Id,
			Name:    t.Name,
			AppId:   t.AppId,
			Env:     t.Env,
			Host:    t.Host,
			ApiKey:  t.ApiKey,
			Created: t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceSourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listStatusSource(c *gin.Context) {
	sources, err := deploysrv.Outer.ListStatusSource(c, deploysrv.ListStatusSourceReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(sources, func(t deploysrv.StatusSourceDTO) (StatusSourceVO, error) {
		return StatusSourceVO{
			Id:   t.Id,
			Name: t.Name,
			Env:  t.Env,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]StatusSourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func createServiceSource(c *gin.Context) {
	var req CreateServiceSourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreateServiceSource(c, deploysrv.CreateServiceSourceReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			Name:     req.Name,
			Host:     req.Host,
			ApiKey:   req.ApiKey,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateServiceSource(c *gin.Context) {
	var req UpdateServiceSourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdateServiceSource(c, deploysrv.UpdateServiceSourceReqDTO{
			SourceId: req.SourceId,
			Name:     req.Name,
			Host:     req.Host,
			ApiKey:   req.ApiKey,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteServiceSource(c *gin.Context) {
	err := deploysrv.Outer.DeleteServiceSource(c, deploysrv.DeleteServiceSourceReqDTO{
		SourceId: cast.ToInt64(c.Param("sourceId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func redoAgentStage(c *gin.Context) {
	err := deploysrv.Outer.RedoAgentStage(c, deploysrv.RedoAgentStageReqDTO{
		StageId:  cast.ToInt64(c.Param("stageId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func confirmInteractStage(c *gin.Context) {
	var req ConfirmInteractStageReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.ConfirmInteractStage(c, deploysrv.ConfirmInteractStageReqDTO{
			PlanId:     req.PlanId,
			StageIndex: req.StageIndex,
			Args:       req.Args,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func forceRedoStage(c *gin.Context) {
	var req ForceRedoStageReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.ForceRedoNotSuccessfulAgentStages(c, deploysrv.ForceRedoNotSuccessfulAgentStagesReqDTO{
			PlanId:     req.PlanId,
			StageIndex: req.StageIndex,
			Args:       req.Args,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func killStage(c *gin.Context) {
	err := deploysrv.Outer.KillStage(c, deploysrv.KillStageReqDTO{
		PlanId:     cast.ToInt64(c.Param("planId")),
		StageIndex: cast.ToInt(c.Param("index")),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func createPlan(c *gin.Context) {
	var req CreatePlanReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreatePlan(c, deploysrv.CreatePlanReqDTO{
			Name:           req.Name,
			PipelineId:     req.PipelineId,
			ProductVersion: req.ProductVersion,
			Operator:       apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func closePlan(c *gin.Context) {
	err := deploysrv.Outer.ClosePlan(c, deploysrv.ClosePlanReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func startPlan(c *gin.Context) {
	err := deploysrv.Outer.StartPlan(c, deploysrv.StartPlanReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listPlan(c *gin.Context) {
	var req ListPlanReqVO
	if util.ShouldBindQuery(&req, c) {
		plans, total, err := deploysrv.Outer.ListPlan(c, deploysrv.ListPlanReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(plans, func(t deploysrv.PlanDTO) (PlanVO, error) {
			return PlanVO{
				Id:             t.Id,
				PipelineId:     t.PipelineId,
				PipelineName:   t.PipelineName,
				Name:           t.Name,
				ProductVersion: t.ProductVersion,
				PlanStatus:     t.PlanStatus,
				Env:            t.Env,
				Creator:        t.Creator,
				Created:        t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[PlanVO]{
			DataResp: ginutil.DataResp[[]PlanVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func getPlanDetail(c *gin.Context) {
	plan, err := deploysrv.Outer.GetPlanDetail(c, deploysrv.GetPlanDetailReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[PlanDetailVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: PlanDetailVO{
			Id:             plan.Id,
			PipelineId:     plan.PipelineId,
			PipelineName:   plan.PipelineName,
			PipelineConfig: plan.PipelineConfig,
			Name:           plan.Name,
			ProductVersion: plan.ProductVersion,
			PlanStatus:     plan.PlanStatus,
			Env:            plan.Env,
			Creator:        plan.Creator,
			Created:        plan.Created.Format(time.DateTime),
		},
	})
}

func listStages(c *gin.Context) {
	stages, err := deploysrv.Outer.ListStages(c, deploysrv.ListStagesReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(stages, func(t deploysrv.StageDTO) (StageVO, error) {
		subStages, _ := listutil.Map(t.SubStages, func(t deploysrv.SubStageDTO) (SubStageVO, error) {
			return SubStageVO{
				Id:          t.Id,
				Agent:       t.Agent,
				AgentHost:   t.AgentHost,
				StageStatus: t.StageStatus,
				ExecuteLog:  t.ExecuteLog,
			}, nil
		})
		return StageVO{
			Name:                             t.Name,
			Percent:                          t.Percent,
			Total:                            t.Total,
			Done:                             t.Done,
			IsAutomatic:                      t.IsAutomatic,
			HasError:                         t.HasError,
			IsRunning:                        t.IsRunning,
			IsAllDone:                        t.IsAllDone,
			WaitInteract:                     t.WaitInteract,
			SubStages:                        subStages,
			Script:                           t.Script,
			Confirm:                          t.Confirm,
			CanForceRedoUnSuccessAgentStages: t.CanForceRedoUnSuccessAgentStages,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]StageVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func createPipeline(c *gin.Context) {
	var req CreatePipelineReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.CreatePipeline(c, deploysrv.CreatePipelineReqDTO{
			AppId:    req.AppId,
			Name:     req.Name,
			Config:   req.Config,
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

func updatePipeline(c *gin.Context) {
	var req UpdatePipelineReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.Outer.UpdatePipeline(c, deploysrv.UpdatePipelineReqDTO{
			PipelineId: req.PipelineId,
			Name:       req.Name,
			Config:     req.Config,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deletePipeline(c *gin.Context) {
	err := deploysrv.Outer.DeletePipeline(c, deploysrv.DeletePipelineReqDTO{
		PipelineId: cast.ToInt64(c.Param("pipelineId")),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listPipeline(c *gin.Context) {
	pipelines, err := deploysrv.Outer.ListPipeline(c, deploysrv.ListPipelineReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(pipelines, func(t deploysrv.PipelineDTO) (PipelineVO, error) {
		return PipelineVO{
			Id:     t.Id,
			AppId:  t.AppId,
			Config: t.Config,
			Env:    t.Env,
			Name:   t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]PipelineVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listPipelineWhenCreatePlan(c *gin.Context) {
	pipelines, err := deploysrv.Outer.ListPipelineWhenCreatePlan(c, deploysrv.ListPipelineWhenCreatePlanReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(pipelines, func(t deploysrv.SimplePipelineDTO) (SimplePipelineVO, error) {
		return SimplePipelineVO{
			Id:   t.Id,
			Env:  t.Env,
			Name: t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimplePipelineVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
