package webhookmd

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertWebhook(ctx context.Context, reqDTO InsertWebhookReqDTO) error {
	hook := Webhook{
		RepoId:     reqDTO.RepoId,
		HookUrl:    reqDTO.HookUrl,
		HookType:   reqDTO.HookType,
		WildBranch: reqDTO.WildBranch,
		WildTag:    reqDTO.WildTag,
	}
	if reqDTO.HttpHeaders != nil {
		m, err := json.Marshal(reqDTO.HttpHeaders)
		if err == nil {
			hook.HttpHeaders = string(m)
		}
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&hook)
	return err
}

func UpdateWebhook(ctx context.Context, reqDTO UpdateWebhookReqDTO) (bool, error) {
	hook := &Webhook{
		HookUrl:    reqDTO.HookUrl,
		WildTag:    reqDTO.WildTag,
		WildBranch: reqDTO.WildBranch,
	}
	if reqDTO.HttpHeaders != nil {
		m, err := json.Marshal(reqDTO.HttpHeaders)
		if err == nil {
			hook.HttpHeaders = string(m)
		}
	}
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("hook_url", "http_headers", "wild_branch", "wild_tag").
		Limit(1).
		Update(hook)
	return rows == 1, err
}

func DeleteById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(Webhook))
	return rows == 1, err
}

func ListWebhook(ctx context.Context, repoId int64, hookType HookType) ([]Webhook, error) {
	ret := make([]Webhook, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		And("hook_type = ?", hookType).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetById(ctx context.Context, id int64) (Webhook, bool, error) {
	var ret Webhook
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}
