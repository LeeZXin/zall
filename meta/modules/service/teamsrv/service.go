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
	// ListUserByTeamId 获取成员账号
	ListUserByTeamId(context.Context, ListUserByTeamIdReqDTO) ([]UserDTO, error)
	// ListRoleUser 展示团队成员角色信息
	ListRoleUser(context.Context, ListRoleUserReqDTO) ([]RoleUserDTO, error)
	// DeleteUser 删除团队成员关系
	DeleteUser(context.Context, DeleteUserReqDTO) error
	// CreateUser 添加角色成员
	CreateUser(context.Context, CreateUserReqDTO) error
	// CreateRole 创建角色
	CreateRole(context.Context, CreateRoleReqDTO) error
	// UpdateRole 编辑角色
	UpdateRole(context.Context, UpdateRoleReqDTO) error
	// DeleteRole 删除角色
	DeleteRole(context.Context, DeleteRoleReqDTO) error
	// ListRole 角色列表
	ListRole(context.Context, ListRoleReqDTO) ([]RoleDTO, error)
	// ListTeam 展示用户所在团队列表
	ListTeam(context.Context, ListTeamReqDTO) ([]teammd.Team, error)
	// ChangeRole 更换角色
	ChangeRole(context.Context, ChangeRoleReqDTO) error
}
