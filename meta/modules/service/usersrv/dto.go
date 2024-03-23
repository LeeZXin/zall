package usersrv

import (
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"regexp"
	"time"
)

var (
	validPasswordPattern     = regexp.MustCompile("\\S{6,}")
	validUserEmailRegPattern = regexp.MustCompile(`^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`)
)

type InsertUserReqDTO struct {
	Account   string              `json:"account"`
	Name      string              `json:"name"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	AvatarUrl string              `json:"avatarUrl"`
	IsAdmin   bool                `json:"isAdmin"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *InsertUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateUserEmail(r.Email) {
		return util.InvalidArgsError()
	}
	if !validatePassword(r.Password) {
		return util.InvalidArgsError()
	}
	if len(r.Name) > 32 || len(r.Name) == 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RegisterUserReqDTO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatarUrl"`
}

func (r *RegisterUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateUserEmail(r.Email) {
		return util.InvalidArgsError()
	}
	if !validPasswordPattern.MatchString(r.Password) {
		return util.InvalidArgsError()
	}
	if !validateUserName(r.Name) {
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
	if !validatePassword(r.Password) {
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

type LoginOutReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *LoginOutReqDTO) IsValid() error {
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
	Cursor   int64               `json:"cursor"`
	Limit    int                 `json:"limit"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListUserReqDTO) IsValid() error {
	if r.Cursor < 0 {
		return util.InvalidArgsError()
	}
	if r.Limit <= 0 || r.Limit > 1000 {
		return util.InvalidArgsError()
	}
	if r.Account != "" && !usermd.IsAccountValid(r.Account) {
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
	Updated      time.Time `json:"updated"`
}

type UpdateUserReqDTO struct {
	Account  string              `json:"account"`
	Name     string              `json:"name"`
	Email    string              `json:"email"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateUserReqDTO) IsValid() error {
	if !usermd.IsAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !validateUserName(r.Name) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !validateUserEmail(r.Email) {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateAdminReqDTO struct {
	Account  string              `json:"account"`
	IsAdmin  bool                `json:"isAdmin"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *UpdateAdminReqDTO) IsValid() error {
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
	if !validatePassword(r.Password) {
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
	if !validatePassword(r.Password) {
		return util.InvalidArgsError()
	}
	return nil
}

func validateUserName(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func validateUserEmail(email string) bool {
	return validUserEmailRegPattern.MatchString(email)
}

func validatePassword(password string) bool {
	return validPasswordPattern.MatchString(password)
}
