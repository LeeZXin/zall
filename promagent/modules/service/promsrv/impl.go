package promsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

func (*outerImpl) InsertScrape(ctx context.Context, reqDTO InsertScrapeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PromSrvKeysVO.InsertScrape),
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
	err = prommd.InsertScrape(ctx, prommd.InsertScrapeReqDTO{
		ServerUrl:  reqDTO.ServerUrl,
		AppId:      reqDTO.AppId,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Env:        reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) UpdateScrape(ctx context.Context, reqDTO UpdateScrapeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PromSrvKeysVO.UpdateScrape),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPermByScrapeId(ctx, reqDTO.Id, reqDTO.Env, reqDTO.Operator); err != nil {
		return
	}
	_, err = prommd.UpdateScrapeById(ctx, prommd.UpdateScrapeByIdReqDTO{
		Id:         reqDTO.Id,
		ServerUrl:  reqDTO.ServerUrl,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Env:        reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListScrape(ctx context.Context, reqDTO ListScrapeReqDTO) ([]ScrapeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	scrapes, err := prommd.GetAllScrapeByAppId(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(scrapes, func(t prommd.Scrape) (ScrapeDTO, error) {
		return ScrapeDTO{
			Id:         t.Id,
			ServerUrl:  t.ServerUrl,
			AppId:      t.AppId,
			Target:     t.Target,
			TargetType: t.TargetType,
			Created:    t.Created,
		}, nil
	})
}

func (*outerImpl) DeleteScrape(ctx context.Context, reqDTO DeleteScrapeReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.PromSrvKeysVO.DeleteScrape),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPermByScrapeId(ctx, reqDTO.Id, reqDTO.Env, reqDTO.Operator); err != nil {
		return
	}
	_, err = prommd.DeleteById(ctx, reqDTO.Id, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
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
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}

func checkPermByScrapeId(ctx context.Context, scrapeId int64, env string, operator apisession.UserInfo) error {
	scrape, b, err := prommd.GetById(ctx, scrapeId, env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	app, b, err := appmd.GetByAppId(ctx, scrape.AppId)
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
	if !b || !p.IsAdmin {
		return util.UnauthorizedError()
	}
	return nil
}
