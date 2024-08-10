package sshkeyapi

type CreateSshKeyReqVO struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type SshKeyVO struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Fingerprint  string `json:"fingerprint"`
	Created      string `json:"created"`
	LastOperated string `json:"lastOperated"`
}

type VerifyTokenReqVO struct {
	Id        string `json:"id"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}
