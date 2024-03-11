package gitnodesrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) InsertNode(ctx context.Context, reqDTO InsertNodeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, b, err := gitnodemd.GetByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b {
		return util.AlreadyExistsError()
	}
	err = gitnodemd.InsertGitNode(ctx, gitnodemd.InsertNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		HttpHosts: reqDTO.HttpHosts,
		SshHosts:  reqDTO.SshHosts,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (*outerImpl) DeleteNode(ctx context.Context, reqDTO DeleteNodeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, b, err := gitnodemd.GetByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = gitnodemd.DeleteNode(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}
