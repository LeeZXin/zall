package lfsmd

import "time"

const (
	MetaObjectTableName = "zgit_lfs_meta"
	LfsLockTableName    = "zgit_lfs_lock"
)

type MetaObject struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	RepoId  int64     `json:"repoId"`
	Oid     string    `json:"oid"`
	Size    int64     `json:"size"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*MetaObject) TableName() string {
	return MetaObjectTableName
}

type LfsLock struct {
	Id     int64  `json:"id" xorm:"pk autoincr"`
	RepoId int64  `json:"repoId"`
	Owner  string `json:"owner"`
	Path   string `json:"path" xorm:"TEXT"`
	// 分支名称
	RefName string    `json:"refName"`
	Created time.Time `json:"created" xorm:"created"`
}

func (l LfsLock) TableName() string {
	return LfsLockTableName
}
