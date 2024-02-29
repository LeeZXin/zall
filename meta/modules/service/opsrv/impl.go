package opsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/opmd"
	"github.com/LeeZXin/zall/pkg/wraperr"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"runtime"
	"time"
)

var (
	// 默认带实现队列长度为1024的协程池
	e, _ = executor.NewExecutor(runtime.GOMAXPROCS(0), 1024, 5*time.Minute, executor.CallerRunsStrategy)
)

type innerImpl struct{}

func (*innerImpl) InsertOpLog(_ context.Context, reqDTO InsertOpLogReqDTO) {
	e.Execute(func() {
		ctx, closer := mysqlstore.Context(nil)
		defer closer.Close()
		if reqDTO.EventTime.IsZero() {
			reqDTO.EventTime = time.Now()
		}
		req := opmd.InsertLogReqDTO{
			Operator:   reqDTO.Account,
			OpDesc:     reqDTO.OpDesc,
			ReqContent: reqDTO.ReqContent,
		}
		if reqDTO.Err != nil {
			ierr, b := reqDTO.Err.(*wraperr.Err)
			if b && ierr.Internal != nil {
				req.ErrMsg = ierr.Internal.Error()
			} else {
				req.ErrMsg = reqDTO.Err.Error()
			}
		}
		err := opmd.InsertLog(ctx, req)
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}
