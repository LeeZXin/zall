package branchapi

import (
	"github.com/LeeZXin/zall/pkg/branch"
)

type InsertProtectedBranchReqVO struct {
	RepoId  int64                     `json:"repoId"`
	Pattern string                    `json:"pattern"`
	Cfg     branch.ProtectedBranchCfg `json:"cfg"`
}

type ProtectedBranchVO struct {
	Id      int64                     `json:"id"`
	Pattern string                    `json:"pattern"`
	Cfg     branch.ProtectedBranchCfg `json:"cfg"`
}

type UpdateProtectedBranchVO struct {
	ProtectedBranchId int64                     `json:"protectedBranchId"`
	Pattern           string                    `json:"pattern"`
	Cfg               branch.ProtectedBranchCfg `json:"cfg"`
}
