package webhooksrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct{}

// CreateWebhook 新增webhook
func (*outerImpl) CreateWebhook(ctx context.Context, reqDTO CreateWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkManageWebhookPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = webhookmd.InsertWebhook(ctx, webhookmd.InsertWebhookReqDTO{
		RepoId:  reqDTO.RepoId,
		HookUrl: reqDTO.HookUrl,
		Secret:  reqDTO.Secret,
		Events:  reqDTO.Events,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateWebhook 编辑webhook
func (*outerImpl) UpdateWebhook(ctx context.Context, reqDTO UpdateWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageWebhookPermByHookId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = webhookmd.UpdateWebhook(ctx, webhookmd.UpdateWebhookReqDTO{
		Id:      reqDTO.Id,
		HookUrl: reqDTO.HookUrl,
		Secret:  reqDTO.Secret,
		Events:  reqDTO.Events,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteWebhook 删除webhook
func (*outerImpl) DeleteWebhook(ctx context.Context, reqDTO DeleteWebhookReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, err = checkManageWebhookPermByHookId(ctx, reqDTO.Id, reqDTO.Operator); err != nil {
		return
	}
	if _, err = webhookmd.DeleteById(ctx, reqDTO.Id); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListWebhook 列表webhook
func (*outerImpl) ListWebhook(ctx context.Context, reqDTO ListWebhookReqDTO) ([]WebhookDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkManageWebhookPerm(ctx, reqDTO.RepoId, reqDTO.Operator); err != nil {
		return nil, err
	}
	webhookList, err := webhookmd.ListWebhookByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(webhookList, func(t webhookmd.Webhook) (WebhookDTO, error) {
		ret := WebhookDTO{
			Id:      t.Id,
			RepoId:  t.RepoId,
			HookUrl: t.HookUrl,
			Secret:  t.Secret,
		}
		if t.Events != nil {
			ret.Events = t.Events.Data
		}
		return ret, nil
	})
}

// PingWebhook ping
func (*outerImpl) PingWebhook(ctx context.Context, reqDTO PingWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	hook, err := checkManageWebhookPermByHookId(ctx, reqDTO.WebhookId, reqDTO.Operator)
	if err != nil {
		return err
	}
	req := &event.PingEvent{
		EventTime: time.Now().UnixMilli(),
	}
	err = webhook.Post(ctx, hook.HookUrl, hook.Secret, req)
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

func checkManageWebhookPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) error {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || (!p.IsAdmin && !p.PermDetail.GetRepoPerm(repo.Id).CanManageWebhook) {
		return util.UnauthorizedError()
	}
	return nil
}

func checkManageWebhookPermByHookId(ctx context.Context, hookId int64, operator apisession.UserInfo) (webhookmd.Webhook, error) {
	hook, b, err := webhookmd.GetWebhookById(ctx, hookId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return webhookmd.Webhook{}, util.InternalError(err)
	}
	if !b {
		return webhookmd.Webhook{}, util.InvalidArgsError()
	}
	return hook, checkManageWebhookPerm(ctx, hook.RepoId, operator)
}
