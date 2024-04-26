package repomd

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
	RepoStatus    RepoStatus
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
