package teammd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/perm"
	"time"
)

const (
	TeamTableName          = "zall_team"
	TeamUserTableName      = "zall_team_user"
	TeamUserGroupTableName = "zall_team_user_group"
)

type Team struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Name    string    `json:"name"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*Team) TableName() string {
	return TeamTableName
}

type TeamUser struct {
	Id      int64  `json:"id" xorm:"pk autoincr"`
	TeamId  int64  `json:"teamId"`
	Account string `json:"account"`
	// 关联用户组
	GroupId int64     `json:"groupId"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*TeamUser) TableName() string {
	return TeamUserTableName
}

type TeamUserGroup struct {
	Id int64 `json:"id" xorm:"pk autoincr"`
	// 项目id
	TeamId int64 `json:"teamId"`
	// 名称
	Name string `json:"name"`
	// 权限json内容
	Perm string `json:"perm"`
	// 是否是管理员用户组
	IsAdmin bool `json:"isAdmin"`
	// 创建时间
	Created time.Time `json:"created" xorm:"created"`
	// 更新时间
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*TeamUserGroup) TableName() string {
	return TeamUserGroupTableName
}

func (p *TeamUserGroup) GetPermDetail() perm.Detail {
	ret := perm.Detail{}
	_ = json.Unmarshal([]byte(p.Perm), &ret)
	return ret
}
