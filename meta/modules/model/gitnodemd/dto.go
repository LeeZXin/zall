package gitnodemd

type InsertNodeReqDTO struct {
	NodeId    string
	HttpHosts []string
	SshHosts  []string
}

type UpdateNodeReqDTO struct {
	NodeId    string
	HttpHosts []string
	SshHosts  []string
}

type GitNodeDTO struct {
	NodeId    string
	HttpHosts []string
	SshHosts  []string
	Version   int64
}
