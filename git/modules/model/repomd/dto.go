package repomd

import "time"

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
	LastOperated  time.Time
}
