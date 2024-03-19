package gitactionsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/gitactionmd"
	"github.com/LeeZXin/zall/pkg/action"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	ExecuteAction(context.Context, action.Hook)
}

type OuterService interface {
	InsertAction(context.Context, InsertActionReqDTO) error
	UpdateAction(context.Context, UpdateActionReqDTO) error
	DeleteAction(context.Context, DeleteActionReqDTO) error
	ListAction(context.Context, ListActionReqDTO) ([]gitactionmd.Action, error)
	TriggerAction(context.Context, TriggerActionReqDTO) error

	InsertNode(context.Context, InsertNodeReqDTO) error
	DeleteNode(context.Context, DeleteNodeReqDTO) error
	UpdateNode(context.Context, UpdateNodeReqDTO) error
	ListNode(context.Context, ListNodeReqDTO) ([]NodeDTO, error)
}
