package tpfeishuapi

import "github.com/LeeZXin/zall/util"

type ListWeworkAccessTokenReqVO struct {
	PageNum int    `json:"pageNum"`
	Key     string `json:"key"`
	TeamId  int64  `json:"teamId"`
}

type AccessTokenVO struct {
	Id          int64     `json:"id"`
	TeamId      int64     `json:"teamId"`
	Name        string    `json:"name"`
	AppId       string    `json:"appId"`
	Creator     util.User `json:"creator"`
	Secret      string    `json:"secret"`
	Token       string    `json:"token"`
	TenantToken string    `json:"tenantToken"`
	Expired     string    `json:"expired"`
	ApiKey      string    `json:"apiKey"`
}

type CreateAccessTokenReqVO struct {
	TeamId int64  `json:"teamId"`
	Name   string `json:"name"`
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}

type UpdateAccessTokenReqVO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}
