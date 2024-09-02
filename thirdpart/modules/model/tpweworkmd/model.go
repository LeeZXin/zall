package tpweworkmd

import "time"

const (
	AccessTokenTableName  = "zall_wework_access_token"
	CollaboratorTableName = "zall_wework_access_token_collaborator"
)

type AccessToken struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	Name       string    `json:"name"`
	CorpId     string    `json:"corpId"`
	Secret     string    `json:"secret"`
	Token      string    `json:"token"`
	ExpireTime int64     `json:"expireTime"`
	ApiKey     string    `json:"apiKey"`
	Created    time.Time `json:"created" xorm:"created"`
	Updated    time.Time `json:"updated" xorm:"updated"`
}

func (*AccessToken) TableName() string {
	return AccessTokenTableName
}

func (t *AccessToken) IsNotExpired() bool {
	return t.ExpireTime >= time.Now().UnixMilli()
}

type Collaborator struct {
	Id        int64     `json:"id" xorm:"pk autoincr"`
	TokenId   int64     `json:"tokenId"`
	Account   string    `json:"account"`
	IsCreator bool      `json:"isCreator"`
	Created   time.Time `json:"created" xorm:"created"`
	Updated   time.Time `json:"updated" xorm:"updated"`
}

func (*Collaborator) TableName() string {
	return CollaboratorTableName
}
