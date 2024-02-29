package branchmd

type InsertProtectedBranchReqDTO struct {
	RepoId int64
	Branch string
	Cfg    ProtectedBranchCfg
}

type ProtectedBranchDTO struct {
	Id     int64
	RepoId int64
	Branch string
	Cfg    ProtectedBranchCfg
}
