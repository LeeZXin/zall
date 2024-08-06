package usersrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
	"time"
)

type CreateUserReqDTO struct {
	Account   string              `json:"account"`
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	AvatarUrl string              `json:"avatarUrl"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *CreateUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !usermd.IsEmailValid(r.Email) {
		return util.InvalidArgsError()
	}
	if !usermd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	if !usermd.IsUsernameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	parsedUrl, err := url.Parse(r.AvatarUrl)
	if err != nil || strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	return nil
}

type RegisterUserReqDTO struct {
	Account  string `json:"account"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !usermd.IsEmailValid(r.Email) {
		return util.InvalidArgsError()
	}
	if !usermd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	if !usermd.IsUsernameValid(r.Name) {
		return util.InvalidArgsError()
	}
	return nil
}

type LoginReqDTO struct {
	Account  string `json:"account"`
	Password string `json:"-"`
}

func (r *LoginReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !usermd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	return nil
}

type RefreshReqDTO struct {
	SessionId string              `json:"sessionId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *RefreshReqDTO) IsValid() error {
	if r.SessionId == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type LogoutReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *LogoutReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SetProhibitedReqDTO struct {
	Account      string              `json:"account"`
	IsProhibited bool                `json:"isProhibited"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *SetProhibitedReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteUserReqDTO struct {
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
	Account  string              `json:"account"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListUserReqDTO) IsValid() error {
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if r.Account != "" && len(r.Account) > 32 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UserDTO struct {
	Account      string    `json:"account"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	IsAdmin      bool      `json:"isAdmin"`
	IsProhibited bool      `json:"isProhibited"`
	AvatarUrl    string    `json:"avatarUrl"`
	Created      time.Time `json:"created"`
	IsDba        bool      `json:"isDba"`
}

type UpdateUserReqDTO struct {
	Account   string              `json:"account"`
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	AvatarUrl string              `json:"avatarUrl"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *UpdateUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !usermd.IsUsernameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !usermd.IsEmailValid(r.Email) {
		return util.InvalidArgsError()
	}
	parsedUrl, err := url.Parse(r.AvatarUrl)
	if err != nil || strings.HasPrefix(parsedUrl.Scheme, "http") {
		return util.InvalidArgsError()
	}
	return nil
}

type SetAdminReqDTO struct {
	Account  string              `json:"account"`
	IsAdmin  bool                `json:"isAdmin"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *SetAdminReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SetDbaReqDTO struct {
	Account  string              `json:"account"`
	IsDba    bool                `json:"isDba"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *SetDbaReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ResetPasswordReqDTO struct {
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ResetPasswordReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdatePasswordReqDTO struct {
	Password string              `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdatePasswordReqDTO) IsValid() error {
	if !usermd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertUserCorpReqDTO struct {
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertUserCorpReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteUserCorpReqDTO struct {
	Account  string              `json:"account"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteUserCorpReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CheckAccountAndPasswordReqDTO struct {
	Account  string `json:"account"`
	Password string `json:"-"`
}

func (r *CheckAccountAndPasswordReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !usermd.IsPasswordValid(r.Password) {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAllUserReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAllUserReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SimpleUserDTO struct {
	Account string
	Name    string
}
