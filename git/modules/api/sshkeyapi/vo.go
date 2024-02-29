package sshkeyapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type InsertSshKeyReqVO struct {
	Name          string `json:"name"`
	PubKeyContent string `json:"pubKeyContent"`
}

type DeleteSshKeyReqVO struct {
	Id int64 `json:"id"`
}

type SshKeyVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
}

type ListSshKeyRespVO struct {
	ginutil.BaseResp
	Data []SshKeyVO `json:"data"`
}

type GetTokenReqVO struct {
	KeyId string `json:"keyId"`
}

type GetTokenRespVO struct {
	ginutil.BaseResp
	Token  string   `json:"token"`
	Guides []string `json:"guides"`
}

type VerifyTokenReqVO struct {
	Id        string `json:"id"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}
