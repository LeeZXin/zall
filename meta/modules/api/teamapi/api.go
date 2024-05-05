package teamapi

import (
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/perm"
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
		// 项目
		group := e.Group("/api/team", apisession.CheckLogin)
		{
			// 创建团队
			group.POST("/create", createTeam)
			// 获取所在团队列表
			group.GET("/list", listTeam)
			group.POST("/delete", deleteTeam)
			group.POST("/update", updateTeam)
			// 是否是团队管理员
			group.GET("/isAdmin/:teamId", isAdmin)
			// 获取团队权限
			group.GET("/getTeamPerm/:teamId", getTeamPerm)
			// 获取团队信息
			group.GET("/get/:teamId", getTeam)
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

func isAdmin(c *gin.Context) {
	teamId := cast.ToInt64(c.Param("teamId"))
	b, err := teamsrv.Outer.IsAdmin(c, teamsrv.IsAdminReqDTO{
		TeamId:   teamId,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[bool]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     b,
	})
}

func getTeamPerm(c *gin.Context) {
	teamId := cast.ToInt64(c.Param("teamId"))
	teamPerm, err := teamsrv.Outer.GetTeamPerm(c, teamsrv.GetTeamPermReqDTO{
		TeamId:   teamId,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[perm.TeamPerm]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     teamPerm,
	})
}

func getTeam(c *gin.Context) {
	teamId := cast.ToInt64(c.Param("teamId"))
	team, err := teamsrv.Outer.GetTeam(c, teamsrv.GetTeamReqDTO{
		TeamId:   teamId,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[TeamVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: TeamVO{
			TeamId: team.Id,
			Name:   team.Name,
		},
	})
}

func createTeam(c *gin.Context) {
	var req CreateTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.CreateTeam(c, teamsrv.CreateTeamReqDTO{
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
	data, _ := listutil.Map(teamList, func(t teammd.Team) (TeamVO, error) {
		return TeamVO{
			TeamId: t.Id,
			Name:   t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]TeamVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
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
		err := teamsrv.Outer.UpsertUser(c, teamsrv.UpsertUserReqDTO{
			TeamId:   req.TeamId,
			Account:  req.Account,
			RoleId:   req.GroupId,
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
		users, next, err := teamsrv.Outer.ListUser(c, teamsrv.ListUserReqDTO{
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
		data, _ := listutil.Map(users, func(t teamsrv.UserDTO) (TeamUserVO, error) {
			return TeamUserVO{
				TeamId:    t.TeamId,
				Account:   t.Account,
				GroupId:   t.RoleId,
				GroupName: t.RoleName,
				Created:   t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[TeamUserVO]{
			DataResp: ginutil.DataResp[[]TeamUserVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}

func deleteTeamUser(c *gin.Context) {
	var req DeleteTeamUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.DeleteUser(c, teamsrv.DeleteUserReqDTO{
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
		err := teamsrv.Outer.InsertRole(c, teamsrv.InsertRoleReqDTO{
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
		err := teamsrv.Outer.UpdateRoleName(c, teamsrv.UpdateRoleNameReqDTO{
			RoleId:   req.GroupId,
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
		err := teamsrv.Outer.UpdateRolePerm(c, teamsrv.UpdateRolePermReqDTO{
			RoleId:   req.GroupId,
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
		err := teamsrv.Outer.DeleteRole(c, teamsrv.DeleteRoleReqDTO{
			RoleId:   req.GroupId,
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
		groups, err := teamsrv.Outer.ListRole(c, teamsrv.ListRoleReqDTO{
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret, _ := listutil.Map(groups, func(t teamsrv.RoleDTO) (TeamUserGroupVO, error) {
			return TeamUserGroupVO{
				GroupId: t.RoleId,
				TeamId:  t.TeamId,
				Name:    t.Name,
				Perm:    t.Perm,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]TeamUserGroupVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}
