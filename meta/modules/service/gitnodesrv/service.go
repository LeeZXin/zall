package gitnodesrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertNode(context.Context, InsertNodeReqDTO) error
	DeleteNode(context.Context, DeleteNodeReqDTO) error
}
