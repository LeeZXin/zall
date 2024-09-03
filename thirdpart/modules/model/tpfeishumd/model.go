package tpfeishumd

import "time"

const (
	AccessTokenTableName = "zall_feishu_access_token"
)

type AccessToken struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	TeamId      int64     `json:"teamId"`
	Name        string    `json:"name"`
	AppId       string    `json:"appId"`
	Secret      string    `json:"secret"`
	Token       string    `json:"token"`
	TenantToken string    `json:"tenantToken"`
	ExpireTime  int64     `json:"expireTime"`
	ApiKey      string    `json:"apiKey"`
	Creator     string    `json:"creator"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (*AccessToken) TableName() string {
	return AccessTokenTableName
}

func (t *AccessToken) IsNotExpired() bool {
	return t.ExpireTime >= time.Now().UnixMilli()
}
