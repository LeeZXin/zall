package discoverysrv

import (
	"context"
)

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outImpl)
	}
}

type OuterService interface {
	// ListDiscoverySource 注册中心来源
	ListDiscoverySource(context.Context, ListDiscoverySourceReqDTO) ([]DiscoverySourceDTO, error)
	// ListSimpleDiscoverySource 注册中心来源
	ListSimpleDiscoverySource(context.Context, ListDiscoverySourceReqDTO) ([]SimpleDiscoverySourceDTO, error)
	// CreateDiscoverySource 创建注册中心来源
	CreateDiscoverySource(context.Context, CreateDiscoverySourceReqDTO) error
	// DeleteDiscoverySource 删除注册中心来源
	DeleteDiscoverySource(context.Context, DeleteDiscoverySourceReqDTO) error
	// UpdateDiscoverySource 编辑注册中心来源
	UpdateDiscoverySource(context.Context, UpdateDiscoverySourceReqDTO) error
	// ListDiscoveryService 服务列表
	ListDiscoveryService(context.Context, ListDiscoveryServiceReqDTO) ([]ServiceDTO, error)
	// DeregisterService 下线服务
	DeregisterService(context.Context, DeregisterServiceReqDTO) error
	// ReRegisterService 上线服务
	ReRegisterService(context.Context, ReRegisterServiceReqDTO) error
	// DeleteDownService 删除下线服务
	DeleteDownService(context.Context, DeleteDownServiceReqDTO) error
}
