package sshkeymd

import (
	"time"
)

const (
	SshKeyTableName = "zgit_ssh_key"
)

type KeyInfo struct {
	Id          int64  `json:"id"`
	Account     string `json:"account"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Content     string `json:"content"`
}

type SshKey struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	Account     string    `json:"account"`
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	Content     string    `json:"content"`
	Verified    bool      `json:"verified"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (k *SshKey) ToKeyInfo() KeyInfo {
	return KeyInfo{
		Id:          k.Id,
		Account:     k.Account,
		Name:        k.Name,
		Fingerprint: k.Fingerprint,
		Content:     k.Content,
	}
}

func (*SshKey) TableName() string {
	return SshKeyTableName
}
