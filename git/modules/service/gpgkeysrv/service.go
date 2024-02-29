package gpgkeysrv

import (
	"context"
	"github.com/LeeZXin/zall/util"
	"github.com/keybase/go-crypto/openpgp"
)

// echo "3837db075e638d9b9a6e2b585322c8c34b7bca6f7a98715f448a3c44da97ba48" | gpg -a --default-key C18172985BDF2F4A --detach-sig

var (
	Outer OuterService = &outerImpl{
		tokenCache: util.NewGoCache(),
	}

	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	GetVerifiedByAccount(context.Context, string) ([]openpgp.EntityList, error)
}

type OuterService interface {
	InsertGpgKey(context.Context, InsertGpgKeyReqDTO) error
	DeleteGpgKey(context.Context, DeleteGpgKeyReqDTO) error
	ListGpgKey(context.Context, ListGpgKeyReqDTO) ([]GpgKeyDTO, error)
	GetToken(context.Context, GetTokenReqDTO) (string, []string, error)
}
