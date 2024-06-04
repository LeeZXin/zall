package teamsrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
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
	RelationId int64               `json:"relationId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *DeleteUserReqDTO) IsValid() error {
	if r.RelationId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListRoleUserReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListRoleUserReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreateUserReqDTO struct {
	RoleId   int64               `json:"roleId"`
	Accounts []string            `json:"accounts"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateUserReqDTO) IsValid() error {
	if len(r.Accounts) == 0 || len(r.Accounts) > 1000 {
		return util.InvalidArgsError()
	}
	for _, account := range r.Accounts {
		if !usermd.IsAccountValid(account) {
			return util.InvalidArgsError()
		}
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RoleId <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type CreateRoleReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Name     string              `json:"name"`
	Perm     perm.Detail         `json:"perm"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateRoleReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !teammd.IsRoleNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Perm.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateRoleReqDTO struct {
	RoleId   int64               `json:"roleId"`
	Name     string              `json:"name"`
	Perm     perm.Detail         `json:"perm"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateRoleReqDTO) IsValid() error {
	if r.RoleId <= 0 {
		return util.InvalidArgsError()
	}
	if !teammd.IsRoleNameValid(r.Name) {
		return util.InvalidArgsError()
	}
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
	RoleId  int64       `json:"roleId"`
	TeamId  int64       `json:"teamId"`
	Name    string      `json:"name"`
	Perm    perm.Detail `json:"perm"`
	IsAdmin bool        `json:"isAdmin"`
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
	Account string `json:"account"`
	Name    string `json:"name"`
}

type RoleUserDTO struct {
	Id       int64  `json:"id"`
	Account  string `json:"account"`
	Name     string `json:"name"`
	RoleId   int64  `json:"roleId"`
	RoleName string `json:"roleName"`
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

type ListUserByTeamIdReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListUserByTeamIdReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ChangeRoleReqDTO struct {
	RelationId int64               `json:"relationId"`
	RoleId     int64               `json:"roleId"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *ChangeRoleReqDTO) IsValid() error {
	if r.RelationId <= 0 {
		return util.InvalidArgsError()
	}
	if r.RoleId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
