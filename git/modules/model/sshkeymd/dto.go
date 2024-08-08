package sshkeymd

import "time"

type InsertSshKeyReqDTO struct {
	Account      string
	Name         string
	Fingerprint  string
	Content      string
	LastOperated time.Time
}

type UpdateVerifiedVarReqDTO struct {
	KeyId    int64
	Verified bool
}
