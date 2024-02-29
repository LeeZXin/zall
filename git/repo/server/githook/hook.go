package githook

import (
	"context"
	"github.com/LeeZXin/zall/pkg/githook"
)

type Hook interface {
	PreReceive(context.Context, githook.Opts) error
	PostReceive(context.Context, githook.Opts)
}
