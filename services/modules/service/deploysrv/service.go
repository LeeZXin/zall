package deploysrv

import (
	"context"
	"github.com/LeeZXin/zall/pkg/deploy"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	// GetDeploy 获取部署配置
	GetDeploy(context.Context, GetDeployReqDTO) (deploy.Config, error)
	// UpdateDeploy 编辑部署配置
	UpdateDeploy(context.Context, UpdateDeployReqDTO) error
}
