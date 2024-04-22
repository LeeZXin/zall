package repomd

import "github.com/LeeZXin/zall/pkg/i18n"

type InsertRepoReqDTO struct {
	Name          string
	Path          string
	Author        string
	TeamId        int64
	RepoDesc      string
	DefaultBranch string
	GitSize       int64
	LfsSize       int64
	Cfg           RepoCfg
}

type RepoInfo struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Path    string  `json:"path"`
	Author  string  `json:"author"`
	TeamId  int64   `json:"teamId"`
	GitSize int64   `json:"gitSize"`
	LfsSize int64   `json:"lfsSize"`
	CfgStr  string  `json:"-"`
	Cfg     RepoCfg `json:"cfg"`
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

type InsertRepoTokenReqDTO struct {
	RepoId  int64
	Account string
	Token   string
}

type GetRepoTokenReqDTO struct {
	RepoId  int64
	Account string
}

type InsertActionReqDTO struct {
	RepoId     int64
	Name       string
	Content    string
	NodeId     int64
	PushBranch string
}

type UpdateActionReqDTO struct {
	Id         int64
	Name       string
	Content    string
	NodeId     int64
	PushBranch string
}
