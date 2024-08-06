package usersrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

var (
	Inner InnerService
	Outer OuterService
)

func Init() {
	if Inner == nil {
		Inner = &innerImpl{
			userCache: util.NewGoCache(),
		}
		Outer = newOuterService()
	}
}

type InnerService interface {
	GetByAccount(context.Context, string) (usermd.UserInfo, bool)
	CheckAccountAndPassword(context.Context, CheckAccountAndPasswordReqDTO) (usermd.UserInfo, bool)
}

type OuterService interface {
	// Login 登录
	Login(context.Context, LoginReqDTO) (apisession.Session, error)
	Refresh(context.Context, RefreshReqDTO) (string, int64, error)
	Logout(context.Context, LogoutReqDTO) error
	// RegisterUser 注册用户
	RegisterUser(context.Context, RegisterUserReqDTO) error
	// CreateUser 管理员创建用户
	CreateUser(context.Context, CreateUserReqDTO) error
	// DeleteUser 删除用户
	DeleteUser(context.Context, DeleteUserReqDTO) error
	// ListUser 管理员查询用户列表
	ListUser(context.Context, ListUserReqDTO) ([]UserDTO, int64, error)
	// UpdateUser 编辑用户
	UpdateUser(context.Context, UpdateUserReqDTO) error
	// SetProhibited 禁用用户
	SetProhibited(context.Context, SetProhibitedReqDTO) error
	UpdatePassword(context.Context, UpdatePasswordReqDTO) error
	// SetAdmin 变更管理员
	SetAdmin(context.Context, SetAdminReqDTO) error
	// SetDba 设置dba角色
	SetDba(context.Context, SetDbaReqDTO) error
	// ListAllUser 管理员角色获取所有用户列表
	ListAllUser(context.Context, ListAllUserReqDTO) ([]SimpleUserDTO, error)
	// ResetPassword 重置密码
	ResetPassword(context.Context, ResetPasswordReqDTO) error
}
