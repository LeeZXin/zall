package lfs

import (
	"context"
	"encoding/base64"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/lfs"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"net/http"
	"os"
	"path/filepath"
)

func NewLfs() Lfs {
	lfsStore, err := lfs.NewLocalStorage(git.LfsDir(), git.TempDir())
	if err != nil {
		logger.Logger.Fatal(err)
	}
	return &lfsImpl{
		lfsStore: lfsStore,
	}
}

type lfsImpl struct {
	lfsStore lfs.Storage
}

func (s *lfsImpl) Exists(ctx context.Context, req reqvo.LfsExistsReq) (bool, error) {
	b, err := s.lfsStore.Exists(ctx, convertPointerPath(req.RepoPath, req.Oid))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false, util.InternalError(err)
	}
	return b, nil
}

func (s *lfsImpl) BatchExists(ctx context.Context, req reqvo.LfsBatchExistsReq) (map[string]bool, error) {
	ret := make(map[string]bool, len(req.OidList))
	for _, oid := range req.OidList {
		b, err := s.lfsStore.Exists(ctx, convertPointerPath(req.RepoPath, oid))
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		ret[oid] = b
	}
	return ret, nil
}

func (s *lfsImpl) Upload(ctx context.Context, req reqvo.LfsUploadReq) {
	defer req.C.Request.Body.Close()
	_, err := s.lfsStore.Save(ctx, convertPointerPath(req.RepoPath, req.Oid), req.C.Request.Body)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		util.LfsRet(req.C, http.StatusInternalServerError, reqvo.LfsErrVO{
			Message: "internal error",
		})
	}
	util.LfsRet(req.C, http.StatusOK, "")
}

func (s *lfsImpl) Stat(ctx context.Context, req reqvo.LfsStatReq) (reqvo.LfsStatResp, error) {
	stat, err := s.lfsStore.Stat(ctx, convertPointerPath(req.RepoPath, req.Oid))
	if err != nil {
		if os.IsNotExist(err) {
			return reqvo.LfsStatResp{
				Exists: false,
				Size:   0,
			}, nil
		}
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.LfsStatResp{}, util.InternalError(err)
	}
	return reqvo.LfsStatResp{
		Exists: true,
		Size:   stat.Size(),
	}, nil
}

func (s *lfsImpl) Download(ctx context.Context, req reqvo.LfsDownloadReq) {
	path := convertPointerPath(req.RepoPath, req.Oid)
	exists, err := s.lfsStore.Exists(ctx, path)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		util.LfsRet(req.C, http.StatusInternalServerError, reqvo.LfsErrVO{
			Message:   "internal error",
			RequestID: logger.GetTraceId(req.C),
		})
		return
	}
	if !exists {
		util.LfsRet(req.C, http.StatusNotFound, reqvo.LfsErrVO{
			Message:   "not found",
			RequestID: logger.GetTraceId(req.C),
		})
		return
	}
	req.C.Header("Content-Type", "application/octet-stream")
	filename := req.C.Query("filename")
	if filename != "" {
		decodedFilename, err := base64.RawURLEncoding.DecodeString(filename)
		if err == nil {
			req.C.Header("Content-Disposition", "attachment; filename=\""+string(decodedFilename)+"\"")
			req.C.Header("Access-Control-Expose-Headers", "Content-Disposition")
		}
	}
	downloadPath := filepath.Join(s.lfsStore.StoreDir(), path)
	req.C.File(downloadPath)
}

func convertPointerPath(repoPath, oid string) string {
	if len(oid) < 5 {
		return filepath.Join(repoPath, oid)
	}
	return filepath.Join(repoPath, oid[0:2], oid[2:4], oid[4:])
}
