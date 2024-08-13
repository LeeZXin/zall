package branchsrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/oplogsrv"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/eventbus"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct{}

func (*outerImpl) CreateProtectedBranch(ctx context.Context, reqDTO CreateProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, err := checkManageProtectedBranchPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
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
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   reqDTO.RepoId,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.BranchSrvKeysVO.CreateProtectedBranch, reqDTO.Pattern),
		Req:      reqDTO,
	})
	// 通知其他领域
	notifyEventBus(
		repo,
		reqDTO.Operator.Account,
		nil,
		&eventbus.ProtectedBranchObj{
			Pattern:            reqDTO.Pattern,
			ProtectedBranchCfg: reqDTO.Cfg,
		},
		webhook.PbCreateAction,
	)
	return nil
}

// UpdateProtectedBranch 编辑保护分支
func (*outerImpl) UpdateProtectedBranch(ctx context.Context, reqDTO UpdateProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	protectedBranch, b, err := branchmd.GetById(ctx, reqDTO.ProtectedBranchId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	repo, err := checkManageProtectedBranchPerm(ctx, protectedBranch.RepoId, reqDTO.Operator)
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
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.BranchSrvKeysVO.UpdateProtectedBranch, reqDTO.Pattern),
		Req:      reqDTO,
	})
	notifyEventBus(
		repo,
		reqDTO.Operator.Account,
		&eventbus.ProtectedBranchObj{
			Pattern:            protectedBranch.Pattern,
			ProtectedBranchCfg: protectedBranch.GetCfg(),
		},
		&eventbus.ProtectedBranchObj{
			Pattern:            reqDTO.Pattern,
			ProtectedBranchCfg: reqDTO.Cfg,
		},
		webhook.PbUpdateAction,
	)
	return nil
}

func (*outerImpl) DeleteProtectedBranch(ctx context.Context, reqDTO DeleteProtectedBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	pb, b, err := branchmd.GetById(ctx, reqDTO.ProtectedBranchId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	repo, err := checkManageProtectedBranchPerm(ctx, pb.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = branchmd.DeleteById(ctx, reqDTO.ProtectedBranchId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 插入日志
	oplogsrv.Inner.InsertOpLog(oplogsrv.OpLog{
		RepoId:   repo.Id,
		Operator: reqDTO.Operator.Account,
		Log:      oplogsrv.FormatI18n(i18n.BranchSrvKeysVO.CreateProtectedBranch, pb.Pattern),
		Req:      reqDTO,
	})
	notifyEventBus(repo, reqDTO.Operator.Account, &eventbus.ProtectedBranchObj{
		Pattern:            pb.Pattern,
		ProtectedBranchCfg: pb.GetCfg(),
	}, nil, webhook.PbDeleteAction)
	return nil
}

func notifyEventBus(repo repomd.Repo, operator string, before, after *eventbus.ProtectedBranchObj, action webhook.ProtectedBranchAction) {
	psub.Publish(eventbus.ProtectedBranchEventTopic, eventbus.ProtectedBranchEvent{
		RepoId:    repo.Id,
		Name:      repo.Name,
		Path:      repo.Path,
		Operator:  operator,
		Action:    string(action),
		Before:    before,
		After:     after,
		EventTime: time.Now(),
	})
}

func (*outerImpl) ListProtectedBranch(ctx context.Context, reqDTO ListProtectedBranchReqDTO) ([]ProtectedBranchDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageProtectedBranchPerm(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	branchList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret, _ := listutil.Map(branchList, func(t branchmd.ProtectedBranch) (ProtectedBranchDTO, error) {
		return ProtectedBranchDTO{
			Id:      t.Id,
			Pattern: t.Pattern,
			Cfg:     *t.Cfg,
		}, nil
	})
	return ret, nil
}

// checkManageProtectedBranchPerm 检查权限
func checkManageProtectedBranchPerm(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, error) {
	// 检查仓库是否存在
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repo, util.InternalError(err)
	}
	if !b {
		return repo, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return repo, nil
	}
	// 如果不是 检查用户组权限
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		return repo, util.InternalError(err)
	}
	if b && (p.IsAdmin || p.PermDetail.GetRepoPerm(repoId).CanManageProtectedBranch) {
		return repo, nil
	}
	return repo, util.UnauthorizedError()
}
