package githook

import (
	"context"
	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/LeeZXin/zall/git/modules/model/actiontaskmd"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/githook"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
	"time"
)

type hookImpl struct{}

func NewHook() Hook {
	return new(hookImpl)
}

func (*hookImpl) PreReceive(ctx context.Context, opts githook.Opts) error {
	ctx, closer := mysqlstore.Context(ctx)
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
	var pbList []branchmd.ProtectedBranchDTO
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
					if opts.PrId > 0 && len(pb.Cfg.DirectPushList) > 0 {
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
	ctx, closer := mysqlstore.Context(ctx)
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
		pushHookList []webhookmd.WebhookDTO
		tagHookList  []webhookmd.WebhookDTO
		err          error
	)
	for _, revInfo := range opts.RevInfoList {
		// 分支push
		if strings.HasPrefix(revInfo.RefName, git.BranchPrefix) {
			if pushHookList == nil {
				pushHookList, err = webhookmd.ListWebhook(ctx, repo.RepoId, webhookmd.PushHook)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return
				}
			}
			// 触发hook
			triggerHook(pushHookList, repo, revInfo, operator, false)
		} else if strings.HasPrefix(revInfo.RefName, git.TagPrefix) {
			// tag commit
			if tagHookList == nil {
				tagHookList, err = webhookmd.ListWebhook(ctx, repo.RepoId, webhookmd.TagHook)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return
				}
			}
			// 触发webhook
			triggerHook(tagHookList, repo, revInfo, operator, true)
		}
	}
	// 处理action
	actions, err := repomd.ListAction(ctx, repo.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	if len(actions) == 0 {
		return
	}
	graphs := make(map[*action.Graph]repomd.Action, len(actions))
	for _, item := range actions {
		var cfg action.GraphCfg
		err = yaml.Unmarshal([]byte(item.Content), &cfg)
		if err == nil {
			if graph, err := cfg.ConvertToGraph(); err == nil {
				graphs[graph] = item
			}
		}
	}
	if len(graphs) == 0 {
		return
	}
	// 解析并触发action
	for graph, actionCfg := range graphs {
		refs, b := graph.GetSupportedRefs(action.PushAction)
		if !b {
			continue
		}
		for _, revInfo := range opts.RevInfoList {
			// 不要删除分支
			if revInfo.NewCommitId != git.ZeroCommitId &&
				util.WildcardMatchBranches(refs.Branches, util.BaseRefName(revInfo.RefName)) {
				// 触发actions
				triggerActions(repo, revInfo, operator, actionCfg.AssignInstance, actionCfg.Content)
			}
		}
	}
}

func triggerHook(hookList []webhookmd.WebhookDTO, repo repomd.RepoInfo, revInfo githook.RevInfo, operator usermd.UserInfo, isTag bool) {
	req := webhook.GitReceiveHook{
		RepoId:    repo.RepoId,
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

func triggerActions(repo repomd.RepoInfo, revInfo githook.RevInfo, operator usermd.UserInfo, instanceId, yamlContent string) {
	req := action.Webhook{
		RepoId:    repo.RepoId,
		RepoName:  repo.Name,
		Ref:       revInfo.RefName,
		EventTime: time.Now().UnixMilli(),
		Operator: git.User{
			Account: operator.Account,
			Email:   operator.Email,
		},
		TriggerType: actiontaskmd.SysTriggerType.Int(),
		YamlContent: yamlContent,
	}
	// 负载均衡选择一个节点
	instance, b, err := actionsrv.SelectAndIncrJobCountInstances(context.Background(), instanceId)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	// 没有可用节点
	if !b {
		return
	}
	action.TriggerActionHook(req, instance.InstanceHost)
}
