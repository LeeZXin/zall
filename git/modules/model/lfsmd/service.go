package lfsmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertLock(ctx context.Context, reqDTO InsertLockReqDTO) (LfsLock, error) {
	ret := LfsLock{
		RepoId:  reqDTO.RepoId,
		Owner:   reqDTO.Owner,
		Path:    reqDTO.Path,
		RefName: reqDTO.RefName,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func GetLockById(ctx context.Context, id int64) (LfsLock, bool, error) {
	var ret LfsLock
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func DeleteLock(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(LfsLock))
	return rows == 1, err
}

func GetMetaObjectByOid(ctx context.Context, oid string, repoId int64) (MetaObject, bool, error) {
	var ret MetaObject
	b, err := xormutil.MustGetXormSession(ctx).
		Where("oid = ?", oid).
		And("repo_id = ?", repoId).
		Get(&ret)
	return ret, b, err
}

func BatchMetaObjectByOidList(ctx context.Context, oidList []string, repoId int64) ([]MetaObject, error) {
	ret := make([]MetaObject, 0, len(oidList))
	err := xormutil.MustGetXormSession(ctx).
		In("oid", oidList).
		And("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func SumSizeWithoutOidList(ctx context.Context, oidList []string, repoId int64) (float64, error) {
	return xormutil.MustGetXormSession(ctx).
		NotIn("oid", oidList).
		And("repo_id = ?", repoId).
		Sum(new(MetaObject), "size")
}

func InsertMetaObject(ctx context.Context, reqDTO InsertMetaObjectReqDTO) (MetaObject, error) {
	ret := MetaObject{
		RepoId: reqDTO.RepoId,
		Oid:    reqDTO.Oid,
		Size:   reqDTO.Size,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func BatchInsertMetaObject(ctx context.Context, reqDTO []InsertMetaObjectReqDTO) ([]MetaObject, error) {
	ret := listutil.MapNe(reqDTO, func(t InsertMetaObjectReqDTO) MetaObject {
		return MetaObject{
			RepoId: t.RepoId,
			Oid:    t.Oid,
			Size:   t.Size,
		}
	})
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func ListLfsLock(ctx context.Context, reqDTO ListLockReqDTO) ([]LfsLock, error) {
	ret := make([]LfsLock, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", reqDTO.RepoId).
		OrderBy("id asc")
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	if reqDTO.Path != "" {
		session.And("path = ?", reqDTO.Path)
	}
	if reqDTO.RefName != "" {
		session.And("ref_name = ?", reqDTO.RefName)
	}
	err := session.Find(&ret)
	return ret, err
}
