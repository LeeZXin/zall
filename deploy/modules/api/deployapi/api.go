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
			// 所有配置来源列表
			group.GET("/listAll/:env", listAllServiceSource)
			// 获取应用服务绑定的配置来源
			group.GET("/listBind", listBindServiceSource)
			// 绑定应用服务和配置来源
			group.POST("/bindApp", bindAppAndServiceSource)
		}
		group = e.Group("/api/service", apisession.CheckLogin)
		{
			// 展示服务状态
			group.GET("/listStatus/:bindId", listServiceStatus)
			// 获取操作列表
			group.GET("/listActions/:bindId", listServiceStatusActions)
			// 操作服务
			group.PUT("/doAction", doServiceStatusAction)

		}
	})
}

func doServiceStatusAction(c *gin.Context) {
	var req DoServiceStatusActionReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.DoStatusAction(c, deploysrv.DoStatusActionReqDTO{
			BindId:    req.BindId,
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
	services, err := deploysrv.ListServiceStatus(c, deploysrv.ListServiceStatusReqDTO{
		BindId:   cast.ToInt64(c.Param("bindId")),
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
	actions, err := deploysrv.ListStatusActions(c, deploysrv.ListStatusActionReqDTO{
		BindId:   cast.ToInt64(c.Param("bindId")),
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
	varsList, err := deploysrv.ListPipelineVars(c, deploysrv.ListPipelineVarsReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(varsList, func(t deploysrv.PipelineVarsWithoutContentDTO) PipelineVarsWithoutContentVO {
		return PipelineVarsWithoutContentVO{
			Id:    t.Id,
			Name:  t.Name,
			AppId: t.AppId,
			Env:   t.Env,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]PipelineVarsWithoutContentVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func createPipelineVars(c *gin.Context) {
	var req CreatePipelineVarsReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.CreatePipelineVars(c, deploysrv.CreatePipelineVarsReqDTO{
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
		err := deploysrv.UpdatePipelineVars(c, deploysrv.UpdatePipelineVarsReqDTO{
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
	vars, err := deploysrv.GetPipelineVarsContent(c, deploysrv.GetPipelineVarsContentReqDTO{
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
	err := deploysrv.DeletePipelineVars(c, deploysrv.DeletePipelineVarsReqDTO{
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
	sources, err := deploysrv.ListServiceSource(c, deploysrv.ListServiceSourceReqDTO{
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(sources, func(t deploysrv.ServiceSourceDTO) ServiceSourceVO {
		return ServiceSourceVO{
			Id:      t.Id,
			Name:    t.Name,
			Env:     t.Env,
			Host:    t.Host,
			ApiKey:  t.ApiKey,
			Created: t.Created.Format(time.DateTime),
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]ServiceSourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listAllServiceSource(c *gin.Context) {
	sources, err := deploysrv.ListAllServiceSource(c, deploysrv.ListAllServiceSourceReqDTO{
		Env:      c.Param("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(sources, func(t deploysrv.SimpleServiceSourceDTO) SimpleServiceSourceVO {
		return SimpleServiceSourceVO{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleServiceSourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func createServiceSource(c *gin.Context) {
	var req CreateServiceSourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.CreateServiceSource(c, deploysrv.CreateServiceSourceReqDTO{
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
		err := deploysrv.UpdateServiceSource(c, deploysrv.UpdateServiceSourceReqDTO{
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
	err := deploysrv.DeleteServiceSource(c, deploysrv.DeleteServiceSourceReqDTO{
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
	err := deploysrv.RedoAgentStage(c, deploysrv.RedoAgentStageReqDTO{
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
		err := deploysrv.ConfirmInteractStage(c, deploysrv.ConfirmInteractStageReqDTO{
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
		err := deploysrv.ForceRedoNotSuccessfulAgentStages(c, deploysrv.ForceRedoNotSuccessfulAgentStagesReqDTO{
			PlanId:     req.PlanId,
			StageIndex: req.StageIndex,
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
	err := deploysrv.KillStage(c, deploysrv.KillStageReqDTO{
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
		err := deploysrv.CreatePlan(c, deploysrv.CreatePlanReqDTO{
			Name:            req.Name,
			PipelineId:      req.PipelineId,
			ArtifactVersion: req.ArtifactVersion,
			Operator:        apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func closePlan(c *gin.Context) {
	err := deploysrv.ClosePlan(c, deploysrv.ClosePlanReqDTO{
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
	err := deploysrv.StartPlan(c, deploysrv.StartPlanReqDTO{
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
		plans, total, err := deploysrv.ListPlan(c, deploysrv.ListPlanReqDTO{
			AppId:    req.AppId,
			Env:      req.Env,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(plans, func(t deploysrv.PlanDTO) PlanVO {
			return PlanVO{
				Id:              t.Id,
				PipelineId:      t.PipelineId,
				PipelineName:    t.PipelineName,
				Name:            t.Name,
				ArtifactVersion: t.ArtifactVersion,
				PlanStatus:      t.PlanStatus,
				Env:             t.Env,
				Creator:         t.Creator,
				Created:         t.Created.Format(time.DateTime),
			}
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
	plan, err := deploysrv.GetPlanDetail(c, deploysrv.GetPlanDetailReqDTO{
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
			Id:              plan.Id,
			PipelineId:      plan.PipelineId,
			PipelineName:    plan.PipelineName,
			PipelineConfig:  plan.PipelineConfig,
			Name:            plan.Name,
			ArtifactVersion: plan.ArtifactVersion,
			PlanStatus:      plan.PlanStatus,
			Env:             plan.Env,
			Creator:         plan.Creator,
			Created:         plan.Created.Format(time.DateTime),
		},
	})
}

func listStages(c *gin.Context) {
	stages, err := deploysrv.ListStages(c, deploysrv.ListStagesReqDTO{
		PlanId:   cast.ToInt64(c.Param("planId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(stages, func(t deploysrv.StageDTO) StageVO {
		subStages := listutil.MapNe(t.SubStages, func(t deploysrv.SubStageDTO) SubStageVO {
			return SubStageVO{
				Id:          t.Id,
				Agent:       t.Agent,
				AgentHost:   t.AgentHost,
				StageStatus: t.StageStatus,
				ExecuteLog:  t.ExecuteLog,
			}
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
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]StageVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func createPipeline(c *gin.Context) {
	var req CreatePipelineReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.CreatePipeline(c, deploysrv.CreatePipelineReqDTO{
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
		err := deploysrv.UpdatePipeline(c, deploysrv.UpdatePipelineReqDTO{
			Id:       req.PipelineId,
			Name:     req.Name,
			Config:   req.Config,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deletePipeline(c *gin.Context) {
	err := deploysrv.DeletePipeline(c, deploysrv.DeletePipelineReqDTO{
		Id:       cast.ToInt64(c.Param("pipelineId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listPipeline(c *gin.Context) {
	pipelines, err := deploysrv.ListPipeline(c, deploysrv.ListPipelineReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(pipelines, func(t deploysrv.PipelineDTO) PipelineVO {
		return PipelineVO{
			Id:     t.Id,
			AppId:  t.AppId,
			Config: t.Config,
			Env:    t.Env,
			Name:   t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]PipelineVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listPipelineWhenCreatePlan(c *gin.Context) {
	pipelines, err := deploysrv.ListPipelineWhenCreatePlan(c, deploysrv.ListPipelineWhenCreatePlanReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(pipelines, func(t deploysrv.SimplePipelineDTO) SimplePipelineVO {
		return SimplePipelineVO{
			Id:   t.Id,
			Env:  t.Env,
			Name: t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimplePipelineVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listBindServiceSource(c *gin.Context) {
	sources, err := deploysrv.ListBindServiceSource(c, deploysrv.ListBindServiceSourceReqDTO{
		AppId:    c.Query("appId"),
		Env:      c.Query("env"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(sources, func(t deploysrv.SimpleBindServiceSourceDTO) SimpleBindServiceSourceVO {
		return SimpleBindServiceSourceVO{
			Id:     t.Id,
			Name:   t.Name,
			BindId: t.BindId,
			Env:    t.Env,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleBindServiceSourceVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func bindAppAndServiceSource(c *gin.Context) {
	var req BindAppAndServiceSourceReqVO
	if util.ShouldBindJSON(&req, c) {
		err := deploysrv.BindAppAndServiceSource(c, deploysrv.BindAppAndServiceSourceReqDTO{
			AppId:        req.AppId,
			SourceIdList: req.SourceIdList,
			Env:          req.Env,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
