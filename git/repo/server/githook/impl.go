package githook

import (
	"context"
	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/githook"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"path/filepath"
	"strings"
	"time"
)

type hookImpl struct{}

func NewHook() Hook {
	return new(hookImpl)
}

func (*hookImpl) PreReceive(ctx context.Context, opts githook.Opts) error {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, opts.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	repoPath := filepath.Join(git.RepoDir(), repo.Path)
	// 检查仓库大小是否大于配置大小
	if repo.Cfg.MaxGitLimitSize > 0 {
		repoSize, err := git.GetRepoSize(repoPath)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if repoSize > repo.Cfg.MaxGitLimitSize {
			return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.RepoSizeExceedLimit, util.VolumeReadable(repo.Cfg.MaxGitLimitSize))
		}
	}
	var pbList []branchmd.ProtectedBranch
	for _, info := range opts.RevInfoList {
		name := info.RefName
		// 是分支
		if strings.HasPrefix(name, git.BranchPrefix) {
			// 检查是否是保护分支
			if pbList == nil {
				// 懒加载一下
				pbList, err = branchmd.ListProtectedBranch(ctx, opts.RepoId)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return util.InternalError(err)
				}
			}
			name = strings.TrimPrefix(name, git.BranchPrefix)
			for _, pb := range pbList {
				// 通配符匹配 是保护分支
				if wildcard.Match(pb.Branch, name) {
					// 只有可推送名单里面才能直接push
					if opts.PrId > 0 && pb.Cfg != nil && len(pb.Cfg.DirectPushList) > 0 {
						// prId为空说明不是来自合并请求的push
						contains, _ := listutil.Contains(pb.Cfg.DirectPushList, func(account string) (bool, error) {
							return account == opts.PusherAccount, nil
						})
						if !contains {
							return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowDirectPush)
						}
					}
					// 不允许删除保护分支
					if info.NewCommitId == git.ZeroCommitId {
						return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowDelete)
					}
					// 检查push -f
					isForcePush, err := git.DetectForcePush(ctx,
						repoPath,
						info.OldCommitId,
						info.NewCommitId,
						git.DetectForcePushEnv{
							ObjectDirectory:              opts.ObjectDirectory,
							AlternativeObjectDirectories: opts.AlternativeObjectDirectories,
							QuarantinePath:               opts.QuarantinePath,
						})
					if err != nil {
						logger.Logger.WithContext(ctx).Error(err)
						return util.InternalError(err)
					}
					if isForcePush {
						// 禁止push -f
						return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowForcePush)
					}
				}
			}
		}
	}
	return nil
}

func (*hookImpl) PostReceive(ctx context.Context, opts githook.Opts) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := reposrv.Inner.GetByRepoId(ctx, opts.RepoId)
	if !b {
		return
	}
	operator, b := usersrv.Inner.GetByAccount(ctx, opts.PusherAccount)
	if !b {
		return
	}
	// 查找webhook
	var (
		pushHookList []webhookmd.Webhook
		tagHookList  []webhookmd.Webhook
		err          error
	)
	for _, revInfo := range opts.RevInfoList {
		// 分支push
		if strings.HasPrefix(revInfo.RefName, git.BranchPrefix) {
			branch := strings.TrimPrefix(revInfo.RefName, git.BranchPrefix)
			if pushHookList == nil {
				pushHookList, err = webhookmd.ListWebhook(ctx, repo.Id, webhookmd.PushHook)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return
				}
			}
			hookList, _ := listutil.Filter(pushHookList, func(t webhookmd.Webhook) (bool, error) {
				return wildcard.Match(t.WildBranch, branch), nil
			})
			// 触发hook
			triggerWebhook(hookList, repo, revInfo, operator, false)
		} else if strings.HasPrefix(revInfo.RefName, git.TagPrefix) {
			tag := strings.TrimPrefix(revInfo.RefName, git.TagPrefix)
			// tag commit
			if tagHookList == nil {
				tagHookList, err = webhookmd.ListWebhook(ctx, repo.Id, webhookmd.TagHook)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return
				}
			}
			hookList, _ := listutil.Filter(tagHookList, func(t webhookmd.Webhook) (bool, error) {
				return wildcard.Match(t.WildTag, tag), nil
			})
			// 触发webhook
			triggerWebhook(hookList, repo, revInfo, operator, true)
		}
	}
}

func triggerWebhook(hookList []webhookmd.Webhook, repo repomd.Repo, revInfo githook.RevInfo, operator usermd.UserInfo, isTag bool) {
	if len(hookList) == 0 {
		return
	}
	req := webhook.GitReceiveHook{
		RepoId:    repo.Id,
		RepoName:  repo.Name,
		IsCreated: revInfo.OldCommitId == git.ZeroCommitId,
		IsDeleted: revInfo.NewCommitId == git.ZeroCommitId,
		Ref:       revInfo.RefName,
		EventTime: time.Now().UnixMilli(),
		Operator: git.User{
			Account: operator.Account,
			Email:   operator.Email,
		},
		IsTagPush: isTag,
	}
	for _, hook := range hookList {
		webhook.TriggerGitHook(hook.HookUrl, hook.HttpHeaders, req)
	}
}
