package zalletsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/zalletmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

// CreateZalletNode 创建node
func (*outerImpl) CreateZalletNode(ctx context.Context, reqDTO CreateZalletNodeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := zalletmd.InsertZalletNode(ctx, zalletmd.InsertZalletNodeReqDTO{
		Name:       reqDTO.Name,
		AgentHost:  reqDTO.AgentHost,
		AgentToken: reqDTO.AgentToken,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// UpdateZalletNode 编辑node
func (*outerImpl) UpdateZalletNode(ctx context.Context, reqDTO UpdateZalletNodeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := zalletmd.UpdateZalletNode(ctx, zalletmd.UpdateZalletNodeReqDTO{
		Id:         reqDTO.NodeId,
		Name:       reqDTO.Name,
		AgentHost:  reqDTO.AgentHost,
		AgentToken: reqDTO.AgentToken,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteZalletNode 删除node
func (*outerImpl) DeleteZalletNode(ctx context.Context, reqDTO DeleteZalletNodeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := zalletmd.DeleteZalletNodeById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListZalletNode node列表
func (*outerImpl) ListZalletNode(ctx context.Context, reqDTO ListZalletNodeReqDTO) ([]ZalletNodeDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	ret, total, err := zalletmd.ListZalletNode(ctx, zalletmd.ListZalletNodeReqDTO{
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		Name:     reqDTO.Name,
		Cols:     []string{"id", "name", "agent_host", "agent_token"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(ret, func(t zalletmd.ZalletNode) (ZalletNodeDTO, error) {
		return ZalletNodeDTO{
			Id:         t.Id,
			Name:       t.Name,
			AgentHost:  t.AgentHost,
			AgentToken: t.AgentToken,
		}, nil
	})
	return data, total, nil
}

// ListAllZalletNode 所有列表
func (*outerImpl) ListAllZalletNode(ctx context.Context, reqDTO ListAllZalletNodeReqDTO) ([]SimpleZalletNodeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := zalletmd.ListAllZalletNode(ctx, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(nodes, func(t zalletmd.ZalletNode) (SimpleZalletNodeDTO, error) {
		return SimpleZalletNodeDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
	return data, nil
}
