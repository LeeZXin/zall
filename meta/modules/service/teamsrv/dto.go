package teamsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
	"time"
)

type InsertTeamReqDTO struct {
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertTeamReqDTO) IsValid() error {
	if !teammd.IsTeamNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTeamReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTeamReqDTO) IsValid() error {
	if !teammd.IsTeamNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTeamReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTeamReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTeamUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTeamUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTeamUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Account  string              `json:"account"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTeamUserReqDTO) IsValid() error {
	if r.Account != "" && !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	return nil
}

type UpsertTeamUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Account  string              `json:"account"`
	GroupId  int64               `json:"groupId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpsertTeamUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.GroupId <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertTeamUserGroupReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Perm     perm.Detail         `json:"perm"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertTeamUserGroupReqDTO) IsValid() error {
	if !teammd.IsGroupNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTeamUserGroupNameReqDTO struct {
	GroupId  int64               `json:"groupId"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTeamUserGroupNameReqDTO) IsValid() error {
	if !teammd.IsGroupNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateTeamUserGroupPermReqDTO struct {
	GroupId  int64               `json:"groupId"`
	Perm     perm.Detail         `json:"perm"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateTeamUserGroupPermReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTeamUserGroupReqDTO struct {
	GroupId  int64               `json:"groupId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTeamUserGroupReqDTO) IsValid() error {
	if r.GroupId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTeamUserGroupReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTeamUserGroupReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TeamUserGroupDTO struct {
	GroupId int64       `json:"groupId"`
	TeamId  int64       `json:"teamId"`
	Name    string      `json:"name"`
	Perm    perm.Detail `json:"perm"`
}

type ListTeamReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTeamReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TeamDTO struct {
	TeamId  int64     `json:"teamId"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type TeamUserDTO struct {
	TeamId    int64     `json:"teamId"`
	Account   string    `json:"account"`
	GroupId   int64     `json:"groupId"`
	GroupName string    `json:"groupName"`
	Created   time.Time `json:"created"`
}
