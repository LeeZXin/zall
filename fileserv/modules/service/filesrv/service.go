package filesrv

import (
	"context"
)

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// UploadAvatar 上传头像
	UploadAvatar(context.Context, UploadAvatarReqDTO) (string, error)
	// GetAvatar 获取头像路径
	GetAvatar(context.Context, GetAvatarReqDTO) (string, error)
	// UploadProduct 上传制品
	UploadProduct(context.Context, UploadProductReqDTO) (string, error)
	// GetProduct 获取制品路径
	GetProduct(context.Context, GetProductReqDTO) (string, error)
	// ListProduct 制品库列表
	ListProduct(context.Context, ListProductReqDTO) ([]ProductDTO, error)
	// DeleteProduct 删除制品
	DeleteProduct(context.Context, DeleteProductReqDTO) error
}
