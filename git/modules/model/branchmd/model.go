package branchmd

import (
	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/LeeZXin/zall/pkg/branch"
	"time"
)

const (
	ProtectedBranchTableName = "zgit_protected_branch"
)

type ProtectedBranchList []ProtectedBranch

func (l ProtectedBranchList) IsProtectedBranch(branch string) (bool, ProtectedBranch) {
	for _, b := range l {
		if wildcard.Match(b.Pattern, branch) {
			return true, b
		}
	}
	return false, ProtectedBranch{}
}

type ProtectedBranch struct {
	Id      int64                      `json:"id" xorm:"pk autoincr"`
	Pattern string                     `json:"pattern"`
	RepoId  int64                      `json:"repoId"`
	Cfg     *branch.ProtectedBranchCfg `json:"cfg"`
	Created time.Time                  `json:"created" xorm:"created"`
	Updated time.Time                  `json:"updated" xorm:"updated"`
}

func (b *ProtectedBranch) GetCfg() branch.ProtectedBranchCfg {
	if b.Cfg == nil {
		return branch.ProtectedBranchCfg{}
	}
	return *b.Cfg
}

func (*ProtectedBranch) TableName() string {
	return ProtectedBranchTableName
}
