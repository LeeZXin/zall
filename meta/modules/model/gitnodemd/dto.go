package gitnodemd

type InsertNodeReqDTO struct {
	NodeId    string
	HttpHosts []string
	SshHosts  []string
}

type GitNodeDTO struct {
	NodeId    string   `json:"nodeId"`
	HttpHosts []string `json:"httpHosts"`
	SshHosts  []string `json:"sshHosts"`
}
