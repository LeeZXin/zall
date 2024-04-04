package usermd

import (
	"time"
)

const (
	UserTableName = "zall_user"
)

type RoleType int

const (
	DeveloperRole RoleType = 1
	DbaRole       RoleType = 1 << 1
)

func (t RoleType) IsDeveloper() bool {
	return t&DeveloperRole == DeveloperRole
}

func (t RoleType) IsDba() bool {
	return t&DbaRole == DbaRole
}

type User struct {
	Id           int64     `json:"id" xorm:"pk autoincr"`
	Account      string    `json:"account"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	IsProhibited bool      `json:"isProhibited"`
	AvatarUrl    string    `json:"avatarUrl"`
	IsAdmin      bool      `json:"isAdmin"`
	RoleType     RoleType  `json:"roleType"`
	Created      time.Time `json:"created" xorm:"created"`
	Updated      time.Time `json:"updated" xorm:"updated"`
}

func (*User) TableName() string {
	return UserTableName
}

func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		Account:      u.Account,
		Name:         u.Name,
		Email:        u.Email,
		IsProhibited: u.IsProhibited,
		AvatarUrl:    u.AvatarUrl,
		IsAdmin:      u.IsAdmin,
	}
}
