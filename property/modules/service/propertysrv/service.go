package propertysrv

import (
	"context"
)

var (
	Outer OuterService
	Inner InnerService
)

func Init() {
	Outer = new(outerImpl)
	Inner = new(innerImpl)
}

type InnerService interface {
	GrantAuth(context.Context, string, string)
	// CheckConsistent 检查数据库和etcd的数据一致性
	CheckConsistent(string)
}

type OuterService interface {
	ListSimpleEtcdNode(context.Context, string) ([]string, error)
	ListEtcdNode(context.Context, ListEtcdNodeReqDTO) ([]EtcdNodeDTO, error)
	InsertEtcdNode(context.Context, InsertEtcdNodeReqDTO) error
	DeleteEtcdNode(context.Context, DeleteEtcdNodeReqDTO) error
	UpdateEtcdNode(context.Context, UpdateEtcdNodeReqDTO) error

	GrantAuth(context.Context, GrantAuthReqDTO) error
	GetAuth(context.Context, GetAuthReqDTO) (string, string, error)
	// CreateFile 创建配置文件
	CreateFile(context.Context, CreateFileReqDTO) error
	// NewVersion 新增版本
	NewVersion(context.Context, NewVersionReqDTO) error
	DeletePropContent(context.Context, DeletePropContentReqDTO) error
	// ListFile 配置文件列表
	ListFile(context.Context, ListFileReqDTO) ([]FileDTO, error)
	DeployPropContent(context.Context, DeployPropContentReqDTO) error
	// GetHistoryByVersion 获取指定版本的配置
	GetHistoryByVersion(context.Context, GetHistoryByVersionReqDTO) (HistoryDTO, bool, error)
	// PageHistory 版本列表
	PageHistory(context.Context, PageHistoryReqDTO) ([]HistoryDTO, int64, error)
	ListDeploy(context.Context, ListDeployReqDTO) ([]DeployDTO, int64, error)
}
