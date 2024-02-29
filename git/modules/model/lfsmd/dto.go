package lfsmd

type InsertLockReqDTO struct {
	RepoId  int64
	RefName string
	Owner   string
	Path    string
}

type InsertMetaObjectReqDTO struct {
	RepoId int64
	Oid    string
	Size   int64
}

type ListLockReqDTO struct {
	RepoId  int64
	Path    string
	Cursor  int64
	Limit   int
	RefName string
}
