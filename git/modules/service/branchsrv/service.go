package branchsrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	InsertProtectedBranch(context.Context, InsertProtectedBranchReqDTO) error
	DeleteProtectedBranch(context.Context, DeleteProtectedBranchReqDTO) error
	ListProtectedBranch(context.Context, ListProtectedBranchReqDTO) ([]ProtectedBranchDTO, error)
}
