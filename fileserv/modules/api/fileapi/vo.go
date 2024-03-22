package fileapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type UploadIconRespVO struct {
	ginutil.BaseResp
	Data string `json:"data"`
}
