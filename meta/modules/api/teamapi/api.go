package teamapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/timeutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		// 项目
		group := e.Group("/api/team", apisession.CheckLogin)
		{
			group.POST("/insert", insertTeam)
			group.GET("/list", listTeam)
			group.POST("/delete", deleteTeam)
			group.POST("/update", updateTeam)
		}
		// 项目用户
		group = e.Group("/api/teamUser", apisession.CheckLogin)
		{
			group.POST("/upsert", upsertTeamUser)
			group.POST("/list", listTeamUser)
			group.POST("/delete", deleteTeamUser)
		}
		// 项目用户组
		group = e.Group("/api/teamUserGroup", apisession.CheckLogin)
		{
			group.POST("/insert", insertTeamUserGroup)
			group.POST("/list", listTeamUserGroup)
			group.POST("/updateName", updateTeamUserGroupName)
			group.POST("/updatePerm", updateTeamUserGroupPerm)
			group.POST("/delete", deleteTeamUserGroup)
		}
	})
}

func insertTeam(c *gin.Context) {
	var req InsertTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.InsertTeam(c, teamsrv.InsertTeamReqDTO{
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

func listTeam(c *gin.Context) {
	teamList, err := teamsrv.Outer.ListTeam(c, teamsrv.ListTeamReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	ret := ListTeamRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
	}
	ret.Data, _ = listutil.Map(teamList, func(t teamsrv.TeamDTO) (TeamVO, error) {
		return TeamVO{
			TeamId:  t.TeamId,
			Name:    t.Name,
			Created: t.Created.Format(timeutil.DefaultTimeFormat),
		}, nil
	})
	c.JSON(http.StatusOK, ret)
}

func deleteTeam(c *gin.Context) {
	var req DeleteTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.DeleteTeam(c, teamsrv.DeleteTeamReqDTO{
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

func updateTeam(c *gin.Context) {
	var req UpdateTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.UpdateTeam(c, teamsrv.UpdateTeamReqDTO{
			TeamId:   req.TeamId,
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

func upsertTeamUser(c *gin.Context) {
	var req UpsertTeamUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.UpsertTeamUser(c, teamsrv.UpsertTeamUserReqDTO{
			TeamId:   req.TeamId,
			Account:  req.Account,
			GroupId:  req.GroupId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listTeamUser(c *gin.Context) {
	var req ListTeamUserReqVO
	if util.ShouldBindJSON(&req, c) {
		users, next, err := teamsrv.Outer.ListTeamUser(c, teamsrv.ListTeamUserReqDTO{
			TeamId:   req.TeamId,
			Account:  req.Account,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListTeamUserRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Next:     next,
		}
		ret.Data, _ = listutil.Map(users, func(t teamsrv.TeamUserDTO) (TeamUserVO, error) {
			return TeamUserVO{
				TeamId:    t.TeamId,
				Account:   t.Account,
				GroupId:   t.GroupId,
				GroupName: t.GroupName,
				Created:   t.Created.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func deleteTeamUser(c *gin.Context) {
	var req DeleteTeamUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.DeleteTeamUser(c, teamsrv.DeleteTeamUserReqDTO{
			TeamId:   req.TeamId,
			Account:  req.Account,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func insertTeamUserGroup(c *gin.Context) {
	var req InsertTeamUserGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.InsertTeamUserGroup(c, teamsrv.InsertTeamUserGroupReqDTO{
			TeamId:   req.TeamId,
			Name:     req.Name,
			Perm:     req.Perm,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func updateTeamUserGroupName(c *gin.Context) {
	var req UpdateTeamUserGroupNameReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.UpdateTeamUserGroupName(c, teamsrv.UpdateTeamUserGroupNameReqDTO{
			GroupId:  req.GroupId,
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

func updateTeamUserGroupPerm(c *gin.Context) {
	var req UpdateTeamUserGroupPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.UpdateTeamUserGroupPerm(c, teamsrv.UpdateTeamUserGroupPermReqDTO{
			GroupId:  req.GroupId,
			Perm:     req.Perm,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteTeamUserGroup(c *gin.Context) {
	var req DeleteTeamUserGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.DeleteTeamUserGroup(c, teamsrv.DeleteTeamUserGroupReqDTO{
			GroupId:  req.GroupId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listTeamUserGroup(c *gin.Context) {
	var req ListTeamUserGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		groups, err := teamsrv.Outer.ListTeamUserGroup(c, teamsrv.ListTeamUserGroupReqDTO{
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret, _ := listutil.Map(groups, func(t teamsrv.TeamUserGroupDTO) (TeamUserGroupVO, error) {
			return TeamUserGroupVO{
				GroupId: t.GroupId,
				TeamId:  t.TeamId,
				Name:    t.Name,
				Perm:    t.Perm,
			}, nil
		})
		c.JSON(http.StatusOK, ListTeamUserGroupRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}
