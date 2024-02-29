package opmd

import "time"

const (
	OpTableName = "op_log"
)

type OpLog struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	Operator   string    `json:"operator"`
	OpDesc     string    `json:"opDesc"`
	ReqContent string    `json:"reqContent"`
	ErrMsg     string    `json:"errMsg"`
	Created    time.Time `json:"created" xorm:"created"`
}

func (*OpLog) TableName() string {
	return OpTableName
}
