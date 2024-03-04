package idapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type IdsRespVO struct {
	ginutil.BaseResp
	Data []int64 `json:"data"`
}

type InsertGeneratorReqVO struct {
	CurrentId int64  `json:"currentId"`
	BizName   string `json:"bizName"`
}
