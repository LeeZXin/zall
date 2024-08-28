package webhooksrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.GitWebhookTopic, func(data any) {
			req, ok := data.(event.GitWebhookEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.GitWebhookEventCreateAction:
						return events.GitWebhook.Create
					case event.GitWebhookEventUpdateAction:
						return events.GitWebhook.Update
					case event.GitWebhookEventDeleteAction:
						return events.GitWebhook.Delete
					default:
						return false
					}
				})
			}
		})
	})
}

// CreateWebhook 新增webhook
func CreateWebhook(ctx context.Context, reqDTO CreateWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, team, err := checkManageWebhookPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	w, err := webhookmd.InsertWebhook(ctx, webhookmd.InsertWebhookReqDTO{
		RepoId:  reqDTO.RepoId,
		HookUrl: reqDTO.HookUrl,
		Secret:  reqDTO.Secret,
		Events:  reqDTO.Events,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(w, repo, team, reqDTO.Operator, event.GitWebhookEventCreateAction)
	return nil
}

// UpdateWebhook 编辑webhook
func UpdateWebhook(ctx context.Context, reqDTO UpdateWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	w, repo, team, err := checkManageWebhookPermByHookId(ctx, reqDTO.Id, reqDTO.Operator)
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
	notifyEvent(w, repo, team, reqDTO.Operator, event.GitWebhookEventUpdateAction)
	return nil
}

// DeleteWebhook 删除webhook
func DeleteWebhook(ctx context.Context, reqDTO DeleteWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	w, repo, team, err := checkManageWebhookPermByHookId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = webhookmd.DeleteById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(w, repo, team, reqDTO.Operator, event.GitWebhookEventDeleteAction)
	return nil
}

// ListWebhook 列表webhook
func ListWebhook(ctx context.Context, reqDTO ListWebhookReqDTO) ([]WebhookDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkManageWebhookPermByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
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
func PingWebhook(ctx context.Context, reqDTO PingWebhookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	hook, _, _, err := checkManageWebhookPermByHookId(ctx, reqDTO.WebhookId, reqDTO.Operator)
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

func checkManageWebhookPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return repo, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, team, util.InternalError(err)
	}
	if !b || (!p.IsAdmin && !p.PermDetail.GetRepoPerm(repo.Id).CanManageWebhook) {
		return repo, team, util.UnauthorizedError()
	}
	return repo, team, nil
}

func checkManageWebhookPermByHookId(ctx context.Context, hookId int64, operator apisession.UserInfo) (webhookmd.Webhook, repomd.Repo, teammd.Team, error) {
	hook, b, err := webhookmd.GetWebhookById(ctx, hookId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return webhookmd.Webhook{}, repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return webhookmd.Webhook{}, repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	repo, team, err := checkManageWebhookPermByRepoId(ctx, hook.RepoId, operator)
	return hook, repo, team, err
}

func notifyEvent(w webhookmd.Webhook, repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, action event.GitWebhookEventAction) {
	initPsub()
	psub.Publish(event.GitWebhookTopic, event.GitWebhookEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoId:   repo.Id,
			RepoPath: repo.Path,
			RepoName: repo.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action:     action,
		WebhookId:  w.Id,
		WebhookUrl: w.HookUrl,
	})
}
