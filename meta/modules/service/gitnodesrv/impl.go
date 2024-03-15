package gitnodesrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) InsertNode(ctx context.Context, reqDTO InsertNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GitNodeSrvKeysVO.InsertNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, b, err := gitnodemd.GetByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = gitnodemd.InsertGitNode(ctx, gitnodemd.InsertNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		HttpHosts: reqDTO.HttpHosts,
		SshHosts:  reqDTO.SshHosts,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteNode(ctx context.Context, reqDTO DeleteNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GitNodeSrvKeysVO.DeleteNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	_, b, err := gitnodemd.GetByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = gitnodemd.DeleteNode(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateNode(ctx context.Context, reqDTO UpdateNodeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.GitNodeSrvKeysVO.UpdateNode),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, b, err := gitnodemd.GetByNodeId(ctx, reqDTO.NodeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = gitnodemd.UpdateNode(ctx, gitnodemd.UpdateNodeReqDTO{
		NodeId:    reqDTO.NodeId,
		HttpHosts: reqDTO.HttpHosts,
		SshHosts:  reqDTO.SshHosts,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListNode(ctx context.Context, reqDTO ListNodeReqDTO) ([]GitNodeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	if !reqDTO.Operator.IsAdmin {
		return nil, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	all, err := gitnodemd.GetAll(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(all, func(t gitnodemd.GitNodeDTO) (GitNodeDTO, error) {
		return GitNodeDTO{
			NodeId:    t.NodeId,
			HttpHosts: t.HttpHosts,
			SshHosts:  t.SshHosts,
		}, nil
	})
}
