package repomd

import "github.com/LeeZXin/zall/pkg/i18n"

type InsertRepoReqDTO struct {
	Name          string
	Path          string
	Author        string
	TeamId        int64
	RepoDesc      string
	DefaultBranch string
	IsEmpty       bool
	GitSize       int64
	LfsSize       int64
	Cfg           RepoCfg
	NodeId        string
}

type RepoInfo struct {
	RepoId  int64   `json:"repoId"`
	Name    string  `json:"name"`
	Path    string  `json:"path"`
	Author  string  `json:"author"`
	TeamId  int64   `json:"teamId"`
	IsEmpty bool    `json:"isEmpty"`
	GitSize int64   `json:"gitSize"`
	LfsSize int64   `json:"lfsSize"`
	CfgStr  string  `json:"-"`
	Cfg     RepoCfg `json:"cfg"`
	NodeId  string  `json:"nodeId"`
}

type RepoStatus int

const (
	OpenRepoStatus RepoStatus = iota + 1
	ClosedRepoStatus
	DeletedRepoStatus
)

func (s RepoStatus) Int() int {
	return int(s)
}

func (s RepoStatus) Readable() string {
	switch s {
	case OpenRepoStatus:
		return i18n.GetByKey(i18n.RepoOpenStatus)
	case ClosedRepoStatus:
		return i18n.GetByKey(i18n.RepoClosedStatus)
	case DeletedRepoStatus:
		return i18n.GetByKey(i18n.RepoDeletedStatus)
	default:
		return i18n.GetByKey(i18n.RepoUnknownStatus)
	}
}

type InsertAccessTokenReqDTO struct {
	RepoId  int64
	Account string
	Token   string
}

type GetAccessTokenReqDTO struct {
	RepoId  int64
	Account string
}

type InsertActionReqDTO struct {
	RepoId         int64
	Content        string
	AssignInstance string
}

type UpdateActionReqDTO struct {
	ActionId       int64
	Content        string
	AssignInstance string
}
