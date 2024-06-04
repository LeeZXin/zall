package webhooksrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/eventbus"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct{}

func newOuterService() OuterService {
	psub.Subscribe(eventbus.PullRequestEventTopic, func(data any) {
		event, ok := data.(eventbus.PullRequestEvent)
		if ok {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			// 触发webhook
			hookList, err := webhookmd.ListWebhook(ctx, event.RepoId)
			if err == nil && len(hookList) > 0 {
				req := &webhook.PullRequestEventReq{
					PrId:    event.PrId,
					PrTitle: event.PrTitle,
					Action:  webhook.PullRequestAction(event.Action),
					BaseRepoReq: webhook.BaseRepoReq{
						RepoId:    event.RepoId,
						RepoName:  event.RepoName,
						Account:   event.Account,
						EventTime: event.EventTime.UnixMilli(),
					},
				}
				for _, hook := range hookList {
					if hook.Events.Has(webhook.PullRequestEvent) {
						webhook.TriggerWebhook(hook.HookUrl, hook.Secret, req)
					}
				}
			}
		}
	})
	psub.Subscribe(eventbus.GitRepoEventTopic, func(data any) {
		event, ok := data.(eventbus.GitRepoEvent)
		if ok {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			// 触发webhook
			hookList, err := webhookmd.ListWebhook(ctx, event.RepoId)
			if err == nil && len(hookList) > 0 {
				req := &webhook.GitRepoEventReq{
					Action: webhook.GitRepoAction(event.Action),
					BaseRepoReq: webhook.BaseRepoReq{
						RepoId:    event.RepoId,
						RepoName:  event.Name,
						Account:   event.Operator,
						EventTime: event.EventTime.UnixMilli(),
					},
				}
				for _, hook := range hookList {
					if hook.Events.Has(webhook.GitRepoEvent) {
						webhook.TriggerWebhook(hook.HookUrl, hook.Secret, req)
					}
				}
			}
		}
	})
	psub.Subscribe(eventbus.ProtectedBranchEventTopic, func(data any) {
		event, ok := data.(eventbus.ProtectedBranchEvent)
		if ok {
			ctx, closer := xormstore.Context(context.Background())
			defer closer.Close()
			// 触发webhook
			hookList, err := webhookmd.ListWebhook(ctx, event.RepoId)
			if err == nil && len(hookList) > 0 {
				req := &webhook.ProtectedBranchEventReq{
					BaseRepoReq: webhook.BaseRepoReq{
						RepoId:    event.RepoId,
						RepoName:  event.Name,
						Account:   event.Operator,
						EventTime: event.EventTime.UnixMilli(),
					},
					Action: webhook.ProtectedBranchAction(event.Action),
				}
				if event.Before != nil {
					req.Before = &webhook.ProtectedBranchObj{
						Pattern:            event.Before.Pattern,
						ProtectedBranchCfg: event.Before.ProtectedBranchCfg,
					}
				}
				if event.After != nil {
					req.After = &webhook.ProtectedBranchObj{
						Pattern:            event.After.Pattern,
						ProtectedBranchCfg: event.After.ProtectedBranchCfg,
					}
				}
				for _, hook := range hookList {
					if hook.Events.Has(webhook.ProtectedBranchEvent) {
						webhook.TriggerWebhook(hook.HookUrl, hook.Secret, req)
					}
				}
			}
		}
	})
	return new(outerImpl)
}

// CreateWebhook 新增webhook
func (*outerImpl) CreateWebhook(ctx context.Context, reqDTO CreateWebhookReqDTO) (err error) {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator); err != nil {
		return
	}
	if err = webhookmd.InsertWebhook(ctx, webhookmd.InsertWebhookReqDTO{
		RepoId:  reqDTO.RepoId,
		HookUrl: reqDTO.HookUrl,
		Secret:  reqDTO.Secret,
		Events:  reqDTO.Events,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

func (*outerImpl) UpdateWebhook(ctx context.Context, reqDTO UpdateWebhookReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.WebhookSrvKeysVO.UpdateWebhook),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, err = checkPermByHookId(ctx, reqDTO.WebhookId, reqDTO.Operator); err != nil {
		return
	}
	_, err = webhookmd.UpdateWebhook(ctx, webhookmd.UpdateWebhookReqDTO{
		Id:      reqDTO.WebhookId,
		HookUrl: reqDTO.HookUrl,
		Secret:  reqDTO.Secret,
		Events:  reqDTO.Events,
	})
	if err != nil {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if _, err = checkPermByHookId(ctx, reqDTO.WebhookId, reqDTO.Operator); err != nil {
		return
	}
	if _, err = webhookmd.DeleteById(ctx, reqDTO.WebhookId); err != nil {
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
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.RepoId, reqDTO.Operator); err != nil {
		return nil, err
	}
	webhookList, err := webhookmd.ListWebhook(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(webhookList, func(t webhookmd.Webhook) (WebhookDTO, error) {
		return WebhookDTO{
			Id:      t.Id,
			RepoId:  t.RepoId,
			HookUrl: t.HookUrl,
			Secret:  t.Secret,
			Events:  *t.Events,
		}, nil
	})
}

// PingWebhook ping
func (*outerImpl) PingWebhook(ctx context.Context, reqDTO PingWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	hook, err := checkPermByHookId(ctx, reqDTO.WebhookId, reqDTO.Operator)
	if err != nil {
		return err
	}
	req := &webhook.PingEventReq{
		EventTime: time.Now().UnixMilli(),
	}
	err = webhook.Post(ctx, hook.HookUrl, hook.Secret, req)
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

func checkPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) error {
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

func checkPermByHookId(ctx context.Context, hookId int64, operator apisession.UserInfo) (webhookmd.Webhook, error) {
	hook, b, err := webhookmd.GetById(ctx, hookId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return webhookmd.Webhook{}, util.InternalError(err)
	}
	if !b {
		return webhookmd.Webhook{}, util.InvalidArgsError()
	}
	return hook, checkPerm(ctx, hook.RepoId, operator)
}
