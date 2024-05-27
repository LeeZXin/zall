package repomd

import "time"

type InsertRepoReqDTO struct {
	Name          string
	Path          string
	TeamId        int64
	RepoDesc      string
	DefaultBranch string
	GitSize       int64
	LfsSize       int64
	Cfg           RepoCfg
	LastOperated  time.Time
}

type UpdateRepoReqDTO struct {
	Id       int64
	RepoDesc string
	Cfg      RepoCfg
}
