package teamhookapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
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
		group := e.Group("/api/teamHook", apisession.CheckLogin)
		{
			// 创建 team hook
			group.POST("/create", createTeamHook)
			// 编辑 team hook
			group.POST("/update", updateTeamHook)
			// 删除 team hook
			group.DELETE("/delete/:hookId", deleteTeamHook)
			// team hook列表
			group.GET("/list/:teamId", listTeamHook)
		}
	})
}

func createTeamHook(c *gin.Context) {
	var req CreateTeamHookReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamhooksrv.CreateTeamHook(c, teamhooksrv.CreateTeamHookReqDTO{
			Name:     req.Name,
			TeamId:   req.TeamId,
			Events:   req.Events,
			HookType: req.HookType,
			HookCfg:  req.HookCfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateTeamHook(c *gin.Context) {
	var req UpdateTeamHookReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamhooksrv.UpdateTeamHook(c, teamhooksrv.UpdateTeamHookReqDTO{
			Id:       req.Id,
			Name:     req.Name,
			Events:   req.Events,
			HookType: req.HookType,
			HookCfg:  req.HookCfg,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteTeamHook(c *gin.Context) {
	err := teamhooksrv.DeleteTeamHook(c, teamhooksrv.DeleteTeamHookReqDTO{
		Id:       cast.ToInt64(c.Param("hookId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listTeamHook(c *gin.Context) {
	hooks, err := teamhooksrv.ListTeamHook(c, teamhooksrv.ListTeamHookReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(hooks, func(t teamhooksrv.TeamHookDTO) (TeamHookVO, error) {
		return TeamHookVO{
			Id:       t.Id,
			Name:     t.Name,
			TeamId:   t.TeamId,
			Events:   t.Events,
			HookType: t.HookType,
			HookCfg:  t.HookCfg,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]TeamHookVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}
