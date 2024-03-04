package idsrv

import "context"

type OuterService interface {
	GenSnowflakeIds(context.Context, int) []int64
	InsertGenerator(context.Context, string, int64) error
	GenerateIdByBizName(context.Context, string, int) ([]int64, error)
}
