package teamsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/perm"
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
	// CreateTeam 创建团队
	CreateTeam(context.Context, CreateTeamReqDTO) error
	UpdateTeam(context.Context, UpdateTeamReqDTO) error
	DeleteTeam(context.Context, DeleteTeamReqDTO) error
	// IsAdmin 是否是团队管理员
	IsAdmin(context.Context, IsAdminReqDTO) (bool, error)
	// GetTeamPerm 获取团队权限
	GetTeamPerm(context.Context, GetTeamPermReqDTO) (perm.TeamPerm, error)
	// GetTeam 获取团队信息
	GetTeam(context.Context, GetTeamReqDTO) (teammd.Team, error)
	// ListAccount 获取成员账号
	ListAccount(context.Context, ListAccountReqDTO) ([]string, error)
	ListUser(context.Context, ListUserReqDTO) ([]UserDTO, int64, error)
	DeleteUser(context.Context, DeleteUserReqDTO) error
	UpsertUser(context.Context, UpsertUserReqDTO) error
	InsertRole(context.Context, InsertRoleReqDTO) error
	UpdateRoleName(context.Context, UpdateRoleNameReqDTO) error
	UpdateRolePerm(context.Context, UpdateRolePermReqDTO) error
	DeleteRole(context.Context, DeleteRoleReqDTO) error
	ListRole(context.Context, ListRoleReqDTO) ([]RoleDTO, error)
	// ListTeam 展示用户所在团队列表
	ListTeam(context.Context, ListTeamReqDTO) ([]teammd.Team, error)
}
