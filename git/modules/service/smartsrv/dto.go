package smartsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type UploadPackReqDTO struct {
	Repo     repomd.Repo     `json:"repo"`
	Operator usermd.UserInfo `json:"operator"`
	C        *gin.Context    `json:"-"`
}

func (r *UploadPackReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.C == nil {
		return util.InvalidArgsError()
	}
	return nil
}

type ReceivePackReqDTO struct {
	Repo     repomd.Repo     `json:"repo"`
	Operator usermd.UserInfo `json:"operator"`
	C        *gin.Context    `json:"-"`
}

func (r *ReceivePackReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.C == nil {
		return util.InvalidArgsError()
	}
	return nil
}

func validateRepo(info repomd.Repo) bool {
	return info.Id > 0
}

type InfoRefsReqDTO struct {
	Repo     repomd.Repo     `json:"repo"`
	Operator usermd.UserInfo `json:"operator"`
	C        *gin.Context    `json:"-"`
}

func (r *InfoRefsReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.C == nil {
		return util.InvalidArgsError()
	}
	return nil
}

type SendFileReqDTO struct {
	Repo     repomd.Repo         `json:"repo"`
	Operator apisession.UserInfo `json:"operator"`
	FilePath string              `json:"filePath"`
}

func (r *SendFileReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.FilePath == "" || util.ContainsParentDirectorySeparator(r.FilePath) {
		return util.InvalidArgsError()
	}
	return nil
}

type SendFileRespDTO struct {
	File    io.ReadCloser
	Size    int64
	ModTime time.Time
}
