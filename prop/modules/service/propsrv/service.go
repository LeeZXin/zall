package propsrv

import (
	"context"
	"github.com/LeeZXin/zsf/http/httptask"
	"net/url"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = new(innerImpl)
)

func init() {
	httptask.AppendHttpTask("checkPropDbEtcdConsistent", func(_ []byte, _ url.Values) {
		Inner.CheckConsistent()
	})
}

type InnerService interface {
	GrantAuth(context.Context, string)
	// CheckConsistent 检查数据库和etcd的数据一致性
	CheckConsistent()
}

type OuterService interface {
	ListSimpleEtcdNode(context.Context) ([]string, error)
	ListEtcdNode(context.Context, ListEtcdNodeReqDTO) ([]EtcdNodeDTO, error)
	InsertEtcdNode(context.Context, InsertEtcdNodeReqDTO) error
	DeleteEtcdNode(context.Context, DeleteEtcdNodeReqDTO) error
	UpdateEtcdNode(context.Context, UpdateEtcdNodeReqDTO) error

	GrantAuth(context.Context, GrantAuthReqDTO) error
	GetAuth(context.Context, GetAuthReqDTO) (string, string, error)
	InsertPropContent(context.Context, InsertPropContentReqDTO) error
	UpdatePropContent(context.Context, UpdatePropContentReqDTO) error
	DeletePropContent(context.Context, DeletePropContentReqDTO) error
	ListPropContent(context.Context, ListPropContentReqDTO) ([]PropContentDTO, error)
	DeployPropContent(context.Context, DeployPropContentReqDTO) error

	ListHistory(context.Context, ListHistoryReqDTO) ([]HistoryDTO, int64, error)
	ListDeploy(context.Context, ListDeployReqDTO) ([]DeployDTO, int64, error)
}
