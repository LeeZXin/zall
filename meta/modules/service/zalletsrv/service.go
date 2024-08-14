package zalletsrv

import "context"

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// CreateZalletNode 创建node
	CreateZalletNode(context.Context, CreateZalletNodeReqDTO) error
	// UpdateZalletNode 编辑node
	UpdateZalletNode(context.Context, UpdateZalletNodeReqDTO) error
	// DeleteZalletNode 删除node
	DeleteZalletNode(context.Context, DeleteZalletNodeReqDTO) error
	// ListZalletNode node列表
	ListZalletNode(context.Context, ListZalletNodeReqDTO) ([]ZalletNodeDTO, int64, error)
	// ListAllZalletNode 所有列表
	ListAllZalletNode(context.Context, ListAllZalletNodeReqDTO) ([]SimpleZalletNodeDTO, error)
}
