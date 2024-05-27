package oplogmd

import "time"

const OpLogTableName = "zgit_op_log"

type OpLog struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	Operator string    `json:"operator"`
	RepoId   int64     `json:"repoId"`
	Content  string    `json:"content"`
	ReqBody  string    `json:"reqBody"`
	Created  time.Time `json:"created"`
}

func (*OpLog) TableName() string {
	return OpLogTableName
}
