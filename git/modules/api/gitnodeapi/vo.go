package gitnodeapi

type InsertNodeReqVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
	SshHost  string `json:"sshHost"`
}

type UpdateNodeReqVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
	SshHost  string `json:"sshHost"`
}

type DeleteNodeReqVO struct {
	Id int64 `json:"id"`
}

type NodeVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
	SshHost  string `json:"sshHost"`
}
