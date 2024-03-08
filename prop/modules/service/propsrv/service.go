package propsrv

import "context"

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	ListSimpleEtcdNode(context.Context) ([]string, error)
	ListEtcdNode(context.Context, ListEtcdNodeReqDTO) ([]EtcdNodeDTO, error)
	InsertEtcdNode(context.Context, InsertEtcdNodeReqDTO) error
	DeleteEtcdNode(context.Context, DeleteEtcdNodeReqDTO) error
	UpdateEtcdNode(context.Context, UpdateEtcdNodeReqDTO) error

	GrantAuth(context.Context, GrantAuthReqDTO) error
	InsertPropContent(context.Context, InsertPropContentReqDTO) error
	UpdatePropContent(context.Context, UpdatePropContentReqDTO) error
	DeletePropContent(context.Context, DeletePropContentReqDTO) error
	ListPropContent(context.Context, ListPropContentReqDTO) ([]PropContentDTO, error)
	DeployPropContent(context.Context, DeployPropContentReqDTO) error

	ListHistory(context.Context, ListHistoryReqDTO) ([]HistoryDTO, int64, error)
	ListDeploy(context.Context, ListDeployReqDTO) ([]DeployDTO, int64, error)
}
