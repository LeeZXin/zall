package propertysrv

import (
	"context"
)

var (
	Outer OuterService
)

func Init() {
	Outer = new(outerImpl)
}

type OuterService interface {
	// ListPropertySource 配置来源
	ListPropertySource(context.Context, ListPropertySourceReqDTO) ([]PropertySourceDTO, error)
	// ListPropertySourceByFileId 配置来源
	ListPropertySourceByFileId(context.Context, ListPropertySourceByFileIdReqDTO) ([]SimplePropertySourceDTO, error)
	// CreatePropertySource 创建配置来源
	CreatePropertySource(context.Context, CreatePropertySourceReqDTO) error
	// DeletePropertySource 删除配置来源
	DeletePropertySource(context.Context, DeletePropertySourceReqDTO) error
	// UpdatePropertySource 编辑配置来源
	UpdatePropertySource(context.Context, UpdatePropertySourceReqDTO) error
	// CreateFile 创建配置文件
	CreateFile(context.Context, CreateFileReqDTO) error
	// NewVersion 新增版本
	NewVersion(context.Context, NewVersionReqDTO) error
	// DeleteFile 删除配置文件
	DeleteFile(context.Context, DeleteFileReqDTO) error
	// ListFile 配置文件列表
	ListFile(context.Context, ListFileReqDTO) ([]FileDTO, error)
	// DeployHistory 部署配置
	DeployHistory(context.Context, DeployHistoryReqDTO) error
	// GetHistoryByVersion 获取指定版本的配置
	GetHistoryByVersion(context.Context, GetHistoryByVersionReqDTO) (HistoryDTO, bool, error)
	// PageHistory 版本列表
	PageHistory(context.Context, PageHistoryReqDTO) ([]HistoryDTO, int64, error)
	// ListDeploy 发布记录
	ListDeploy(context.Context, ListDeployReqDTO) ([]DeployDTO, error)
}
