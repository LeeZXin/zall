package userapi

import (
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
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
		group := e.Group("/api/login")
		{
			// 登录
			group.POST("/login", login)
			// 注册用户
			group.POST("/register", register)
			// 退出登录
			group.Any("/loginOut", apisession.CheckLogin, loginOut)
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
		sessionId, err := usersrv.Outer.Login(c, usersrv.LoginReqDTO{
			Account:  req.Account,
			Password: req.Password,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.SetCookie(apisession.LoginCookie, sessionId, int(usersrv.LoginSessionExpiry.Seconds()), "/", "", false, true)
		c.JSON(http.StatusOK, LoginRespVO{
			BaseResp:  ginutil.DefaultSuccessResp,
			SessionId: sessionId,
		})
	}
}

func loginOut(c *gin.Context) {
	err := usersrv.Outer.LoginOut(c, usersrv.LoginOutReqDTO{
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
		respDTO, err := usersrv.Outer.ListUser(c, usersrv.ListUserReqDTO{
			Account:  req.Account,
			Cursor:   req.Cursor,
			Limit:    req.Limit,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListUserRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   respDTO.Cursor,
		}
		ret.UserList, _ = listutil.Map(respDTO.UserList, func(t usersrv.UserDTO) (UserVO, error) {
			return UserVO{
				Account:      t.Account,
				Name:         t.Name,
				Email:        t.Email,
				IsAdmin:      t.IsAdmin,
				IsProhibited: t.IsProhibited,
				AvatarUrl:    t.AvatarUrl,
				Created:      t.Created.Format(timeutil.DefaultTimeFormat),
				Updated:      t.Updated.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, ret)
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
			Account:   req.Account,
			Name:      req.Name,
			Email:     req.Email,
			Password:  req.Password,
			AvatarUrl: req.AvatarUrl,
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
