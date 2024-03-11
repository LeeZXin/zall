package detectsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/tcpdetect/modules/model/detectmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) InsertDetect(ctx context.Context, reqDTO InsertDetectReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TcpDetectSrvKeysVO.InsertDetect),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var b bool
	_, b, err = detectmd.GetDetectByIpPort(ctx, reqDTO.Ip, reqDTO.Port)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b {
		err = util.AlreadyExistsError()
		return
	}
	err = detectmd.InsertDetect(ctx, detectmd.InsertDetectReqDTO{
		Ip:            reqDTO.Ip,
		Port:          reqDTO.Port,
		Name:          reqDTO.Name,
		Enabled:       false,
		HeartbeatTime: 0,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateDetect(ctx context.Context, reqDTO UpdateDetectReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TcpDetectSrvKeysVO.UpdateDetect),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		detect detectmd.TcpDetect
		b      bool
	)
	detect, b, err = detectmd.GetDetectByIpPort(ctx, reqDTO.Ip, reqDTO.Port)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if b && detect.Id != reqDTO.Id {
		return util.AlreadyExistsError()
	}
	_, err = detectmd.UpdateDetect(ctx, detectmd.UpdateDetectReqDTO{
		Id:   reqDTO.Id,
		Ip:   reqDTO.Ip,
		Port: reqDTO.Port,
		Name: reqDTO.Name,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListDetect(ctx context.Context, reqDTO ListDetectReqDTO) ([]DetectDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	detects, err := detectmd.ListDetect(ctx, detectmd.ListDetectReqDTO{
		Name:   reqDTO.Name,
		Cursor: reqDTO.Cursor,
		Limit:  reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var cursor int64 = 0
	if reqDTO.Limit > 0 && len(detects) == reqDTO.Limit {
		cursor = detects[len(detects)-1].Id
	}
	ret, _ := listutil.Map(detects, func(t detectmd.TcpDetect) (DetectDTO, error) {
		return DetectDTO{
			Id:            t.Id,
			Ip:            t.Ip,
			Port:          t.Port,
			Name:          t.Name,
			HeartbeatTime: t.HeartbeatTime,
			Enabled:       t.Enabled,
		}, nil
	})
	return ret, cursor, nil
}

func (*outerImpl) DeleteDetect(ctx context.Context, reqDTO DeleteDetectReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TcpDetectSrvKeysVO.DeleteDetect),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err = detectmd.DeleteDetect(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	detectmd.DeleteLog(ctx, reqDTO.Id)
	return
}

func (*outerImpl) ListLog(ctx context.Context, reqDTO ListLogReqDTO) ([]LogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	logs, err := detectmd.ListLog(ctx, detectmd.ListLogReqDTO{
		Id:     reqDTO.Id,
		Cursor: reqDTO.Cursor,
		Limit:  reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var cursor int64 = 0
	if reqDTO.Limit > 0 && len(logs) == reqDTO.Limit {
		cursor = logs[len(logs)-1].Id
	}
	ret, _ := listutil.Map(logs, func(t detectmd.DetectLog) (LogDTO, error) {
		return LogDTO{
			Ip:      t.Ip,
			Port:    t.Port,
			Valid:   t.Valid,
			Created: t.Created,
		}, nil
	})
	return ret, cursor, nil
}

func (*outerImpl) EnabledDetect(ctx context.Context, reqDTO EnableDetectReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TcpDetectSrvKeysVO.EnableDetect),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = detectmd.SetDetectEnabled(ctx, reqDTO.Id, true)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DisableDetect(ctx context.Context, reqDTO DisableDetectReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TcpDetectSrvKeysVO.DisableDetect),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = detectmd.SetDetectEnabled(ctx, reqDTO.Id, false)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}
