package filesrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	UploadIcon(context.Context, UploadIconReqDTO) (string, error)
	GetIcon(context.Context, GetIconReqDTO) (string, error)
	UploadNormal(context.Context, UploadNormalReqDTO) (string, error)
	GetNormal(context.Context, GetNormalReqDTO) (string, error)
	UploadAvatar(context.Context, UploadAvatarReqDTO) (string, error)
	GetAvatar(context.Context, GetAvatarReqDTO) (string, error)
}
