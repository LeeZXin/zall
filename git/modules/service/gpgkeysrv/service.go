package gpgkeysrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/gpgkeymd"
)

// echo "3837db075e638d9b9a6e2b585322c8c34b7bca6f7a98715f448a3c44da97ba48" | gpg -a --default-key C18172985BDF2F4A --detach-sig

var (
	Outer OuterService
	Inner InnerService
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
	// ListValidByAccount 获取有效未过期的
	ListValidByAccount(context.Context, string) []gpgkeymd.GpgKey
	// GetByKeyId 根据keyId获取gpg密钥
	GetByKeyId(context.Context, string) (gpgkeymd.GpgKey, bool)
}

type OuterService interface {
	// CreateGpgKey 创建gpg key
	CreateGpgKey(context.Context, CreateGpgKeyReqDTO) error
	// DeleteGpgKey 删除gpg密钥
	DeleteGpgKey(context.Context, DeleteGpgKeyReqDTO) error
	// ListGpgKey gpg密钥列表
	ListGpgKey(context.Context, ListGpgKeyReqDTO) ([]GpgKeyDTO, error)
}
