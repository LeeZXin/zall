package alertsrv

import (
	"context"
	"github.com/LeeZXin/zall/alert/modules/model/alertmd"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct{}

// InsertConfig 新增配置
func (*outerImpl) InsertConfig(ctx context.Context, reqDTO InsertConfigReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AlertSrvKeysVO.InsertConfig),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return
	}
	err = alertmd.InsertConfig(ctx, alertmd.InsertConfigReqDTO{
		Name:        reqDTO.Name,
		Alert:       reqDTO.Alert,
		AppId:       reqDTO.AppId,
		IntervalSec: reqDTO.IntervalSec,
		SilenceSec:  reqDTO.SilenceSec,
		Enabled:     reqDTO.Enabled,
		NextTime:    time.Now().Add(time.Duration(reqDTO.IntervalSec) * time.Second).UnixMilli(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// UpdateConfig 修改配置
func (*outerImpl) UpdateConfig(ctx context.Context, reqDTO UpdateConfigReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AlertSrvKeysVO.UpdateConfig),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator); err != nil {
		return
	}
	_, err = alertmd.UpdateConfig(ctx, alertmd.UpdateConfigReqDTO{
		Id:          reqDTO.Id,
		Name:        reqDTO.Name,
		Alert:       reqDTO.Alert,
		IntervalSec: reqDTO.IntervalSec,
		SilenceSec:  reqDTO.SilenceSec,
		Enabled:     reqDTO.Enabled,
		NextTime:    time.Now().Add(time.Duration(reqDTO.IntervalSec) * time.Second).UnixMilli(),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// DeleteConfig 删除配置
func (*outerImpl) DeleteConfig(ctx context.Context, reqDTO DeleteConfigReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.AlertSrvKeysVO.DeleteConfig),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPermByConfigId(ctx, reqDTO.Id, reqDTO.Operator); err != nil {
		return
	}
	_, err = alertmd.DeleteConfig(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListConfig 展示配置
func (*outerImpl) ListConfig(ctx context.Context, reqDTO ListConfigReqDTO) ([]ConfigDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, 0, err
	}
	configs, err := alertmd.ListConfig(ctx, alertmd.ListConfigReqDTO{
		Cursor: reqDTO.Cursor,
		Limit:  reqDTO.Limit,
		AppId:  reqDTO.AppId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var next int64 = 0
	if len(configs) == reqDTO.Limit {
		next = configs[len(configs)-1].Id
	}
	data, _ := listutil.Map(configs, func(t alertmd.Config) (ConfigDTO, error) {
		return ConfigDTO{
			Id:          t.Id,
			Name:        t.Name,
			AppId:       t.AppId,
			Content:     t.Content,
			IntervalSec: t.IntervalSec,
			SilenceSec:  t.SilenceSec,
			Enabled:     t.Enabled,
			NextTime:    t.NextTime,
			Created:     t.Created,
		}, nil
	})
	return data, next, nil
}

func checkPermByConfigId(ctx context.Context, id int64, operator apisession.UserInfo) error {
	cfg, b, err := alertmd.GetConfigById(ctx, id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkPermByAppId(ctx, cfg.AppId, operator)
}

func checkPermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) error {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	if p.PermDetail.DevelopAppList.Contains(appId) {
		return nil
	}
	return util.UnauthorizedError()
}
