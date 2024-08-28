package promsrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/appmd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"time"
)

var (
	initPsubOnce = sync.Once{}
)

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.AppPromScrapeTopic, func(data any) {
			req, ok := data.(event.AppPromScrapeEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.AppPromScrapeCreateAction:
							return cfg.AppPromScrape.Create
						case event.AppPromScrapeUpdateAction:
							return cfg.AppPromScrape.Update
						case event.AppPromScrapeDeleteAction:
							return cfg.AppPromScrape.Delete
						default:
							return false
						}
					}
					return false
				})
			}
		})
	})
}

// CreateScrapeBySa 创建抓取配置超级管理员
func CreateScrapeBySa(ctx context.Context, reqDTO CreateScrapeReqDTO) error {
	_, err := createScrape(ctx, reqDTO, func(_ context.Context, reqDTO CreateScrapeReqDTO) error {
		if reqDTO.Operator.IsAdmin {
			return nil
		}
		return util.UnauthorizedError()
	})
	return err
}

// CreateScrapeByTeam 创建抓取配置团队权限
func CreateScrapeByTeam(ctx context.Context, reqDTO CreateScrapeReqDTO) error {
	var (
		app  appmd.App
		team teammd.Team
	)
	scrape, err := createScrape(ctx, reqDTO, func(ctx context.Context, reqDTO CreateScrapeReqDTO) error {
		var err error
		app, team, err = checkManagePromScrapePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
		return err
	})
	if err != nil {
		return err
	}
	// 通知team hook
	notifyPromScrapeEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppPromScrapeCreateAction,
		scrape,
	)
	return nil
}

