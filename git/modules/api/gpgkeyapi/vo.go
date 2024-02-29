package gpgkeyapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type GetTokenReqVO struct {
	Content string `json:"content"`
}

type InsertGpgKeyReqVO struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
	Content   string `json:"content"`
}

type DeleteGpgKeyReqVO struct {
	Id int64 `json:"id"`
}

type GetTokenRespVO struct {
	ginutil.BaseResp
	Token  string
	Guides []string
}

type GpgKeyVO struct {
	Id         int64    `json:"id"`
	Name       string   `json:"name"`
	PubKeyId   string   `json:"pubKeyId"`
	ExpireTime string   `json:"expireTime"`
	EmailList  []string `json:"emailList"`
}

type ListGpgKeyRespVO struct {
	ginutil.BaseResp
	Data []GpgKeyVO `json:"data"`
}
