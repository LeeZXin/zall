package sshkeysrv

import (
	"context"
)

var (
	Inner InnerService
	Outer OuterService
)

func InitInner() {
	if Inner == nil {
		Inner = new(innerImpl)
	}
}

func InitOuter() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

func Init() {
	InitInner()
	InitOuter()
}

type InnerService interface {
	ListAllPubKeyByAccount(context.Context, string) ([]string, error)
}

type OuterService interface {
	// DeleteSshKey 删除ssh密钥
	DeleteSshKey(context.Context, DeleteSshKeyReqDTO) error
	// CreateSshKey 创建ssh密钥
	CreateSshKey(context.Context, CreateSshKeyReqDTO) error
	// ListSshKey 展示ssh密钥
	ListSshKey(context.Context, ListSshKeyReqDTO) ([]SshKeyDTO, error)
}
