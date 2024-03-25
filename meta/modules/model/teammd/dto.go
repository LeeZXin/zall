package teammd

import "github.com/LeeZXin/zall/pkg/perm"

type InsertTeamReqDTO struct {
	Name string
}

type InsertUserReqDTO struct {
	TeamId  int64
	Account string
	RoleId  int64
}

type UpdateUserReqDTO struct {
	TeamId  int64
	Account string
	RoleId  int64
}

type ListUserReqDTO struct {
	TeamId  int64
	Account string
	Cursor  int64
	Limit   int
}

type InsertRoleReqDTO struct {
	Name       string
	TeamId     int64
	PermDetail perm.Detail
	IsAdmin    bool
}

type UserPermDetailDTO struct {
	RoleId     int64
	IsAdmin    bool
	PermDetail perm.Detail
}

type UpdateTeamReqDTO struct {
	TeamId int64
	Name   string
}
