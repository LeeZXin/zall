package tpweworksrv

import (
	"context"
	"github.com/LeeZXin/zall/pkg/weworkapi"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpweworkmd"
	"github.com/LeeZXin/zsf-utils/lease"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

func InitGetAccessTokenTask() {
	leaser, _ := lease.NewDbLease(
		"wework-access-token-lock",
		common.GetInstanceId(),
		"zall_lock",
		xormstore.GetEngine(),
		20*time.Second,
	)
	stopFunc, _ := taskutil.RunMainLoopTask(
		taskutil.MainLoopTaskOpts{
			Handler: func(ctx context.Context) {
				for ctx.Err() == nil {
					doExecuteTask(ctx)
					time.Sleep(time.Minute)
				}
			},
			Leaser: leaser,
			// 抢锁失败 空转等待时间
			WaitDuration: 30 * time.Second,
			// 锁过期时间有20秒 每8秒续命 至少2次续命成功的机会
			RenewDuration: 8 * time.Second,
			GrantCallback: func(err error, b bool) {
				if err != nil {
					logger.Logger.Errorf("wework access token task %s grant lease failed with err: %v", common.GetInstanceId(), err)
					return
				}
				if b {
					logger.Logger.Infof("wework access token task grant lease success: %v", common.GetInstanceId())
				}
			},
		},
	)
	quit.AddShutdownHook(quit.ShutdownHook(stopFunc), true)
}

func doExecuteTask(runCtx context.Context) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	// 未来十分钟内要过期的token
	err := tpweworkmd.IterateAccessToken(ctx, time.Now().Add(10*time.Minute).UnixMilli(), func(at *tpweworkmd.AccessToken) error {
		return refreshAccessToken(runCtx, at)
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func refreshAccessToken(ctx context.Context, at *tpweworkmd.AccessToken) error {
	token, expireIn, err := weworkapi.GetAccessToken(ctx, static.GetString("wework.accessToken.url"), at.CorpId, at.Secret)
	if err != nil {
		// 忽略api错误
		logger.Logger.Errorf("wework access token corpId: %v failed with err: %v", at.CorpId, err)
		return nil
	}
	return updateAccessToken(at.Id, token, expireIn)
}

func updateAccessToken(id int64, token string, expireIn int) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, err := tpweworkmd.UpdateAccessTokenToken(ctx, tpweworkmd.UpdateAccessTokenTokenReqDTO{
		Id:         id,
		Token:      token,
		ExpireTime: time.Now().Add(time.Duration(expireIn) * time.Second).UnixMilli(),
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return err
}
