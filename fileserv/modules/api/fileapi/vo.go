package fileapi

type ArtifactVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
	Created string `json:"created"`
}

type ListArtifactReqVO struct {
	AppId   string `json:"appId"`
	Env     string `json:"env"`
	PageNum int    `json:"pageNum"`
}
