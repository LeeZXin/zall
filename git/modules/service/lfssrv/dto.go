package lfssrv

import (
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/util"
	"github.com/gin-gonic/gin"
	"io"
	"regexp"
	"time"
)

var (
	oidPattern = regexp.MustCompile(`^[a-f\d]{64}$`)
)

type LfsLockDTO struct {
	LockId int64  `json:"lockId"`
	RepoId int64  `json:"repoId"`
	Owner  string `json:"owner"`
	Path   string `json:"path"`
	// 分支名称
	RefName string    `json:"refName"`
	Created time.Time `json:"created"`
}

type LockReqDTO struct {
	RefName  string          `json:"refName"`
	Repo     repomd.RepoInfo `json:"repo"`
	Operator usermd.UserInfo `json:"operator"`
	Path     string          `json:"path"`
}

func (r *LockReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Path == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type LockRespDTO struct {
	AlreadyExists bool
	Lock          LfsLockDTO
}

type ListLockReqDTO struct {
	Repo            repomd.RepoInfo `json:"repo"`
	Operator        usermd.UserInfo `json:"operator"`
	Path            string          `json:"path"`
	Cursor          int64           `json:"cursor"`
	Limit           int             `json:"limit"`
	RefName         string          `json:"refName"`
	FromAccessToken bool            `json:"fromAccessToken"`
}

func (r *ListLockReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListLockRespDTO struct {
	LockList []LfsLockDTO `json:"lockList"`
	Next     int64        `json:"next"`
}

type UnlockReqDTO struct {
	Repo     repomd.RepoInfo `json:"repo"`
	LockId   int64           `json:"lockId"`
	Force    bool            `json:"force"`
	Operator usermd.UserInfo `json:"operator"`
}

func (r *UnlockReqDTO) IsValid() error {
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PointerDTO struct {
	Oid  string `json:"oid"`
	Size int64  `json:"size"`
}

func (p *PointerDTO) IsValid() error {
	if !oidPattern.MatchString(p.Oid) {
		return util.InvalidArgsError()
	}
	if p.Size < 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type VerifyReqDTO struct {
	PointerDTO      `json:"pointerDTO"`
	Repo            repomd.RepoInfo `json:"repo"`
	Operator        usermd.UserInfo `json:"operator"`
	FromAccessToken bool            `json:"fromAccessToken"`
}

func (r *VerifyReqDTO) IsValid() error {
	if err := r.PointerDTO.IsValid(); err != nil {
		return util.InvalidArgsError()
	}
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DownloadReqDTO struct {
	Oid             string          `json:"oid"`
	Repo            repomd.RepoInfo `json:"repo"`
	Operator        usermd.UserInfo `json:"operator"`
	FromAccessToken bool            `json:"fromAccessToken"`
	C               *gin.Context    `json:"-"`
}

func (r *DownloadReqDTO) IsValid() error {
	if !oidPattern.MatchString(r.Oid) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	return nil
}

type UploadReqDTO struct {
	Oid      string          `json:"oid"`
	Size     int64           `json:"size"`
	Repo     repomd.RepoInfo `json:"repo"`
	Operator usermd.UserInfo `json:"operator"`
	C        *gin.Context    `json:"-"`
}

func (r *UploadReqDTO) IsValid() error {
	if !oidPattern.MatchString(r.Oid) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	return nil
}

type DownloadRespDTO struct {
	io.ReadCloser
	FromByte int64 `json:"fromByte"`
	ToByte   int64 `json:"toByte"`
	Length   int64 `json:"length"`
}

type BatchReqDTO struct {
	Repo     repomd.RepoInfo `json:"repo"`
	Operator usermd.UserInfo `json:"operator"`
	Objects  []PointerDTO    `json:"objects"`
	IsUpload bool            `json:"isUpload"`
	RefName  string          `json:"refName"`
}

func (r *BatchReqDTO) IsValid() error {
	// too long
	if len(r.Objects) > 1000 {
		return util.InvalidArgsError()
	}
	for _, obj := range r.Objects {
		if err := obj.IsValid(); err != nil {
			return err
		}
	}
	if !validateRepo(r.Repo) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}

	return nil
}

type LinkDTO struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header"`
	ExpiresAt *time.Time        `json:"expiresAt"`
}

// ObjectErrDTO defines the JSON structure returned to the client in case of an error.
type ObjectErrDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ObjectDTO struct {
	PointerDTO
	ErrObjDTO
}

type BatchRespDTO struct {
	ObjectList []ObjectDTO `json:"objectList"`
}

type ErrObjDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func validateRepo(repo repomd.RepoInfo) bool {
	return repo.RepoId > 0
}
