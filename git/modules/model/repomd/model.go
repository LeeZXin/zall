package repomd

import (
	"encoding/json"
	"time"
)

const (
	RepoTableName      = "zgit_repo"
	RepoTokenTableName = "zgit_repo_token"
)

type Repo struct {
	Id            int64     `json:"id" xorm:"pk autoincr"`
	Path          string    `json:"path"`
	Name          string    `json:"name"`
	Author        string    `json:"author"`
	TeamId        int64     `json:"teamId"`
	RepoDesc      string    `json:"repoDesc"`
	DefaultBranch string    `json:"defaultBranch"`
	RepoStatus    int       `json:"repoStatus"`
	GitSize       int64     `json:"gitSize"`
	LfsSize       int64     `json:"lfsSize"`
	Cfg           string    `json:"cfg"`
	Created       time.Time `json:"created" xorm:"created"`
	Updated       time.Time `json:"updated" xorm:"updated"`
}

func (*Repo) TableName() string {
	return RepoTableName
}

type RepoToken struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	RepoId  int64     `json:"repoId"`
	Account string    `json:"account"`
	Token   string    `json:"token"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*RepoToken) TableName() string {
	return RepoTokenTableName
}

func (r *Repo) GetCfg() RepoCfg {
	var ret RepoCfg
	_ = json.Unmarshal([]byte(r.Cfg), &ret)
	return ret
}

func (r *Repo) ToRepoInfo() RepoInfo {
	return RepoInfo{
		Id:      r.Id,
		Name:    r.Name,
		Path:    r.Path,
		Author:  r.Author,
		TeamId:  r.TeamId,
		GitSize: r.GitSize,
		LfsSize: r.LfsSize,
		CfgStr:  r.Cfg,
		Cfg:     r.GetCfg(),
	}
}

type RepoCfg struct {
	// 单个lfs size大小限制
	SingleLfsFileLimitSize int64 `json:"singleLfsFileLimitSize"`
	// 整个lfs仓库大小限制
	MaxLfsLimitSize int64 `json:"maxLfsLimitSize"`
	// 整个仓库大小限制
	MaxGitLimitSize int64 `json:"maxGitLimitSize"`
}

func (c *RepoCfg) ToString() string {
	m, _ := json.Marshal(c)
	return string(m)
}
