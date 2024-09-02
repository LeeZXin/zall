package tpfeishumd

type InsertAccessTokenReqDTO struct {
	AppId       string
	Secret      string
	Token       string
	ExpireTime  int64
	ApiKey      string
	TenantToken string
}

type UpdateAccessTokenTokenReqDTO struct {
	Id          int64
	Token       string
	TenantToken string
	ExpireTime  int64
}

type InsertCollaboratorReqDTO struct {
	TokenId   int64
	Account   string
	IsCreator bool
}

type ListCollaboratorReqDTO struct {
	PageNum  int
	PageSize int
	Account  string
}
