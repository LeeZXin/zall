package gpgkeymd

import "time"

type InsertGpgKeyReqDTO struct {
	Name    string
	Account string
	KeyId   string
	Content string
	SubKeys []InsertGpgSubKeyReqDTO
	Email   string
	Expired time.Time
}

type InsertGpgSubKeyReqDTO struct {
	KeyId   string
	Content string
}
