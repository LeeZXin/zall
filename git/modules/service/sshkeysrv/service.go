package sshkeysrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/sshkeymd"
	"github.com/LeeZXin/zall/util"
)

var (
	Inner InnerService = &innerImpl{
		sshKeyCache: util.NewGoCache(),
	}
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	GetAccountByFingerprint(context.Context, string) (string, bool, error)
	GetVerifiedByAccount(context.Context, string) ([]string, error)
}

type OuterService interface {
	DeleteSshKey(context.Context, DeleteSshKeyReqDTO) error
	InsertSshKey(context.Context, InsertSshKeyReqDTO) error
	ListSshKey(context.Context, ListSshKeyReqDTO) ([]sshkeymd.SshKey, error)
	GetToken(context.Context, GetTokenReqDTO) (string, []string, error)
}
