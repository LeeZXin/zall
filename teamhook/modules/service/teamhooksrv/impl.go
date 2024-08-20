package teamhooksrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/notify/modules/model/notifymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/model/teamhookmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

// CreateTeamHook 创建team hook
func (*outerImpl) CreateTeamHook(ctx context.Context, reqDTO CreateTeamHookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	err := checkManageTeamHookByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	switch reqDTO.HookType {
	case teamhook.NotifyType:
		b, err := notifymd.ExistTplById(ctx, reqDTO.HookCfg.NotifyTplId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
	}
	err = teamhookmd.InsertTeamHook(ctx, teamhookmd.InsertTeamHookReqDTO{
		Name:     reqDTO.Name,
		TeamId:   reqDTO.TeamId,
		Events:   reqDTO.Events,
		HookType: reqDTO.HookType,
		HookCfg:  reqDTO.HookCfg,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateTeamHook 编辑team hook
func (*outerImpl) UpdateTeamHook(ctx context.Context, reqDTO UpdateTeamHookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	err := checkManageTeamHookByHookId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	switch reqDTO.HookType {
	case teamhook.NotifyType:
		b, err := notifymd.ExistTplById(ctx, reqDTO.HookCfg.NotifyTplId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		if !b {
			return util.InvalidArgsError()
		}
	}
	_, err = teamhookmd.UpdateTeamHook(ctx, teamhookmd.UpdateTeamHookReqDTO{
		Id:       reqDTO.Id,
		Name:     reqDTO.Name,
		Events:   reqDTO.Events,
		HookType: reqDTO.HookType,
		HookCfg:  reqDTO.HookCfg,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteTeamHook 删除team hook
func (*outerImpl) DeleteTeamHook(ctx context.Context, reqDTO DeleteTeamHookReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	err := checkManageTeamHookByHookId(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = teamhookmd.DeleteTeamHookById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListTeamHook team hook 列表
func (*outerImpl) ListTeamHook(ctx context.Context, reqDTO ListTeamHookReqDTO) ([]TeamHookDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 检查权限
	err := checkManageTeamHookByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	hooks, err := teamhookmd.ListTeamHookByTeamId(ctx, reqDTO.TeamId, []string{
		"id", "name", "team_id", "events", "hook_type", "hook_cfg",
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(hooks, func(t teamhookmd.TeamHook) (TeamHookDTO, error) {
		return TeamHookDTO{
			Id:       t.Id,
			Name:     t.Name,
			TeamId:   t.TeamId,
			Events:   t.GetEvents(),
			HookType: t.HookType,
			HookCfg:  t.GetHookCfg(),
		}, nil
	})
}

func checkManageTeamHookByHookId(ctx context.Context, hookId int64, operator apisession.UserInfo) error {
	hook, b, err := teamhookmd.GetTeamHookById(ctx, hookId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkManageTeamHookByTeamId(ctx, hook.TeamId, operator)
}

func checkManageTeamHookByTeamId(ctx context.Context, teamId int64, operator apisession.UserInfo) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManageNotifyTpl {
		return nil
	}
	return util.UnauthorizedError()
}
