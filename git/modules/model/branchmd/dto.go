package branchmd

import "github.com/LeeZXin/zall/pkg/branch"

type InsertProtectedBranchReqDTO struct {
	RepoId  int64
	Pattern string
	Cfg     branch.ProtectedBranchCfg
}

type UpdateProtectedBranchReqDTO struct {
	Id      int64
	Pattern string
	Cfg     branch.ProtectedBranchCfg
}
