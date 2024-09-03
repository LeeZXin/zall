package tpfeishumd

type InsertAccessTokenReqDTO struct {
	TeamId      int64
	Name        string
	AppId       string
	Secret      string
	Token       string
	ExpireTime  int64
	ApiKey      string
	TenantToken string
	Creator     string
}

type UpdateAccessTokenTokenReqDTO struct {
	Id          int64
	Token       string
	TenantToken string
	ExpireTime  int64
}

type UpdateAccessTokenReqDTO struct {
	Id     int64
	Name   string
	AppId  string
	Secret string
}

type ListAccessTokenReqDTO struct {
	PageNum  int
	PageSize int
	Key      string
	TeamId   int64
}
