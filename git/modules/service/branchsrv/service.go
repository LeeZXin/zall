package branchsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
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
				// 触发teamhook
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.ProtectedBranchEventCreateAction:
						return events.ProtectedBranch.Create
					case event.ProtectedBranchEventUpdateAction:
						return events.ProtectedBranch.Update
					case event.ProtectedBranchEventDeleteAction:
						return events.ProtectedBranch.Delete
					default:
						return false
					}
				})
			}
		})
	})
}

// CreateProtectedBranch 添加保护分支
func CreateProtectedBranch(ctx context.Context, reqDTO CreateProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, team, err := checkManageProtectedBranchPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 检查审批人
	if len(reqDTO.Cfg.ReviewerList) > 0 {
		accounts, err := teammd.ListUserAccountByTeamId(ctx, repo.TeamId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		// 是否包含所有的reviewerList
		reviewerSet := hashset.NewHashSet([]string(reqDTO.Cfg.ReviewerList)...)
		reviewerSet.Remove(accounts...)
		if reviewerSet.Size() != 0 {
			return util.InvalidArgsError()
		}
	}
	err = branchmd.InsertProtectedBranch(ctx, branchmd.InsertProtectedBranchReqDTO{
		RepoId:  reqDTO.RepoId,
		Pattern: reqDTO.Pattern,
		Cfg:     reqDTO.Cfg,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 通知其他领域
	notifyEvent(
		repo,
		team,
		reqDTO.Operator,
		nil,
		&branch.ProtectedBranch{
			Pattern:            reqDTO.Pattern,
			ProtectedBranchCfg: reqDTO.Cfg,
		},
		event.ProtectedBranchEventCreateAction,
	)
	return nil
}

// UpdateProtectedBranch 编辑保护分支
func UpdateProtectedBranch(ctx context.Context, reqDTO UpdateProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pb, b, err := branchmd.GetProtectedBranchById(ctx, reqDTO.ProtectedBranchId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	repo, team, err := checkManageProtectedBranchPerm(ctx, pb.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = branchmd.UpdateProtectedBranch(ctx, branchmd.UpdateProtectedBranchReqDTO{
		Id:      reqDTO.ProtectedBranchId,
		Pattern: reqDTO.Pattern,
		Cfg:     reqDTO.Cfg,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 上报事件
	notifyEvent(
		repo,
		team,
		reqDTO.Operator,
		&branch.ProtectedBranch{
			Pattern:            pb.Pattern,
			ProtectedBranchCfg: pb.GetCfg(),
		},
		&branch.ProtectedBranch{
			Pattern:            reqDTO.Pattern,
			ProtectedBranchCfg: reqDTO.Cfg,
		},
		event.ProtectedBranchEventUpdateAction,
	)
	return nil
}

// DeleteProtectedBranch 删除保护分支
func DeleteProtectedBranch(ctx context.Context, reqDTO DeleteProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pb, b, err := branchmd.GetProtectedBranchById(ctx, reqDTO.ProtectedBranchId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	repo, team, err := checkManageProtectedBranchPerm(ctx, pb.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = branchmd.DeleteProtectedBranchById(ctx, reqDTO.ProtectedBranchId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(
		repo,
		team,
		reqDTO.Operator,
		&branch.ProtectedBranch{
			Pattern:            pb.Pattern,
			ProtectedBranchCfg: pb.GetCfg(),
		},
		nil,
		event.ProtectedBranchEventDeleteAction,
	)
	return nil
}

func notifyEvent(repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, before, after *branch.ProtectedBranch, action event.ProtectedBranchEventAction) {
	initPsub()
	psub.Publish(event.ProtectedBranchTopic, event.ProtectedBranchEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoPath: repo.Path,
			RepoId:   repo.Id,
			RepoName: repo.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
		},
		Action: action,
		Before: before,
		After:  after,
	})
}

// ListProtectedBranch 保护分支列表
func ListProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) ([]ProtectedBranchDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkManageProtectedBranchPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	branchList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := listutil.MapNe(branchList, func(t branchmd.ProtectedBranch) ProtectedBranchDTO {
		return ProtectedBranchDTO{
			Id:      t.Id,
			Pattern: t.Pattern,
			Cfg:     *t.Cfg,
		}
	})
	return ret, nil
}

// checkManageProtectedBranchPerm 检查权限
func checkManageProtectedBranchPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	// 检查仓库是否存在
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
	// 如果不是 检查用户组权限
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		return repo, team, util.InternalError(err)
	}
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanManageProtectedBranch) {
		return repo, team, nil
	}
	return repo, team, util.UnauthorizedError()
}
