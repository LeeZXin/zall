package fileapi

type ProductVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
	Created string `json:"created"`
}

type ListProductReqVO struct {
	AppId   string `json:"appId"`
	Env     string `json:"env"`
	PageNum int    `json:"pageNum"`
}
