package actionsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/actionmd"
)

var (
	Outer OuterService = new(outerImpl)
	Inner InnerService = new(innerImpl)
)

type InnerService interface {
	ExecuteAction(string, string, actionmd.TriggerType)
}

type OuterService interface {
	InsertAction(context.Context, InsertActionReqDTO) error
	UpdateAction(context.Context, UpdateActionReqDTO) error
	DeleteAction(context.Context, DeleteActionReqDTO) error
	ListAction(context.Context, ListActionReqDTO) ([]actionmd.Action, error)
	TriggerAction(context.Context, TriggerActionReqDTO) error
	InsertNode(context.Context, InsertNodeReqDTO) error
	DeleteNode(context.Context, DeleteNodeReqDTO) error
	UpdateNode(context.Context, UpdateNodeReqDTO) error
	ListNode(context.Context, ListNodeReqDTO) ([]NodeDTO, error)
	AllNode(context.Context, AllNodeReqDTO) ([]SimpleNodeDTO, error)
	ListTask(context.Context, ListTaskReqDTO) ([]TaskDTO, int64, error)
	ListStep(context.Context, ListStepReqDTO) ([]StepDTO, error)
}
