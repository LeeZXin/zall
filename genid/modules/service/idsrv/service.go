package idsrv

import "context"

type OuterService interface {
	GenSnowflakeIds(context.Context, int) []int64
}
