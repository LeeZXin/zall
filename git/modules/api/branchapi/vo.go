package branchapi

import (
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
)

type InsertProtectedBranchReqVO struct {
	RepoId int64                       `json:"repoId"`
	Branch string                      `json:"branch"`
	Cfg    branchmd.ProtectedBranchCfg `json:"cfg"`
}

type DeleteProtectedBranchReqVO struct {
	Id int64 `json:"id"`
}

type ListProtectedBranchReqVO struct {
	RepoId int64 `json:"repoId"`
}

type ProtectedBranchVO struct {
	Id     int64  `json:"id"`
	Branch string `json:"branch"`
	Cfg    branchmd.ProtectedBranchCfg
}
