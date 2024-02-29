package teamapi

import (
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InsertTeamReqVO struct {
	Name string `json:"name"`
}

type UpdateTeamReqVO struct {
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
}

type UpsertTeamUserReqVO struct {
	TeamId  int64  `json:"teamId"`
	Account string `json:"account"`
	GroupId int64  `json:"groupId"`
}

type DeleteTeamUserReqVO struct {
	TeamId  int64  `json:"teamId"`
	Account string `json:"account"`
}

type ListTeamUserReqVO struct {
	TeamId  int64  `json:"teamId"`
	Account string `json:"account"`
	Cursor  int64  `json:"cursor"`
	Limit   int    `json:"limit"`
}

type ListTeamUserRespVO struct {
	ginutil.BaseResp
	Data []TeamUserVO `json:"data"`
	Next int64        `json:"next"`
}

type TeamUserVO struct {
	TeamId    int64  `json:"teamId"`
	Account   string `json:"account"`
	GroupId   int64  `json:"groupId"`
	GroupName string `json:"groupName"`
	Created   string `json:"created"`
}

type InsertTeamUserGroupReqVO struct {
	TeamId int64       `json:"teamId"`
	Name   string      `json:"name"`
	Perm   perm.Detail `json:"perm"`
}

type UpdateTeamUserGroupNameReqVO struct {
	GroupId int64  `json:"groupId"`
	Name    string `json:"name"`
}

type UpdateTeamUserGroupPermReqVO struct {
	GroupId int64       `json:"groupId"`
	Perm    perm.Detail `json:"perm"`
}

type DeleteTeamUserGroupReqVO struct {
	GroupId int64 `json:"groupId"`
}

type ListTeamUserGroupReqVO struct {
	TeamId int64 `json:"teamId"`
}

type TeamUserGroupVO struct {
	GroupId int64       `json:"groupId"`
	TeamId  int64       `json:"teamId"`
	Name    string      `json:"name"`
	Perm    perm.Detail `json:"perm"`
}

type ListTeamUserGroupRespVO struct {
	ginutil.BaseResp
	Data []TeamUserGroupVO `json:"data"`
}

type TeamVO struct {
	TeamId  int64  `json:"teamId"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Created string `json:"created"`
}

type ListTeamRespVO struct {
	ginutil.BaseResp
	Data []TeamVO `json:"data"`
}

type DeleteTeamReqVO struct {
	TeamId int64 `json:"teamId"`
}
