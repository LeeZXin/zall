package productsrv

import "context"

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	ListProduct(context.Context, ListProductReqDTO) ([]ProductDTO, error)
}
