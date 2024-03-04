package idapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type SnowFlakeRespVO struct {
	ginutil.BaseResp
	Data []int64 `json:"data"`
}
