package userapi

import (
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type LoginReqVO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginRespVO struct {
	ginutil.BaseResp
	Session apisession.Session `json:"session"`
}

type RefreshRespVO struct {
	ginutil.BaseResp
	SessionId string `json:"sessionId"`
	ExpireAt  int64  `json:"expireAt"`
}

type CreateUserReqVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatarUrl"`
}

type RegisterUserReqVO struct {
	Account  string `json:"account"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserVO struct {
	Account      string `json:"account"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"isAdmin"`
	IsProhibited bool   `json:"isProhibited"`
	AvatarUrl    string `json:"avatarUrl"`
	Created      string `json:"created"`
	IsDba        bool   `json:"isDba"`
}

type ListUserReqVO struct {
	Account string `json:"account"`
	PageNum int    `json:"pageNum"`
}

type UpdateUserReqVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatarUrl"`
}

type SetAdminReqVO struct {
	Account string `json:"account"`
	IsAdmin bool   `json:"isAdmin"`
}

type UpdatePasswordReqVO struct {
	Origin   string `json:"origin"`
	Password string `json:"password"`
}

type SetProhibitedReqVO struct {
	Account      string `json:"account"`
	IsProhibited bool   `json:"isProhibited"`
}

type SetDbaReqVO struct {
	Account string `json:"account"`
	IsDba   bool   `json:"isDba"`
}

type SimpleUserVO struct {
	Account string `json:"account"`
	Name    string `json:"name"`
}
