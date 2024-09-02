package tpfeishusrv

import (
	"context"
	"github.com/LeeZXin/zall/pkg/feishuapi"
	"github.com/LeeZXin/zall/thirdpart/modules/model/tpfeishumd"
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
		"feishu-access-token-lock",
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
					time.Sleep(30 * time.Second)
				}
			},
			Leaser: leaser,
			// 抢锁失败 空转等待时间
			WaitDuration: 30 * time.Second,
			// 锁过期时间有20秒 每8秒续命 至少2次续命成功的机会
			RenewDuration: 8 * time.Second,
			GrantCallback: func(err error, b bool) {
				if err != nil {
					logger.Logger.Errorf("feishu access token task %s grant lease failed with err: %v", common.GetInstanceId(), err)
					return
				}
				if b {
					logger.Logger.Infof("feishu access token task grant lease success: %v", common.GetInstanceId())
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
	err := tpfeishumd.IterateAccessToken(ctx, time.Now().Add(10*time.Minute).UnixMilli(), func(at *tpfeishumd.AccessToken) error {
		token, tenantToken, expireIn, err := feishuapi.GetAppAccessToken(runCtx, static.GetString("feishu.accessToken.url"), at.AppId, at.Secret)
		if err != nil {
			// 忽略api错误
			logger.Logger.Errorf("feishu access token appId: %v failed with err: %v", at.AppId, err)
			return nil
		}
		return updateAppAccessToken(at.Id, token, tenantToken, expireIn)
	})
	if err != nil {
		logger.Logger.Error(err)
	}
}

func updateAppAccessToken(id int64, token, tenantToken string, expireIn int) error {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, err := tpfeishumd.UpdateAccessTokenToken(ctx, tpfeishumd.UpdateAccessTokenTokenReqDTO{
		Id:          id,
		Token:       token,
		TenantToken: tenantToken,
		ExpireTime:  time.Now().Add(time.Duration(expireIn) * time.Second).UnixMilli(),
	})
	if err != nil {
		logger.Logger.Error(err)
	}
	return err
}
