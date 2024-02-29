package teammd

import "github.com/LeeZXin/zall/pkg/perm"

type InsertTeamReqDTO struct {
	Name string
}

type InsertTeamUserReqDTO struct {
	TeamId  int64
	Account string
	GroupId int64
}

type UpdateTeamUserReqDTO struct {
	TeamId  int64
	Account string
	GroupId int64
}

type ListTeamUserReqDTO struct {
	TeamId  int64
	Account string
	Cursor  int64
	Limit   int
}

type InsertTeamUserGroupReqDTO struct {
	Name       string
	TeamId     int64
	PermDetail perm.Detail
	IsAdmin    bool
}

type TeamUserPermDetailDTO struct {
	GroupId    int64
	IsAdmin    bool
	PermDetail perm.Detail
}

type UpdateTeamReqDTO struct {
	TeamId int64
	Name   string
}
