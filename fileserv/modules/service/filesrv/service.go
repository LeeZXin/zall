package filesrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	UploadAvatar(context.Context, UploadAvatarReqDTO) (string, error)
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
