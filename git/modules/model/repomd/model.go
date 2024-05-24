package repomd

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	RepoTableName = "zgit_repo"
)

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

type Repo struct {
	Id            int64      `json:"id" xorm:"pk autoincr"`
	Path          string     `json:"path"`
	Name          string     `json:"name"`
	Author        string     `json:"author"`
	TeamId        int64      `json:"teamId"`
	RepoDesc      string     `json:"repoDesc"`
	DefaultBranch string     `json:"defaultBranch"`
	RepoStatus    RepoStatus `json:"repoStatus"`
	GitSize       int64      `json:"gitSize"`
	LfsSize       int64      `json:"lfsSize"`
	// 仓库配置
	Cfg *RepoCfg `json:"cfg"`
	// LastOperated 最后操作时间
	LastOperated time.Time `json:"lastOperated"`
	Created      time.Time `json:"created" xorm:"created"`
	Updated      time.Time `json:"updated" xorm:"updated"`
}

func (*Repo) TableName() string {
	return RepoTableName
}

type RepoCfg struct {
	// DisableLfs 是否禁用lfs
	DisableLfs bool `json:"disableLfs"`
	// 单个lfs size大小限制
	SingleLfsFileLimitSize int64 `json:"singleLfsFileLimitSize"`
	// 整个lfs仓库大小限制
	MaxLfsLimitSize int64 `json:"maxLfsLimitSize"`
	// 整个仓库大小限制
	MaxGitLimitSize int64 `json:"maxGitLimitSize"`
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
