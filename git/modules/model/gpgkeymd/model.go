package gpgkeymd

import "time"

const (
	GpgKeyTableName = "zgit_gpg_key"
)

type GpgKey struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	Account   string    `json:"account"`
	Name      string    `json:"name"`
	PubKeyId  string    `json:"pubKeyId"`
	Content   string    `json:"content"`
	Expiry    int64     `json:"expiry"`
	EmailList string    `json:"emailList"`
	Created   time.Time `json:"created" xorm:"created"`
}

func (*GpgKey) TableName() string {
	return GpgKeyTableName
}

type SimpleGpgKey struct {
	Id        int64  `json:"id" xorm:"pk autoincr"`
	Name      string `json:"name"`
	PubKeyId  string `json:"pubKeyId"`
	Expiry    int64  `json:"expiry"`
	EmailList string `json:"emailList"`
}

func (*SimpleGpgKey) TableName() string {
	return GpgKeyTableName
}
