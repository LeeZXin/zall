package sshkeyapi

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

type GetTokenReqVO struct {
	KeyId string `json:"keyId"`
}

type TokenVO struct {
	Token  string   `json:"token"`
	Guides []string `json:"guides"`
}

type VerifyTokenReqVO struct {
	Id        string `json:"id"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}
