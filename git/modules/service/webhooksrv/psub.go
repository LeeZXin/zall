package webhooksrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
)

var (
	psubOnce = sync.Once{}
)

func InitPsub() {
	psubOnce.Do(func() {
		psub.Subscribe(event.PullRequestTopic, func(data any) {
			req, ok := data.(event.PullRequestEvent)
			if ok {
				ctx, closer := xormstore.Context(context.Background())
				// 触发webhook
				hookList, err := webhookmd.ListWebhookByRepoId(ctx, req.RepoId)
				closer.Close()
				if err == nil && len(hookList) > 0 {
					for _, hook := range hookList {
						if hook.GetEvents().PullRequest {
							webhook.TriggerWebhook(hook.HookUrl, hook.Secret, &req)
						}
					}
				}
			}
		})
		psub.Subscribe(event.GitRepoTopic, func(data any) {
			req, ok := data.(event.GitRepoEvent)
			if ok {
				ctx, closer := xormstore.Context(context.Background())
				// 触发webhook
				hookList, err := webhookmd.ListWebhookByRepoId(ctx, req.RepoId)
				closer.Close()
				if err == nil && len(hookList) > 0 {
					for _, hook := range hookList {
						if hook.GetEvents().GitRepo {
							webhook.TriggerWebhook(hook.HookUrl, hook.Secret, &req)
						}
					}
				}
			}
		})
		psub.Subscribe(event.ProtectedBranchTopic, func(data any) {
			req, ok := data.(event.ProtectedBranchEvent)
			if ok {
				ctx, closer := xormstore.Context(context.Background())
				// 触发webhook
				hookList, err := webhookmd.ListWebhookByRepoId(ctx, req.RepoId)
				closer.Close()
				if err == nil && len(hookList) > 0 {
					for _, hook := range hookList {
						if hook.GetEvents().ProtectedBranch {
							webhook.TriggerWebhook(hook.HookUrl, hook.Secret, &req)
						}
					}
				}
			}
		})
	})
}
