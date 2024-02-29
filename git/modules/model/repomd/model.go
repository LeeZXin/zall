package repomd

import (
	"encoding/json"
	"time"
)

const (
	RepoTableName        = "repo"
	AccessTokenTableName = "repo_access_token"
	ActionTableName      = "repo_actions"
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
	IsEmpty       bool      `json:"isEmpty"`
	GitSize       int64     `json:"gitSize"`
	LfsSize       int64     `json:"lfsSize"`
	Cfg           string    `json:"cfg"`
	NodeId        string    `json:"nodeId"`
	Created       time.Time `json:"created" xorm:"created"`
	Updated       time.Time `json:"updated" xorm:"updated"`
}

func (*Repo) TableName() string {
	return RepoTableName
}

type AccessToken struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	RepoId  int64     `json:"repoId"`
	Account string    `json:"account"`
	Token   string    `json:"token"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*AccessToken) TableName() string {
	return AccessTokenTableName
}

func (r *Repo) GetCfg() RepoCfg {
	var ret RepoCfg
	_ = json.Unmarshal([]byte(r.Cfg), &ret)
	return ret
}

func (r *Repo) ToRepoInfo() RepoInfo {
	return RepoInfo{
		RepoId:  r.Id,
		Name:    r.Name,
		Path:    r.Path,
		Author:  r.Author,
		TeamId:  r.TeamId,
		IsEmpty: r.IsEmpty,
		GitSize: r.GitSize,
		LfsSize: r.LfsSize,
		CfgStr:  r.Cfg,
		Cfg:     r.GetCfg(),
		NodeId:  r.NodeId,
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

type Action struct {
	Id             int64     `json:"id" xorm:"pk autoincr"`
	RepoId         int64     `json:"repoId"`
	Content        string    `json:"content"`
	AssignInstance string    `json:"assignInstance"`
	Created        time.Time `json:"created" xorm:"created"`
	Updated        time.Time `json:"updated" xorm:"updated"`
}

func (*Action) TableName() string {
	return ActionTableName
}
