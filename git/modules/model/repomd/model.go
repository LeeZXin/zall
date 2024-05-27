package repomd

import (
	"encoding/json"
	"time"
)

const (
	RepoTableName = "zgit_repo"
)

type Repo struct {
	Id            int64  `json:"id" xorm:"pk autoincr"`
	Path          string `json:"path"`
	Name          string `json:"name"`
	TeamId        int64  `json:"teamId"`
	RepoDesc      string `json:"repoDesc"`
	DefaultBranch string `json:"defaultBranch"`
	GitSize       int64  `json:"gitSize"`
	LfsSize       int64  `json:"lfsSize"`
	// 仓库配置
	Cfg *RepoCfg `json:"cfg"`
	// LastOperated 最后操作时间
	LastOperated time.Time `json:"lastOperated"`
	// IsArchived 是否归档
	IsArchived bool `json:"isArchived"`
	// IsDeleted 是否删除
	IsDeleted bool      `json:"isDeleted"`
	Deleted   time.Time `json:"deleted"`
	// Created Updated 创建时间和更新时间
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

func (*Repo) TableName() string {
	return RepoTableName
}

func (r *Repo) GetCfg() RepoCfg {
	if r.Cfg != nil {
		return *r.Cfg
	}
	return RepoCfg{}
}

type RepoCfg struct {
	// DisableLfs 是否禁用lfs
	DisableLfs bool `json:"disableLfs"`
	// LfsLimitSize lfs大小限制
	LfsLimitSize int64 `json:"lfsLimitSize"`
	// GitLimitSize 仓库大小限制
	GitLimitSize int64 `json:"gitLimitSize"`
}

func (c *RepoCfg) FromDB(content []byte) error {
	if c == nil {
		*c = RepoCfg{}
	}
	return json.Unmarshal(content, c)
}

func (c *RepoCfg) ToDB() ([]byte, error) {
	return json.Marshal(c)
}
