package idsrv

import (
	"context"
	"github.com/LeeZXin/zall/genid/modules/model/idmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"regexp"
	"time"
)

var (
	Outer               = newOuterImpl()
	validBizNamePattern = regexp.MustCompile("^\\w+$")
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

func (o *outerImpl) InsertGenerator(ctx context.Context, bizName string, currentId int64) error {
	if !validBizNamePattern.MatchString(bizName) {
		return util.InvalidArgsError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := idmd.GetByBizName(ctx, bizName)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		err = idmd.InsertGenerator(ctx, idmd.InsertGeneratorReqDTO{
			BizName:   bizName,
			CurrentId: currentId,
		})
		if err != nil && !xormutil.IsDuplicatedEntryError(err) {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
	}
	return nil
}

func (o *outerImpl) GenerateIdByBizName(ctx context.Context, bizName string, step int) ([]int64, error) {
	if step <= 0 {
		return []int64{}, nil
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	for i := 0; i < 10; i++ {
		gen, b, err := idmd.GetByBizName(ctx, bizName)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.NotExistsError()
		}
		b, err = idmd.UpdateCurrentId(ctx, idmd.UpdateCurrentIdReqDTO{
			BizName:   bizName,
			CurrentId: gen.CurrentId + int64(step),
			Version:   gen.Version,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if b {
			ret := make([]int64, 0, step)
			for j := 1; j < step+1; j++ {
				ret = append(ret, gen.CurrentId+int64(j))
			}
			return ret, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil, util.OperationFailedError()
}
