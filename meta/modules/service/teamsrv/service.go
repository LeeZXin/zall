package teamsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/util"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = &innerImpl{
		permCache: util.NewGoCache(),
	}
)

type InnerService interface {
	GetUserPermDetail(context.Context, int64, string) (teammd.UserPermDetailDTO, bool)
}

type OuterService interface {
	InsertTeam(context.Context, InsertTeamReqDTO) error
	UpdateTeam(context.Context, UpdateTeamReqDTO) error
	DeleteTeam(context.Context, DeleteTeamReqDTO) error
	ListUser(context.Context, ListUserReqDTO) ([]UserDTO, int64, error)
	DeleteUser(context.Context, DeleteUserReqDTO) error
	UpsertUser(context.Context, UpsertUserReqDTO) error
	InsertRole(context.Context, InsertRoleReqDTO) error
	UpdateRoleName(context.Context, UpdateRoleNameReqDTO) error
	UpdateRolePerm(context.Context, UpdateRolePermReqDTO) error
	DeleteRole(context.Context, DeleteRoleReqDTO) error
	ListRole(context.Context, ListRoleReqDTO) ([]RoleDTO, error)
	ListTeam(context.Context, ListTeamReqDTO) ([]TeamDTO, error)
}
