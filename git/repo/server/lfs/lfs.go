package lfs

import (
	"context"
	"github.com/LeeZXin/zall/git/repo/reqvo"
)

type Lfs interface {
	Exists(context.Context, reqvo.LfsExistsReq) (bool, error)
	BatchExists(context.Context, reqvo.LfsBatchExistsReq) (map[string]bool, error)
	Upload(context.Context, reqvo.LfsUploadReq)
	Stat(context.Context, reqvo.LfsStatReq) (reqvo.LfsStatResp, error)
	Download(context.Context, reqvo.LfsDownloadReq)
}
