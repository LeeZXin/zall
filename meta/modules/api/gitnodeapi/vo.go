package gitnodeapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type InsertNodeReqVO struct {
	NodeId    string   `json:"nodeId"`
	HttpHosts []string `json:"httpHosts"`
	SshHosts  []string `json:"sshHosts"`
}

type UpdateNodeReqVO struct {
	NodeId    string   `json:"nodeId"`
	HttpHosts []string `json:"httpHosts"`
	SshHosts  []string `json:"sshHosts"`
}

type DeleteNodeReqVO struct {
	NodeId string `json:"nodeId"`
}

type GitNodeVO struct {
	NodeId    string   `json:"nodeId"`
	HttpHosts []string `json:"httpHosts"`
	SshHosts  []string `json:"sshHosts"`
}

type ListGitNodeRespVO struct {
	ginutil.BaseResp
	Data []GitNodeVO `json:"data"`
}
