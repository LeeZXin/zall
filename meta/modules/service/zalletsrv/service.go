package zalletsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/zalletmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

// CreateZalletNode 创建node
func CreateZalletNode(ctx context.Context, reqDTO CreateZalletNodeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	// 检查权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := zalletmd.GetZalletNodeByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = zalletmd.InsertZalletNode(ctx, zalletmd.InsertZalletNodeReqDTO{
		NodeId:     reqDTO.NodeId,
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
func UpdateZalletNode(ctx context.Context, reqDTO UpdateZalletNodeReqDTO) error {
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
		Id:         reqDTO.Id,
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
func DeleteZalletNode(ctx context.Context, reqDTO DeleteZalletNodeReqDTO) error {
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
func ListZalletNode(ctx context.Context, reqDTO ListZalletNodeReqDTO) ([]ZalletNodeDTO, int64, error) {
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
		Cols:     []string{"id", "node_id", "name", "agent_host", "agent_token"},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data := listutil.MapNe(ret, func(t zalletmd.ZalletNode) ZalletNodeDTO {
		return ZalletNodeDTO{
			Id:         t.Id,
			NodeId:     t.NodeId,
			Name:       t.Name,
			AgentHost:  t.AgentHost,
			AgentToken: t.AgentToken,
		}
	})
	return data, total, nil
}

// ListAllZalletNode 所有列表
func ListAllZalletNode(ctx context.Context, reqDTO ListAllZalletNodeReqDTO) ([]SimpleZalletNodeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	nodes, err := zalletmd.ListAllZalletNode(ctx, []string{"id", "name", "node_id"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.MapNe(nodes, func(t zalletmd.ZalletNode) SimpleZalletNodeDTO {
		return SimpleZalletNodeDTO{
			Id:     t.Id,
			NodeId: t.NodeId,
			Name:   t.Name,
		}
	}), nil
}
