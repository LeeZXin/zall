package reqvo

import (
	"github.com/gin-gonic/gin"
)

type LfsUploadReq struct {
	RepoPath string       `json:"repoPath"`
	Oid      string       `json:"oid"`
	C        *gin.Context `json:"-"`
}

type LfsExistsReq struct {
	RepoPath string `json:"repoPath"`
	Oid      string `json:"oid"`
}

type LfsBatchExistsReq struct {
	RepoPath string   `json:"repoPath"`
	OidList  []string `json:"oidList"`
}

type LfsStatReq struct {
	RepoPath string `json:"repoPath"`
	Oid      string `json:"oid"`
}

type LfsStatResp struct {
	Exists bool  `json:"exists"`
	Size   int64 `json:"size"`
}

type LfsDownloadReq struct {
	RepoPath string       `json:"repoPath"`
	Oid      string       `json:"oid"`
	C        *gin.Context `json:"-"`
}

type LfsErrVO struct {
	Message       string `json:"message"`
	Documentation string `json:"documentation_url,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
}
