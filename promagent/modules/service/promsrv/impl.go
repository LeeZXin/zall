package promsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
)

type outerImpl struct{}

// CreateScrape 创建抓取配置
func (*outerImpl) CreateScrape(ctx context.Context, reqDTO CreateScrapeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	b, err := appmd.ExistByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
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
		return util.InternalError(err)
	}
	return nil
}

// UpdateScrape 编辑抓取配置
func (*outerImpl) UpdateScrape(ctx context.Context, reqDTO UpdateScrapeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	_, err := prommd.UpdateScrapeById(ctx, prommd.UpdateScrapeByIdReqDTO{
		Id:         reqDTO.ScrapeId,
		Endpoint:   reqDTO.Endpoint,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListScrape 展示抓取列表
func (*outerImpl) ListScrape(ctx context.Context, reqDTO ListScrapeReqDTO) ([]ScrapeDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	scrapes, total, err := prommd.ListScrape(ctx, prommd.ListScrapeReqDTO{
		AppId:    reqDTO.AppId,
		Env:      reqDTO.Env,
		Endpoint: reqDTO.Endpoint,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var appIdNamMap = make(map[string]string)
	if len(scrapes) > 0 {
		appIdList, _ := listutil.Map(scrapes, func(t prommd.Scrape) (string, error) {
			return t.AppId, nil
		})
		apps, err := appmd.GetByAppIdList(ctx, listutil.Distinct(appIdList...))
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		for _, app := range apps {
			appIdNamMap[app.AppId] = app.Name
		}
	}

	data, _ := listutil.Map(scrapes, func(t prommd.Scrape) (ScrapeDTO, error) {
		return ScrapeDTO{
			Id:         t.Id,
			Endpoint:   t.Endpoint,
			AppId:      t.AppId,
			AppName:    appIdNamMap[t.AppId],
			Target:     t.Target,
			TargetType: t.TargetType,
			Created:    t.Created,
			Env:        t.Env,
		}, nil
	})
	return data, total, nil
}

// DeleteScrape 删除抓取配置
func (*outerImpl) DeleteScrape(ctx context.Context, reqDTO DeleteScrapeReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	_, err := prommd.DeleteScrapeById(ctx, reqDTO.ScrapeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}
