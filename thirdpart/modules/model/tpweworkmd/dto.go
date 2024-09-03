package tpweworkmd

type InsertAccessTokenReqDTO struct {
	Name       string
	TeamId     int64
	CorpId     string
	Secret     string
	Token      string
	ExpireTime int64
	ApiKey     string
	Creator    string
}

type UpdateAccessTokenTokenReqDTO struct {
	Id         int64
	Token      string
	ExpireTime int64
}

type ListAccessTokenReqDTO struct {
	PageNum  int
	PageSize int
	Key      string
	TeamId   int64
}

type UpdateAccessTokenReqDTO struct {
	Id     int64
	Name   string
	CorpId string
	Secret string
}
