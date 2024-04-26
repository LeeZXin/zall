package teamsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
	"time"
)

type CreateTeamReqDTO struct {
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateTeamReqDTO) IsValid() error {
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
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type IsAdminReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *IsAdminReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetTeamPermReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetTeamPermReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Account  string              `json:"account"`
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListUserReqDTO) IsValid() error {
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

type UpsertUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Account  string              `json:"account"`
	RoleId   int64               `json:"roleId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpsertUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RoleId <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertRoleReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Perm     perm.Detail         `json:"perm"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertRoleReqDTO) IsValid() error {
	if !teammd.IsRoleNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateRoleNameReqDTO struct {
	RoleId   int64               `json:"roleId"`
	Name     string              `json:"name"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateRoleNameReqDTO) IsValid() error {
	if !teammd.IsRoleNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateRolePermReqDTO struct {
	RoleId   int64               `json:"roleId"`
	Perm     perm.Detail         `json:"perm"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateRolePermReqDTO) IsValid() error {
	if !r.Perm.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteRoleReqDTO struct {
	RoleId   int64               `json:"roleId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteRoleReqDTO) IsValid() error {
	if r.RoleId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListRoleReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListRoleReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RoleDTO struct {
	RoleId int64       `json:"roleId"`
	TeamId int64       `json:"teamId"`
	Name   string      `json:"name"`
	Perm   perm.Detail `json:"perm"`
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

type UserDTO struct {
	TeamId   int64     `json:"teamId"`
	Account  string    `json:"account"`
	RoleId   int64     `json:"roleId"`
	RoleName string    `json:"roleName"`
	Created  time.Time `json:"created"`
}

type GetTeamReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetTeamReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
