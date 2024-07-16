package promsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

// CreateScrape 创建抓取配置
func (*outerImpl) CreateScrape(ctx context.Context, reqDTO CreateScrapeReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkManagePromAgentPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return
	}
	err = prommd.InsertScrape(ctx, prommd.InsertScrapeReqDTO{
		Endpoint:   reqDTO.Endpoint,
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

// UpdateScrape 编辑抓取配置
func (*outerImpl) UpdateScrape(ctx context.Context, reqDTO UpdateScrapeReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if err = checkManagePromAgentPermByScrapeId(ctx, reqDTO.ScrapeId, reqDTO.Operator); err != nil {
		return
	}
	_, err = prommd.UpdateScrapeById(ctx, prommd.UpdateScrapeByIdReqDTO{
		Id:         reqDTO.ScrapeId,
		Endpoint:   reqDTO.Endpoint,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListScrape 展示抓取列表
func (*outerImpl) ListScrape(ctx context.Context, reqDTO ListScrapeReqDTO) ([]ScrapeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkManagePromAgentPermByAppId(ctx, reqDTO.AppId, reqDTO.Operator); err != nil {
		return nil, err
	}
	scrapes, err := prommd.ListScrapeByAppIdAndEnv(ctx, reqDTO.AppId, reqDTO.Env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	data, _ := listutil.Map(scrapes, func(t prommd.Scrape) (ScrapeDTO, error) {
		return ScrapeDTO{
			Id:         t.Id,
			Endpoint:   t.Endpoint,
			AppId:      t.AppId,
			Target:     t.Target,
			TargetType: t.TargetType,
			Created:    t.Created,
			Env:        t.Env,
		}, nil
	})
	return data, nil
}

// DeleteScrape 删除抓取配置
func (*outerImpl) DeleteScrape(ctx context.Context, reqDTO DeleteScrapeReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkManagePromAgentPermByScrapeId(ctx, reqDTO.ScrapeId, reqDTO.Operator); err != nil {
		return
	}
	_, err = prommd.DeleteScrapeById(ctx, reqDTO.ScrapeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func checkManagePromAgentPermByAppId(ctx context.Context, appId string, operator apisession.UserInfo) error {
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
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManagePromAgent {
		return nil
	}
	return util.UnauthorizedError()
}

func checkManagePromAgentPermByScrapeId(ctx context.Context, scrapeId int64, operator apisession.UserInfo) error {
	scrape, b, err := prommd.GetScrapeById(ctx, scrapeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	return checkManagePromAgentPermByAppId(ctx, scrape.AppId, operator)
}
