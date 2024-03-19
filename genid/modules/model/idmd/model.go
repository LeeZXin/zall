package idmd

import "time"

const (
	IdTableName = "zid_generator"
)

type Generator struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	BizName   string    `json:"bizName"`
	CurrentId int64     `json:"currentId"`
	Version   int64     `json:"version"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (*Generator) TableName() string {
	return IdTableName
}
