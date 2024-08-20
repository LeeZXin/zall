package branchsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/service/webhooksrv"
)

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		webhooksrv.InitPsub()
		Outer = new(outerImpl)
	}
}

type OuterService interface {
	// CreateProtectedBranch 添加保护分支
	CreateProtectedBranch(context.Context, CreateProtectedBranchReqDTO) error
	// DeleteProtectedBranch 删除保护分支
	DeleteProtectedBranch(context.Context, DeleteProtectedBranchReqDTO) error
	// ListProtectedBranch 保护分支列表
	ListProtectedBranch(context.Context, ListProtectedBranchReqDTO) ([]ProtectedBranchDTO, error)
	// UpdateProtectedBranch 编辑保护分支
	UpdateProtectedBranch(context.Context, UpdateProtectedBranchReqDTO) error
}
