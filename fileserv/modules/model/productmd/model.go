package productmd

import "time"

const (
	ProductTableName = "zfile_product"
)

type Product struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Name    string    `json:"name"`
	Creator string    `json:"creator"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*Product) TableName() string {
	return ProductTableName
}
