package fileapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type UploadIconRespVO struct {
	ginutil.BaseResp
	Data string `json:"data"`
}

type ProductVO struct {
	Name    string `json:"name"`
	Creator string `json:"creator"`
	Created string `json:"created"`
}
