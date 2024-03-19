package gitnodemd

type InsertNodeReqDTO struct {
	Name     string
	HttpHost string
	SshHost  string
}

type UpdateNodeReqDTO struct {
	Id       int64
	Name     string
	HttpHost string
	SshHost  string
}
