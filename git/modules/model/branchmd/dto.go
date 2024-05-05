package branchmd

type InsertProtectedBranchReqDTO struct {
	RepoId int64
	Branch string
	Cfg    ProtectedBranchCfg
}
