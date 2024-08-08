package sshkeymd

import (
	"time"
)

const (
	SshKeyTableName = "zgit_ssh_key"
)

type SshKey struct {
	Id           int64     `json:"id" xorm:"pk autoincr"`
	Account      string    `json:"account"`
	Name         string    `json:"name"`
	Fingerprint  string    `json:"fingerprint"`
	Content      string    `json:"content"`
	LastOperated time.Time `json:"lastOperated" xorm:"created"`
	Created      time.Time `json:"created" xorm:"created"`
}

func (*SshKey) TableName() string {
	return SshKeyTableName
}
