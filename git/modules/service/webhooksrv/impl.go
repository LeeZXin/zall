package webhooksrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
)

type outerImpl struct{}

func (*outerImpl) InsertWebHook(ctx context.Context, reqDTO InsertWebhookReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.WebhookSrvKeysVO.InsertWebhook),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator); err != nil {
		return
	}
	if err = webhookmd.InsertWebhook(ctx, webhookmd.InsertWebhookReqDTO{
		RepoId:      reqDTO.RepoId,
		HookUrl:     reqDTO.HookUrl,
		HttpHeaders: reqDTO.HttpHeaders,
		HookType:    reqDTO.HookType,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) DeleteWebhook(ctx context.Context, reqDTO DeleteWebhookReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.WebhookSrvKeysVO.DeleteWebhook),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	hook, b, err := webhookmd.GetById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	if err = checkPerm(ctx, hook.RepoId, reqDTO.Operator); err != nil {
		return
	}
	if _, err = webhookmd.DeleteById(ctx, reqDTO.Id); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListWebhook(ctx context.Context, reqDTO ListWebhookReqDTO) ([]WebhookDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := mysqlstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator); err != nil {
		return nil, err
	}
	webhookList, err := webhookmd.ListWebhook(ctx, reqDTO.RepoId, reqDTO.HookType)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(webhookList, func(t webhookmd.WebhookDTO) (WebhookDTO, error) {
		return WebhookDTO{
			Id:          t.Id,
			RepoId:      t.RepoId,
			HookUrl:     t.HookUrl,
			HttpHeaders: t.HttpHeaders,
			HookType:    t.HookType,
		}, nil
	})
}

func checkPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) error {
	repo, b := reposrv.Inner.GetByRepoId(ctx, repoId)
	if !b {
		return util.InvalidArgsError()
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b || !p.PermDetail.GetRepoPerm(repoId).CanHandleWebhook {
		return util.UnauthorizedError()
	}
	return nil
}
