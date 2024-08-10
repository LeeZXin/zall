package gpgkeymd

import (
	"encoding/json"
	"time"
)

const (
	GpgKeyTableName = "zgit_gpg_key"
)

type GpgKey struct {
	Id      int64      `json:"id" xorm:"pk autoincr"`
	Account string     `json:"account"`
	Name    string     `json:"name"`
	KeyId   string     `json:"keyId"`
	SubKeys GpgSubKeys `json:"subKeys"`
	Content string     `json:"content"`
	Email   string     `json:"email"`
	Expired time.Time  `json:"expired"`
	Created time.Time  `json:"created" xorm:"created"`
}

func (*GpgKey) TableName() string {
	return GpgKeyTableName
}

func (k *GpgKey) IsExpired() bool {
	return k.Expired.Before(time.Now())
}

type GpgSubKey struct {
	KeyId   string `json:"keyId"`
	Content string `json:"content"`
}

type GpgSubKeys []GpgSubKey

func (k *GpgSubKeys) FromDB(content []byte) error {
	if k == nil {
		*k = make(GpgSubKeys, 0)
	}
	return json.Unmarshal(content, k)
}

func (k *GpgSubKeys) ToDB() ([]byte, error) {
	return json.Marshal(k)
}
