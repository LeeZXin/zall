package deployapi

import (
	"github.com/LeeZXin/zall/deploy/modules/service/deploysrv"
	"github.com/LeeZXin/zall/pkg/apisession"
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
	})
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
			Hosts:   t.Hosts,
			ApiKey:  t.ApiKey,
			Created: t.Created.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceSourceVO]{
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
			Hosts:    req.Hosts,
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
			Hosts:    req.Hosts,
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
