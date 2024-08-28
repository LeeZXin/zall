package userapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/login")
		{
			// 登录
			group.POST("/login", login)
			// 注册用户
			group.POST("/register", register)
			// 获取登录信息
			group.Any("/userInfo", apisession.CheckLogin, getUserInfo)
			// 刷新token
			group.Any("/refresh", apisession.CheckLogin, refresh)
			// 退出登录
			group.Any("/logout", apisession.CheckLogin, logout)
		}
		group = e.Group("/api/user", apisession.CheckLogin)
		{
			// 新增用户
			group.POST("/create", createUser)
			// 删除用户
			group.DELETE("/delete/:account", deleteUser)
			// 更新用户
			group.POST("/update", updateUser)
			// 展示用户列表 管理员权限
			group.GET("/list", listUser)
			// 更新密码
			group.POST("/updatePassword", updatePassword)
			// 系统管理员设置
			group.PUT("/setAdmin", setAdmin)
			// 禁用用户
			group.PUT("/setProhibited", setProhibited)
			// 设置dba
			group.PUT("/setDba", setDba)
			// 展示所有用户列表
			group.GET("/listAll", listAll)
			// 重置密码
			group.PUT("/resetPassword/:account", resetPassword)
		}

	})
}

func listAll(c *gin.Context) {
	users, err := usersrv.ListAllUser(c, usersrv.ListAllUserReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data := listutil.MapNe(users, func(t usersrv.SimpleUserDTO) SimpleUserVO {
		return SimpleUserVO{
			Account: t.Account,
			Name:    t.Name,
		}
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleUserVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func login(c *gin.Context) {
	var req LoginReqVO
	if util.ShouldBindJSON(&req, c) {
		session, err := usersrv.Login(c, usersrv.LoginReqDTO{
			Account:  req.Account,
			Password: req.Password,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.SetCookie(apisession.LoginCookie, session.SessionId, int(usersrv.LoginSessionExpiry.Seconds()), "/", "", false, true)
		c.JSON(http.StatusOK, LoginRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Session:  session,
		})
	}
}

func getUserInfo(c *gin.Context) {
	session := apisession.MustGetSession(c)
	c.JSON(http.StatusOK, LoginRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Session:  session,
	})
}

func refresh(c *gin.Context) {
	sessionId, expireAt, err := usersrv.Refresh(c, usersrv.RefreshReqDTO{
		SessionId: apisession.GetSessionId(c),
		Operator:  apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.SetCookie(apisession.LoginCookie, sessionId, int(usersrv.LoginSessionExpiry.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, RefreshRespVO{
		BaseResp:  ginutil.DefaultSuccessResp,
		SessionId: sessionId,
		ExpireAt:  expireAt,
	})
}

func logout(c *gin.Context) {
	err := usersrv.Logout(c, usersrv.LogoutReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func createUser(c *gin.Context) {
	var req CreateUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.CreateUser(c, usersrv.CreateUserReqDTO{
			Account:   req.Account,
			Name:      req.Name,
			Email:     req.Email,
			Password:  req.Password,
			AvatarUrl: req.AvatarUrl,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func deleteUser(c *gin.Context) {
	err := usersrv.DeleteUser(c, usersrv.DeleteUserReqDTO{
		Account:  c.Param("account"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func updateUser(c *gin.Context) {
	var req UpdateUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.UpdateUser(c, usersrv.UpdateUserReqDTO{
			Account:   req.Account,
			Name:      req.Name,
			Email:     req.Email,
			AvatarUrl: req.AvatarUrl,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func listUser(c *gin.Context) {
	var req ListUserReqVO
	if util.ShouldBindQuery(&req, c) {
		users, total, err := usersrv.ListUser(c, usersrv.ListUserReqDTO{
			Account:  req.Account,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data := listutil.MapNe(users, func(t usersrv.UserDTO) UserVO {
			return UserVO{
				Account:      t.Account,
				Name:         t.Name,
				Email:        t.Email,
				IsAdmin:      t.IsAdmin,
				IsProhibited: t.IsProhibited,
				AvatarUrl:    t.AvatarUrl,
				Created:      t.Created.Format(time.DateTime),
				IsDba:        t.IsDba,
			}
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[UserVO]{
			DataResp: ginutil.DataResp[[]UserVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func updatePassword(c *gin.Context) {
	var req UpdatePasswordReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.UpdatePassword(c, usersrv.UpdatePasswordReqDTO{
			Origin:   req.Origin,
			Password: req.Password,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func register(c *gin.Context) {
	var req RegisterUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.RegisterUser(c, usersrv.RegisterUserReqDTO{
			Account:  req.Account,
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func setAdmin(c *gin.Context) {
	var req SetAdminReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.SetAdmin(c, usersrv.SetAdminReqDTO{
			Account:  req.Account,
			IsAdmin:  req.IsAdmin,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func setProhibited(c *gin.Context) {
	var req SetProhibitedReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.SetProhibited(c, usersrv.SetProhibitedReqDTO{
			Account:      req.Account,
			IsProhibited: req.IsProhibited,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func setDba(c *gin.Context) {
	var req SetDbaReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.SetDba(c, usersrv.SetDbaReqDTO{
			Account:  req.Account,
			IsDba:    req.IsDba,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
		} else {
			util.DefaultOkResponse(c)
		}
	}
}

func resetPassword(c *gin.Context) {
	err := usersrv.ResetPassword(c, usersrv.ResetPasswordReqDTO{
		Account:  c.Param("account"),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
	} else {
		util.DefaultOkResponse(c)
	}
}
