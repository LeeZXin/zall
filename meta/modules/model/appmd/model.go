package appmd

import "time"

const (
	AppTableName = "zall_app"
)

type App struct {
	Id      int64     `xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	TeamId  int64     `json:"teamId"`
	Name    string    `json:"name"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*App) TableName() string {
	return AppTableName
}
