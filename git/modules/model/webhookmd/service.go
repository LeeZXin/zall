package webhookmd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertWebhook(ctx context.Context, reqDTO InsertWebhookReqDTO) (Webhook, error) {
	ret := Webhook{
		RepoId:  reqDTO.RepoId,
		HookUrl: reqDTO.HookUrl,
		Secret:  reqDTO.Secret,
		Events: &xormutil.Conversion[webhook.Events]{
			Data: reqDTO.Events,
		},
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateWebhook(ctx context.Context, reqDTO UpdateWebhookReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("hook_url", "secret", "events").
		Update(&Webhook{
			HookUrl: reqDTO.HookUrl,
			Secret:  reqDTO.Secret,
			Events: &xormutil.Conversion[webhook.Events]{
				Data: reqDTO.Events,
			},
		})
	return rows == 1, err
}

func DeleteById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Webhook))
	return rows == 1, err
}

func ListWebhookByRepoId(ctx context.Context, repoId int64) ([]Webhook, error) {
	ret := make([]Webhook, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Find(&ret)
	return ret, err
}

func GetWebhookById(ctx context.Context, id int64) (Webhook, bool, error) {
	var ret Webhook
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}
