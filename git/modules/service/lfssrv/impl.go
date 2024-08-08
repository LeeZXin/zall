package lfssrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/lfsmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/oplogsrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"net/http"
	"strings"
)

const (
	accessRepo = iota
	updateRepo
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
	err := checkPerm(ctx, reqDTO.Repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return LockRespDTO{}, err
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
	err := checkPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return ListLockRespDTO{}, err
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
	err := checkPerm(ctx, reqDTO.Repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return LfsLockDTO{}, err
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
	err := checkPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return false, false, err
	}
	stat, err := client.LfsStat(ctx, reqvo.LfsStatReq{
		RepoPath: reqDTO.Repo.Path,
		Oid:      reqDTO.Oid,
	})
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

func (s *outerImpl) Download(ctx context.Context, reqDTO DownloadReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO.Repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	err = client.LfsDownload(reqvo.LfsDownloadReq{
		RepoPath: reqDTO.Repo.Path,
		Oid:      reqDTO.Oid,
		C:        reqDTO.C,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   reqDTO.Repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.LfsSrvKeysVO.Download, reqDTO.Oid),
		Req:      reqDTO,
	})
	return nil
}

func (s *outerImpl) Upload(ctx context.Context, reqDTO UploadReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查仓库访问权限
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO.Repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return err
	}
	/*
		这个接口是可以绕过batch接口绕过lfs限制并直接上传文件的
		所以依然要判断lfs大小限制
	*/
	limitSize := reqDTO.Repo.GetCfg().LfsLimitSize
	if limitSize > 0 {
		var lfsTotalSize float64
		lfsTotalSize, err = lfsmd.SumSizeWithoutOidList(ctx, []string{reqDTO.Oid}, reqDTO.Repo.Id)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if limitSize < int64(lfsTotalSize)+reqDTO.Size {
			return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.RepoSizeExceedLimit, util.VolumeReadable(limitSize))
		}
	}
	// 检查oid是否落库
	_, b, err := lfsmd.GetMetaObjectByOid(ctx, reqDTO.Oid, reqDTO.Repo.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		_, err = lfsmd.InsertMetaObject(ctx, lfsmd.InsertMetaObjectReqDTO{
			RepoId: reqDTO.Repo.Id,
			Oid:    reqDTO.Oid,
			Size:   reqDTO.Size,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
	}
	err = client.LfsUpload(reqvo.LfsUploadReq{
		RepoPath: reqDTO.Repo.Path,
		Oid:      reqDTO.Oid,
		C:        reqDTO.C,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   reqDTO.Repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.LfsSrvKeysVO.Upload, reqDTO.Oid),
		Req:      reqDTO,
	})
	return nil
}

func (s *outerImpl) Batch(ctx context.Context, reqDTO BatchReqDTO) (BatchRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return BatchRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查保护分支
	ref := strings.TrimPrefix(reqDTO.RefName, git.BranchPrefix)
	pbCfg, b, err := branchmd.IsProtectedBranch(ctx, reqDTO.Repo.Id, ref)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return BatchRespDTO{}, util.InternalError(err)
	}
	if b {
		switch pbCfg.PushOption {
		case branch.NotAllowPush:
			return BatchRespDTO{}, util.NewBizErr(apicode.ProtectedBranchNotAllowPushCode, i18n.ProtectedBranchNotAllowPush)
		case branch.WhiteListPush:
			if !pbCfg.PushWhiteList.Contains(reqDTO.Operator.Account) {
				return BatchRespDTO{}, util.NewBizErr(apicode.ProtectedBranchNotAllowPushCode, i18n.ProtectedBranchNotAllowPush)
			}
		}
	}
	ret := make([]ObjectDTO, 0, len(reqDTO.Objects))
	oidList := make([]string, 0, len(reqDTO.Objects))
	var reqTotalSize int64 = 0
	for _, obj := range reqDTO.Objects {
		oidList = append(oidList, obj.Oid)
		reqTotalSize += obj.Size
	}
	// 检查lfs大小限制
	if reqDTO.IsUpload {
		// 超过限制大小
		limitSize := reqDTO.Repo.GetCfg().LfsLimitSize
		if limitSize > 0 {
			// 获取已入库lfs大小
			lfsExistTotalSize, err := lfsmd.SumSizeWithoutOidList(ctx, oidList, reqDTO.Repo.Id)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return BatchRespDTO{}, util.InternalError(err)
			}
			if limitSize < (reqTotalSize + int64(lfsExistTotalSize)) {
				return BatchRespDTO{}, util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.RepoSizeExceedLimit, util.VolumeReadable(limitSize))
			}
		}
	}
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
	})
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
					Code:    http.StatusUnprocessableEntity,
					Message: i18n.GetByKey(i18n.SystemInvalidArgs),
				},
			})
			continue
		}
		// 文件存在 但没有落库
		exists := existsMap[object.Oid]
		if reqDTO.IsUpload {
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
						Code:    http.StatusNotFound,
						Message: i18n.GetByKey(i18n.SystemNotExists),
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

func checkPerm(ctx context.Context, repo repomd.Repo, operator usermd.UserInfo, permCode int) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case accessRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo
	case updateRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanPushRepo
	}
	if !pass {
		return util.UnauthorizedError()
	}
	return nil
}
