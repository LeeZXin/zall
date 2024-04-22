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
			group.POST("/insert", insertUser)
			// 删除用户
			group.POST("/deleteUser", deleteUser)
			// 更新用户
			group.POST("/update", updateUser)
			// 展示用户列表
			group.POST("/list", listUser)
			// 更新密码
			group.POST("/updatePassword", updatePassword)
			// 系统管理员设置
			group.POST("/updateAdmin", updateAdmin)
			// 禁用用户
			group.POST("/setProhibited", setProhibited)
		}
	})
}

func login(c *gin.Context) {
	var req LoginReqVO
	if util.ShouldBindJSON(&req, c) {
		session, err := usersrv.Outer.Login(c, usersrv.LoginReqDTO{
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
	sessionId, expireAt, err := usersrv.Outer.Refresh(c, usersrv.RefreshReqDTO{
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
	err := usersrv.Outer.Logout(c, usersrv.LogoutReqDTO{
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func insertUser(c *gin.Context) {
	var reqVO InsertUserReqVO
	if util.ShouldBindJSON(&reqVO, c) {
		err := usersrv.Outer.InsertUser(c, usersrv.InsertUserReqDTO{
			Account:   reqVO.Account,
			Name:      reqVO.Name,
			Email:     reqVO.Email,
			Password:  reqVO.Password,
			AvatarUrl: reqVO.AvatarUrl,
			IsAdmin:   reqVO.IsAdmin,
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
	var req DeleteUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.Outer.DeleteUser(c, usersrv.DeleteUserReqDTO{
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

func updateUser(c *gin.Context) {
	var req UpdateUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.Outer.UpdateUser(c, usersrv.UpdateUserReqDTO{
			Account:  req.Account,
			Name:     req.Name,
			Email:    req.Email,
			Operator: apisession.MustGetLoginUser(c),
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
	if util.ShouldBindJSON(&req, c) {
		users, next, err := usersrv.Outer.ListUser(c, usersrv.ListUserReqDTO{
			Account:  req.Account,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(users, func(t usersrv.UserDTO) (UserVO, error) {
			return UserVO{
				Account:      t.Account,
				Name:         t.Name,
				Email:        t.Email,
				IsAdmin:      t.IsAdmin,
				IsProhibited: t.IsProhibited,
				AvatarUrl:    t.AvatarUrl,
				Created:      t.Created.Format(time.DateTime),
				Updated:      t.Updated.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.PageResp[[]UserVO]{
			DataResp: ginutil.DataResp[[]UserVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			Next: next,
		})
	}
}

func updatePassword(c *gin.Context) {
	var req UpdatePasswordReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.Outer.UpdatePassword(c, usersrv.UpdatePasswordReqDTO{
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
		err := usersrv.Outer.RegisterUser(c, usersrv.RegisterUserReqDTO{
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

func updateAdmin(c *gin.Context) {
	var req UpdateAdminReqVO
	if util.ShouldBindJSON(&req, c) {
		err := usersrv.Outer.UpdateAdmin(c, usersrv.UpdateAdminReqDTO{
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
		err := usersrv.Outer.SetProhibited(c, usersrv.SetProhibitedReqDTO{
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
