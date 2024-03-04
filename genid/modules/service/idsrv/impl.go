package idsrv

import (
	"context"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/bwmarrin/snowflake"
	"math/rand"
)

var (
	Outer = newOuterImpl()
)

type outerImpl struct {
	snode *snowflake.Node
}

func newOuterImpl() OuterService {
	node := static.GetInt("snowflake.node")
	if node < 0 || node > 1023 {
		node = rand.Int() % 1024
	}
	snode, err := snowflake.NewNode(int64(node))
	if err != nil {
		logger.Logger.Fatalf("new snowflake.node: %v,  err: %v", node, err)
	}
	return &outerImpl{
		snode: snode,
	}
}

func (o *outerImpl) GenSnowflakeIds(_ context.Context, batchNum int) []int64 {
	if batchNum <= 0 {
		batchNum = 1
	}
	ret := make([]int64, 0, batchNum)
	for i := 0; i < batchNum; i++ {
		id := o.snode.Generate()
		ret = append(ret, id.Int64())
	}
	return ret
}
