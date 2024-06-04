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
		// 团队用户
		group = e.Group("/api/teamUser", apisession.CheckLogin)
		{
			// 创建团队成员关系
			group.POST("/create", createTeamUser)
			// 展示团队角色列表
			group.GET("/listRoleUser/:teamId", listRoleUser)
			// 获取团队成员名称
			group.GET("/listByTeamId/:teamId", listUserByTeamId)
			// 删除成员团队绑定关系
			group.DELETE("/delete/:relationId", deleteTeamUser)
			// 更换角色
			group.POST("/change", changeRole)
		}
		// 项目用户组
		group = e.Group("/api/teamRole", apisession.CheckLogin)
		{
			// 创建角色
			group.POST("/create", createRole)
			// 团队角色列表
			group.GET("/list/:teamId", listRole)
			// 编辑角色
			group.POST("/update", updateRole)
			// 删除角色
			group.DELETE("/delete/:roleId", deleteRole)
		}
	})
}

func changeRole(c *gin.Context) {
	var req ChangeRoleReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.ChangeRole(c, teamsrv.ChangeRoleReqDTO{
			RelationId: req.RelationId,
			RoleId:     req.RoleId,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
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

func listUserByTeamId(c *gin.Context) {
	teamId := cast.ToInt64(c.Param("teamId"))
	users, err := teamsrv.Outer.ListUserByTeamId(c, teamsrv.ListUserByTeamIdReqDTO{
		TeamId:   teamId,
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(users, func(t teamsrv.UserDTO) (UserVO, error) {
		return UserVO{
			Account: t.Account,
			Name:    t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]UserVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
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

func createTeamUser(c *gin.Context) {
	var req CreateTeamUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.CreateUser(c, teamsrv.CreateUserReqDTO{
			Accounts: req.Accounts,
			RoleId:   req.RoleId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listRoleUser(c *gin.Context) {
	users, err := teamsrv.Outer.ListRoleUser(c, teamsrv.ListRoleUserReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(users, func(t teamsrv.RoleUserDTO) (RoleUserVO, error) {
		return RoleUserVO{
			Id:       t.Id,
			Account:  t.Account,
			Name:     t.Name,
			RoleId:   t.RoleId,
			RoleName: t.RoleName,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]RoleUserVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func deleteTeamUser(c *gin.Context) {
	err := teamsrv.Outer.DeleteUser(c, teamsrv.DeleteUserReqDTO{
		RelationId: cast.ToInt64(c.Param("relationId")),
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func createRole(c *gin.Context) {
	var req CreateRoleReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.CreateRole(c, teamsrv.CreateRoleReqDTO{
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

func updateRole(c *gin.Context) {
	var req UpdateRoleReqVO
	if util.ShouldBindJSON(&req, c) {
		err := teamsrv.Outer.UpdateRole(c, teamsrv.UpdateRoleReqDTO{
			RoleId:   req.RoleId,
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

func deleteRole(c *gin.Context) {
	err := teamsrv.Outer.DeleteRole(c, teamsrv.DeleteRoleReqDTO{
		RoleId:   cast.ToInt64(c.Param("roleId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listRole(c *gin.Context) {
	groups, err := teamsrv.Outer.ListRole(c, teamsrv.ListRoleReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	ret, _ := listutil.Map(groups, func(t teamsrv.RoleDTO) (RoleVO, error) {
		return RoleVO{
			RoleId:  t.RoleId,
			TeamId:  t.TeamId,
			Name:    t.Name,
			Perm:    t.Perm,
			IsAdmin: t.IsAdmin,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]RoleVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     ret,
	})
}
