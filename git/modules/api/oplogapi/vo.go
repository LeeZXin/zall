package oplogapi

type PageLogReqVO struct {
	RepoId  int64  `json:"repoId"`
	PageNum int    `json:"pageNum"`
	Account string `json:"account"`
	DateStr string `json:"dateStr"`
}

type OpLogVO struct {
	Id      int64  `json:"id"`
	Account string `json:"account"`
	Created string `json:"created"`
	Content string `json:"content"`
	ReqBody string `json:"reqBody"`
}
