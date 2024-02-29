package usersrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/util"
)

var (
	Inner InnerService = &innerImpl{
		userCache: util.NewGoCache(),
	}
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	GetByAccount(context.Context, string) (usermd.UserInfo, bool)
	CheckAccountAndPassword(context.Context, CheckAccountAndPasswordReqDTO) (usermd.UserInfo, bool)
}

type OuterService interface {
	Login(context.Context, LoginReqDTO) (string, error)
	LoginOut(context.Context, LoginOutReqDTO) error
	RegisterUser(context.Context, RegisterUserReqDTO) error
	InsertUser(context.Context, InsertUserReqDTO) error
	DeleteUser(context.Context, DeleteUserReqDTO) error
	ListUser(context.Context, ListUserReqDTO) (ListUserRespDTO, error)
	UpdateUser(context.Context, UpdateUserReqDTO) error
	// SetProhibited 禁用用户
	SetProhibited(context.Context, SetProhibitedReqDTO) error
	UpdatePassword(context.Context, UpdatePasswordReqDTO) error
	// UpdateAdmin 变更管理员
	UpdateAdmin(context.Context, UpdateAdminReqDTO) error
}
