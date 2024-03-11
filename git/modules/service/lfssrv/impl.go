package lfssrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/lfsmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct {
}

func (*outerImpl) Lock(ctx context.Context, reqDTO LockReqDTO) (LockRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return LockRespDTO{}, err
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return LockRespDTO{}, err
	}
	if !p.GetRepoPerm(reqDTO.Repo.Id).CanUpdateRepo {
		return LockRespDTO{}, util.UnauthorizedError()
	}
	lock, err := lfsmd.InsertLock(ctx, lfsmd.InsertLockReqDTO{
		RepoId:  reqDTO.Repo.Id,
		Owner:   reqDTO.Operator.Account,
		Path:    reqDTO.Path,
		RefName: reqDTO.RefName,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return LockRespDTO{}, util.InternalError(err)
	}
	// 添加锁
	return LockRespDTO{
		AlreadyExists: false,
		Lock:          convertDTO(lock),
	}, nil
}

func convertDTO(lock lfsmd.LfsLock) LfsLockDTO {
	return LfsLockDTO{
		LockId:  lock.Id,
		RepoId:  lock.Id,
		Owner:   lock.Owner,
		Path:    lock.Path,
		RefName: lock.RefName,
		Created: lock.Created,
	}
}

func (*outerImpl) ListLock(ctx context.Context, reqDTO ListLockReqDTO) (ListLockRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return ListLockRespDTO{}, err
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.FromAccessToken {
		p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
		if err != nil {
			return ListLockRespDTO{}, err
		}
		if !p.GetRepoPerm(reqDTO.Repo.Id).CanAccessRepo {
			return ListLockRespDTO{}, util.UnauthorizedError()
		}
	}
	if reqDTO.Limit <= 0 || reqDTO.Limit > 1000 {
		reqDTO.Limit = 1000
	}
	locks, err := lfsmd.ListLfsLock(ctx, lfsmd.ListLockReqDTO{
		RepoId:  reqDTO.Repo.Id,
		Path:    reqDTO.Path,
		Cursor:  reqDTO.Cursor,
		Limit:   reqDTO.Limit,
		RefName: reqDTO.RefName,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return ListLockRespDTO{}, util.InternalError(err)
	}
	ret := ListLockRespDTO{}
	if len(locks) == reqDTO.Limit {
		ret.Next = locks[len(locks)-1].Id
	}
	ret.LockList, _ = listutil.Map(locks, func(t lfsmd.LfsLock) (LfsLockDTO, error) {
		return convertDTO(t), nil
	})
	// 查询lock
	return ret, nil
}

func (*outerImpl) Unlock(ctx context.Context, reqDTO UnlockReqDTO) (LfsLockDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return LfsLockDTO{}, err
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return LfsLockDTO{}, err
	}
	if !p.GetRepoPerm(reqDTO.Repo.Id).CanUpdateRepo {
		return LfsLockDTO{}, util.UnauthorizedError()
	}
	// 查找lock是否存在
	lock, b, err := lfsmd.GetLockById(ctx, reqDTO.LockId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return LfsLockDTO{}, util.InternalError(err)
	}
	if !b {
		return LfsLockDTO{}, util.InvalidArgsError()
	}
	if reqDTO.Force || lock.Owner == reqDTO.Operator.Account {
		_, err = lfsmd.DeleteLock(ctx, reqDTO.LockId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return LfsLockDTO{}, util.InternalError(err)
		}
	} else {
		return LfsLockDTO{}, util.UnauthorizedError()
	}
	return convertDTO(lock), nil
}

func (s *outerImpl) Verify(ctx context.Context, reqDTO VerifyReqDTO) (bool, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return false, false, err
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.FromAccessToken {
		p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
		if err != nil {
			return false, false, err
		}
		if !p.GetRepoPerm(reqDTO.Repo.Id).CanAccessRepo {
			return false, false, util.UnauthorizedError()
		}
	}
	stat, err := client.LfsStat(ctx, reqvo.LfsStatReq{
		RepoPath: reqDTO.Repo.Path,
		Oid:      reqDTO.Oid,
	}, reqDTO.Repo.NodeId)
	if err != nil {
		return false, false, util.InternalError(err)
	}
	if !stat.Exists {
		return false, false, nil
	}
	if stat.Size != reqDTO.Size {
		return true, false, util.InvalidArgsError()
	}
	return true, true, nil
}

