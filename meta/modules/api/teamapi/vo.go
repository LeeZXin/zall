package teamapi

import (
	"github.com/LeeZXin/zall/pkg/perm"
)

type CreateTeamReqVO struct {
	Name string `json:"name"`
}

type UpdateTeamReqVO struct {
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
}

type CreateTeamUserReqVO struct {
	RoleId   int64    `json:"roleId"`
	Accounts []string `json:"accounts"`
}

type UserVO struct {
	Account string `json:"account"`
	Name    string `json:"name"`
}

type RoleUserVO struct {
	Id       int64  `json:"id"`
	Account  string `json:"account"`
	Name     string `json:"name"`
	RoleId   int64  `json:"roleId"`
	RoleName string `json:"roleName"`
}

type CreateRoleReqVO struct {
	TeamId int64       `json:"teamId"`
	Name   string      `json:"name"`
	Perm   perm.Detail `json:"perm"`
}

type UpdateRoleReqVO struct {
	RoleId int64       `json:"roleId"`
	Name   string      `json:"name"`
	Perm   perm.Detail `json:"perm"`
}

type RoleVO struct {
	RoleId  int64       `json:"roleId"`
	TeamId  int64       `json:"teamId"`
	Name    string      `json:"name"`
	Perm    perm.Detail `json:"perm"`
	IsAdmin bool        `json:"isAdmin"`
}

type TeamVO struct {
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
}

type TeamWithPermVO struct {
	TeamId  int64         `json:"teamId"`
	Name    string        `json:"name"`
	IsAdmin bool          `json:"isAdmin"`
	Perm    perm.TeamPerm `json:"perm"`
}

type ChangeRoleReqVO struct {
	RelationId int64 `json:"relationId"`
	RoleId     int64 `json:"roleId"`
}