func createScrape(ctx context.Context, reqDTO CreateScrapeReqDTO, checkPerm func(context.Context, CreateScrapeReqDTO) error) (prommd.Scrape, error) {
	if err := reqDTO.IsValid(); err != nil {
		return prommd.Scrape{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO)
	if err != nil {
		return prommd.Scrape{}, err
	}
	b, err := appmd.ExistByAppId(ctx, reqDTO.AppId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return prommd.Scrape{}, util.InternalError(err)
	}
	if !b {
		return prommd.Scrape{}, util.InvalidArgsError()
	}
	ret, err := prommd.InsertScrape(ctx, prommd.InsertScrapeReqDTO{
		Endpoint:   reqDTO.Endpoint,
		AppId:      reqDTO.AppId,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Env:        reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return prommd.Scrape{}, util.InternalError(err)
	}
	return ret, nil
}

// UpdateScrapeBySa 编辑抓取配置
func UpdateScrapeBySa(ctx context.Context, reqDTO UpdateScrapeReqDTO) error {
	return updateScrape(ctx, reqDTO, func(_ context.Context, reqDTO UpdateScrapeReqDTO) error {
		if reqDTO.Operator.IsAdmin {
			return nil
		}
		return util.UnauthorizedError()
	})
}

// UpdateScrapeByTeam 编辑抓取配置
func UpdateScrapeByTeam(ctx context.Context, reqDTO UpdateScrapeReqDTO) error {
	var (
		app    appmd.App
		team   teammd.Team
		scrape prommd.Scrape
	)
	err := updateScrape(ctx, reqDTO, func(ctx context.Context, reqDTO UpdateScrapeReqDTO) error {
		var err error
		scrape, app, team, err = checkManagePromScrapePermByScrapeId(ctx, reqDTO.Operator, reqDTO.Id)
		return err
	})
	if err != nil {
		return err
	}
	// 通知team hook
	scrape.Endpoint = reqDTO.Endpoint
	scrape.Target = reqDTO.Target
	scrape.TargetType = reqDTO.TargetType
	notifyPromScrapeEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppPromScrapeUpdateAction,
		scrape,
	)
	return nil
}

func updateScrape(ctx context.Context, reqDTO UpdateScrapeReqDTO, checkPerm func(context.Context, UpdateScrapeReqDTO) error) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkPerm(ctx, reqDTO)
	if err != nil {
		return err
	}
	_, err = prommd.UpdateScrapeById(ctx, prommd.UpdateScrapeByIdReqDTO{
		Id:         reqDTO.Id,
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

// ListScrapeBySa 展示抓取列表
func ListScrapeBySa(ctx context.Context, reqDTO ListScrapeReqDTO) ([]ScrapeBySaDTO, int64, error) {
	if err := reqDTO.IsValidBySa(); err != nil {
		return nil, 0, err
	}
	if !reqDTO.Operator.IsAdmin {
		return nil, 0, util.UnauthorizedError()
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
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
	appIdMap := make(map[string]appmd.App)
	teamIdNameMap := make(map[int64]string)
	if len(scrapes) > 0 {
		appIdList := listutil.MapNe(scrapes, func(t prommd.Scrape) string {
			return t.AppId
		})
		apps, err := appmd.GetByAppIdList(ctx, listutil.Distinct(appIdList...), []string{"app_id", "name", "team_id"})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		teamIdList := make([]int64, 0)
		for _, app := range apps {
			appIdMap[app.AppId] = app
			teamIdList = append(teamIdList, app.TeamId)
		}
		teams, err := teammd.ListTeamByIdList(ctx, listutil.Distinct(teamIdList...), []string{"id", "name"})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		for _, team := range teams {
			teamIdNameMap[team.Id] = team.Name
		}
	}
	data := listutil.MapNe(scrapes, func(t prommd.Scrape) ScrapeBySaDTO {
		teamId := appIdMap[t.AppId].TeamId
		return ScrapeBySaDTO{
			Id:         t.Id,
			Endpoint:   t.Endpoint,
			AppId:      t.AppId,
			AppName:    appIdMap[t.AppId].Name,
			TeamId:     teamId,
			TeamName:   teamIdNameMap[teamId],
			Target:     t.Target,
			TargetType: t.TargetType,
			Env:        t.Env,
		}
	})
	return data, total, nil
}

// ListScrapeByTeam 展示抓取列表
func ListScrapeByTeam(ctx context.Context, reqDTO ListScrapeReqDTO) ([]ScrapeByTeamDTO, int64, error) {
	if err := reqDTO.IsValidByTeam(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, _, err := checkManagePromScrapePermByAppId(ctx, reqDTO.Operator, reqDTO.AppId)
	if err != nil {
		return nil, 0, err
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
	data := listutil.MapNe(scrapes, func(t prommd.Scrape) ScrapeByTeamDTO {
		return ScrapeByTeamDTO{
			Id:         t.Id,
			Endpoint:   t.Endpoint,
			Target:     t.Target,
			TargetType: t.TargetType,
			Env:        t.Env,
		}
	})
	return data, total, nil
}

// DeleteScrapeBySa 删除抓取配置
func DeleteScrapeBySa(ctx context.Context, reqDTO DeleteScrapeReqDTO) error {
	return deleteScrape(ctx, reqDTO, func(_ context.Context, reqDTO DeleteScrapeReqDTO) error {
		if reqDTO.Operator.IsAdmin {
			return nil
		}
		return util.UnauthorizedError()
	})
}

// DeleteScrapeByTeam 删除抓取配置
func DeleteScrapeByTeam(ctx context.Context, reqDTO DeleteScrapeReqDTO) error {
	var (
		app    appmd.App
		team   teammd.Team
		scrape prommd.Scrape
	)
	err := deleteScrape(ctx, reqDTO, func(ctx context.Context, reqDTO DeleteScrapeReqDTO) error {
		var err error
		scrape, app, team, err = checkManagePromScrapePermByScrapeId(ctx, reqDTO.Operator, reqDTO.Id)
		return err
	})
	if err != nil {
		return err
	}
	// 通知team hook
	notifyPromScrapeEvent(
		reqDTO.Operator,
		team,
		app,
		event.AppPromScrapeDeleteAction,
		scrape,
	)
	return nil
}

func deleteScrape(ctx context.Context, reqDTO DeleteScrapeReqDTO, checkPerm func(context.Context, DeleteScrapeReqDTO) error) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	err := checkPerm(ctx, reqDTO)
	if err != nil {
		return err
	}
	_, err = prommd.DeleteScrapeById(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkManagePromScrapePermByAppId(ctx context.Context, operator apisession.UserInfo, appId string) (appmd.App, teammd.Team, error) {
	app, b, err := appmd.GetByAppId(ctx, appId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, app.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return appmd.App{}, teammd.Team{}, util.ThereHasBugErr()
	}
	if operator.IsAdmin {
		return app, team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, app.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return app, team, util.InternalError(err)
	}
	if !b {
		return app, team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManagePromScrape {
		return app, team, nil
	}
	return app, team, util.UnauthorizedError()
}

func checkManagePromScrapePermByScrapeId(ctx context.Context, operator apisession.UserInfo, scrapeId int64) (prommd.Scrape, appmd.App, teammd.Team, error) {
	scrape, b, err := prommd.GetScrapeById(ctx, scrapeId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return prommd.Scrape{}, appmd.App{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return prommd.Scrape{}, appmd.App{}, teammd.Team{}, util.InvalidArgsError()
	}
	app, team, err := checkManagePromScrapePermByAppId(ctx, operator, scrape.AppId)
	return scrape, app, team, err
}

func notifyPromScrapeEvent(operator apisession.UserInfo, team teammd.Team, app appmd.App, action event.AppPromScrapeEventAction, scrape prommd.Scrape) {
	initPsub()
	psub.Publish(event.AppPromScrapeTopic, event.AppPromScrapeEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseApp: event.BaseApp{
			AppId:   app.AppId,
			AppName: app.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action:   action,
		Endpoint: scrape.Endpoint,
		Env:      scrape.Env,
	})
}