func (s *outerImpl) Download(ctx context.Context, reqDTO DownloadReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.LfsSrvKeysVO.Download),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.FromAccessToken {
		var p perm.Detail
		p, err = getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
		if err != nil {
			return
		}
		if !p.GetRepoPerm(reqDTO.Repo.Id).CanAccessRepo {
			err = util.UnauthorizedError()
			return
		}
	}
	err = client.LfsDownload(reqvo.LfsDownloadReq{
		RepoPath: reqDTO.Repo.Path,
		Oid:      reqDTO.Oid,
		C:        reqDTO.C,
	}, reqDTO.Repo.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (s *outerImpl) Upload(ctx context.Context, reqDTO UploadReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.LfsSrvKeysVO.Upload),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	p, err := getPerm(ctx, reqDTO.Repo, reqDTO.Operator)
	if err != nil {
		return
	}
	if !p.GetRepoPerm(reqDTO.Repo.Id).CanUpdateRepo {
		err = util.UnauthorizedError()
		return
	}
	// 检查oid是否落库
	_, b, err := lfsmd.GetMetaObjectByOid(ctx, reqDTO.Oid, reqDTO.Repo.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		_, err = lfsmd.InsertMetaObject(ctx, lfsmd.InsertMetaObjectReqDTO{
			RepoId: reqDTO.Repo.Id,
			Oid:    reqDTO.Oid,
			Size:   reqDTO.Size,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	}
	err = client.LfsUpload(reqvo.LfsUploadReq{
		RepoPath: reqDTO.Repo.Path,
		Oid:      reqDTO.Oid,
		C:        reqDTO.C,
	}, reqDTO.Repo.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 忽略lfs size计算异常
	lfsSize, err := lfsmd.SumRepoLfsSize(ctx, reqDTO.Repo.Id)
	if err == nil {
		err = repomd.UpdateLfsSize(ctx, reqDTO.Repo.Id, lfsSize)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
	} else {
		logger.Logger.WithContext(ctx).Error(err)
	}
	return
}

func (s *outerImpl) Batch(ctx context.Context, reqDTO BatchReqDTO) (BatchRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return BatchRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret := make([]ObjectDTO, 0, len(reqDTO.Objects))
	oidList, _ := listutil.Map(reqDTO.Objects, func(t PointerDTO) (string, error) {
		return t.Oid, nil
	})
	objList, err := lfsmd.BatchMetaObjectByOidList(ctx, oidList, reqDTO.Repo.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return BatchRespDTO{}, util.InternalError(err)
	}
	objMap, _ := listutil.CollectToMap(objList, func(t lfsmd.MetaObject) (string, error) {
		return t.Oid, nil
	}, func(t lfsmd.MetaObject) (lfsmd.MetaObject, error) {
		return t, nil
	})
	existsMap, err := client.LfsBatchExists(ctx, reqvo.LfsBatchExistsReq{
		RepoPath: reqDTO.Repo.Path,
		OidList:  oidList,
	}, reqDTO.Repo.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return BatchRespDTO{}, util.InternalError(err)
	}
	shouldInsert := make([]lfsmd.InsertMetaObjectReqDTO, 0)
	for _, object := range reqDTO.Objects {
		meta, b := objMap[object.Oid]
		if b && meta.Size != object.Size {
			// 大小不一致
			ret = append(ret, ObjectDTO{
				ErrObjDTO: ErrObjDTO{
					Code:    422,
					Message: "size is invalid",
				},
			})
			continue
		}
		// 文件存在 但没有落库
		exists := existsMap[object.Oid]
		if reqDTO.IsUpload {
			// 检查是否超过单个lfs文件配置大小
			if !exists && reqDTO.Repo.Cfg.SingleLfsFileLimitSize > 0 && object.Size > reqDTO.Repo.Cfg.SingleLfsFileLimitSize {
				ret = append(ret, ObjectDTO{
					ErrObjDTO: ErrObjDTO{
						Code:    422,
						Message: "size exceeds",
					},
				})
				continue
			}
			if exists && !b {
				shouldInsert = append(shouldInsert, lfsmd.InsertMetaObjectReqDTO{
					RepoId: reqDTO.Repo.Id,
					Oid:    object.Oid,
					Size:   object.Size,
				})
			}
			ret = append(ret, ObjectDTO{
				PointerDTO: object,
			})
		} else {
			if !exists || !b {
				ret = append(ret, ObjectDTO{
					ErrObjDTO: ErrObjDTO{
						Code:    404,
						Message: "not found",
					},
				})
			} else {
				ret = append(ret, ObjectDTO{
					PointerDTO: object,
				})
			}
		}
	}
	if len(shouldInsert) > 0 {
		_, err = lfsmd.BatchInsertMetaObject(ctx, shouldInsert)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return BatchRespDTO{}, util.InternalError(err)
		}
	}
	return BatchRespDTO{
		ObjectList: ret,
	}, nil
}

func getPerm(ctx context.Context, repo repomd.RepoInfo, operator usermd.UserInfo) (perm.Detail, error) {
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return perm.Detail{}, util.UnauthorizedError()
	}
	return p.PermDetail, nil
}
