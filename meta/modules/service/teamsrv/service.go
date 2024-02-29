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
	GetTeamUserPermDetail(context.Context, int64, string) (teammd.TeamUserPermDetailDTO, bool)
}

type OuterService interface {
	InsertTeam(context.Context, InsertTeamReqDTO) error
	UpdateTeam(context.Context, UpdateTeamReqDTO) error
	DeleteTeam(context.Context, DeleteTeamReqDTO) error
	ListTeamUser(context.Context, ListTeamUserReqDTO) ([]TeamUserDTO, int64, error)
	DeleteTeamUser(context.Context, DeleteTeamUserReqDTO) error
	UpsertTeamUser(context.Context, UpsertTeamUserReqDTO) error
	InsertTeamUserGroup(context.Context, InsertTeamUserGroupReqDTO) error
	UpdateTeamUserGroupName(context.Context, UpdateTeamUserGroupNameReqDTO) error
	UpdateTeamUserGroupPerm(context.Context, UpdateTeamUserGroupPermReqDTO) error
	DeleteTeamUserGroup(context.Context, DeleteTeamUserGroupReqDTO) error
	ListTeamUserGroup(context.Context, ListTeamUserGroupReqDTO) ([]TeamUserGroupDTO, error)
	ListTeam(context.Context, ListTeamReqDTO) ([]TeamDTO, error)
}
