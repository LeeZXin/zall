package sshkeymd

type InsertSshKeyReqDTO struct {
	Account     string
	Name        string
	Fingerprint string
	Content     string
}

type UpdateVerifiedVarReqDTO struct {
	KeyId    int64
	Verified bool
}
