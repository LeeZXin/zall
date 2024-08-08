package githook

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/git/modules/service/workflowsrv"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/git/repo/server/store"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/githook"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
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
	// 检查仓库是否归档
	if repo.IsArchived {
		return util.NewBizErr(apicode.OperationFailedErrCode, i18n.SystemForbidden)
	}
	repoPath := filepath.Join(git.RepoDir(), repo.Path)
	// 检查仓库大小是否大于配置大小
	if repo.GetCfg().GitLimitSize > 0 {
		repoSize, err := git.GetDirSize(repoPath)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if repoSize > repo.Cfg.GitLimitSize {
			return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.RepoSizeExceedLimit, util.VolumeReadable(repo.Cfg.GitLimitSize))
		}
	}
	var pbList branchmd.ProtectedBranchList
	for _, info := range opts.RevInfoList {
		ref := info.Ref
		// 是分支
		if strings.HasPrefix(ref, git.BranchPrefix) {
			// 检查是否是保护分支
			if pbList == nil {
				// 懒加载一下
				pbList, err = branchmd.ListProtectedBranch(ctx, opts.RepoId)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return util.InternalError(err)
				}
			}
			isProtectedBranch, protectedBranch := pbList.IsProtectedBranch(strings.TrimPrefix(ref, git.BranchPrefix))
			if isProtectedBranch {
				// 检查该分支是否可推送
				cfg := protectedBranch.GetCfg()
				switch cfg.PushOption {
				case branch.NotAllowPush:
					return util.NewBizErr(apicode.ProtectedBranchNotAllowPushCode, i18n.ProtectedBranchNotAllowPush)
				case branch.WhiteListPush:
					if !cfg.PushWhiteList.Contains(opts.PusherAccount) {
						return util.NewBizErr(apicode.ProtectedBranchNotAllowPushCode, i18n.ProtectedBranchNotAllowPush)
					}
				case branch.AllowPush:
				}
				// 不允许删除保护分支
				if info.NewCommitId == git.ZeroCommitId {
					return util.NewBizErr(apicode.ForcePushForbiddenCode, i18n.ProtectedBranchNotAllowDelete)
				}
				if info.OldCommitId != git.ZeroCommitId {
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
	// 获取仓库信息
	repo, b, err := repomd.GetByRepoId(ctx, opts.RepoId)
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		return
	}
	// 更新最后操作时间
	_, err = repomd.UpdateLastOperated(ctx, opts.RepoId, time.Now())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	gitSize, lfsSize, err := store.Srv.GetRepoSize(ctx, reqvo.GetRepoSizeReq{
		RepoPath: repo.Path,
	})
	if err == nil {
		fmt.Println(gitSize, lfsSize)
		// 更新仓库大小
		err = repomd.UpdateGitSizeAndLfsSize(ctx, opts.RepoId, gitSize, lfsSize)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	} else {
		logger.Logger.WithContext(ctx).Error(err)
	}
	// 查找webhook
	webhookList, err := webhookmd.ListWebhook(ctx, repo.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return
	}
	// 查找工作流
	workflowList, err := workflowmd.ListWorkflowByRepoId(ctx, repo.Id)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	now := time.Now()
	for _, revInfo := range opts.RevInfoList {
		var refType string
		if strings.HasPrefix(revInfo.Ref, git.BranchPrefix) {
			refType = "commit"
		} else if strings.HasPrefix(revInfo.Ref, git.TagPrefix) {
			refType = "tag"
		} else {
			continue
		}
		// 触发webhook
		for _, hook := range webhookList {
			if hook.Events.Has(webhook.GitPushEvent) {
				webhook.TriggerWebhook(hook.HookUrl, hook.Secret, &webhook.GitPushEventReq{
					RefType:     refType,
					Ref:         revInfo.Ref,
					OldCommitId: revInfo.OldCommitId,
					NewCommitId: revInfo.NewCommitId,
					BaseRepoReq: webhook.BaseRepoReq{
						RepoId:    repo.Id,
						RepoName:  repo.Name,
						Account:   opts.PusherAccount,
						EventTime: now.UnixMilli(),
					},
				})
			}
		}
		// 触发工作流
		for _, wf := range workflowList {
			if refType == "commit" {
				ref := strings.TrimPrefix(revInfo.Ref, git.BranchPrefix)
				if wf.Source.MatchBranchBySource(workflowmd.BranchTriggerSource, ref) {
					workflowsrv.Inner.Execute(wf, workflowsrv.ExecuteWorkflowReqDTO{
						RepoPath:    repo.Path,
						Operator:    opts.PusherAccount,
						TriggerType: workflowmd.HookTriggerType,
						Branch:      ref,
						PrId:        0,
					})
				}
			}
		}
	}
}
