package notifysrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/notify/modules/model/notifymd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/notify/notify"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
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
		psub.Subscribe(event.NotifyTplTopic, func(data any) {
			req, ok := data.(event.NotifyTplEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.NotifyTplCreateAction:
						return events.NotifyTpl.Create
					case event.NotifyTplUpdateAction:
						return events.NotifyTpl.Update
					case event.NotifyTplDeleteAction:
						return events.NotifyTpl.Delete
					case event.NotifyTplChangeApiKeyAction:
						return events.NotifyTpl.ChangeApiKey
					default:
						return false
					}
				})
			}
		})
	})
}

// SendNotificationByTplId 通过模板id发送通知
func SendNotificationByTplId(tplId int64, params any) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	tpl, b, err := notifymd.GetTplById(ctx, tplId, nil)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if !b {
		return
	}
	err = notify.SendNotification(tpl.GetNotifyCfg(), params)
	if err != nil {
		logger.Logger.Errorf("send tplId: %d failed with error: %v", tplId, err)
	}
}

// CreateTpl 创建通知模板
func CreateTpl(ctx context.Context, reqDTO CreateTplReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	team, err := checkManageNotifyTplPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return err
	}
	tpl, err := notifymd.InsertTpl(ctx, notifymd.InsertTplReqDTO{
		Name:      reqDTO.Name,
		ApiKey:    idutil.RandomUuid(),
		NotifyCfg: reqDTO.Cfg,
		TeamId:    reqDTO.TeamId,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, tpl, reqDTO.Cfg, reqDTO.Operator, event.NotifyTplCreateAction)
	return nil
}

// UpdateTpl 编辑通知模板
func UpdateTpl(ctx context.Context, reqDTO UpdateTplReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	tpl, team, err := checkManageNotifyTplPermByTplId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = notifymd.UpdateTpl(ctx, notifymd.UpdateTplReqDTO{
		Id:        reqDTO.Id,
		Name:      reqDTO.Name,
		NotifyCfg: reqDTO.Cfg,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, tpl, reqDTO.Cfg, reqDTO.Operator, event.NotifyTplUpdateAction)
	return nil
}

// DeleteTpl 删除通知模板
func DeleteTpl(ctx context.Context, reqDTO DeleteTplReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	tpl, team, err := checkManageNotifyTplPermByTplId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = notifymd.DeleteTpl(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, tpl, tpl.GetNotifyCfg(), reqDTO.Operator, event.NotifyTplDeleteAction)
	return nil
}

// ListTpl 通知模板列表
func ListTpl(ctx context.Context, reqDTO ListTplReqDTO) ([]TplDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageNotifyTplPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, 0, err
	}
	tpls, total, err := notifymd.ListTpl(ctx, notifymd.ListTplReqDTO{
		Name:     reqDTO.Name,
		PageNum:  reqDTO.PageNum,
		TeamId:   reqDTO.TeamId,
		PageSize: 10,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret := listutil.MapNe(tpls, func(t notifymd.Tpl) TplDTO {
		ret := TplDTO{
			Id:     t.Id,
			Name:   t.Name,
			ApiKey: t.ApiKey,
			TeamId: t.TeamId,
		}
		if t.NotifyCfg != nil {
			ret.NotifyCfg = t.NotifyCfg.Data
		}
		return ret
	})
	return ret, total, nil
}

// ChangeTplApiKey 更换api key
func ChangeTplApiKey(ctx context.Context, reqDTO ChangeTplApiKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	tpl, team, err := checkManageNotifyTplPermByTplId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = notifymd.UpdateTplApiKeyById(ctx, reqDTO.Id, idutil.RandomUuid())
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(team, tpl, tpl.GetNotifyCfg(), reqDTO.Operator, event.NotifyTplChangeApiKeyAction)
	return nil
}

// SendNotificationByApiKey 通过api key发送通知
func SendNotificationByApiKey(ctx context.Context, reqDTO SendNotifyByApiKeyReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	tpl, b, err := notifymd.GetTplByApiKey(ctx, reqDTO.ApiKey)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	err = notify.SendNotification(tpl.GetNotifyCfg(), reqDTO.Params)
	if err != nil {
		logger.Logger.WithContext(ctx).Errorf("send tpl id: %d failed with error: %v", tpl.Id, err)
		return util.OperationFailedError()
	}
	return nil
}

// ListAllTpl 通过团队获取模板列表
func ListAllTpl(ctx context.Context, reqDTO ListAllTplReqDTO) ([]SimpleTplDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 只要属于同一个团队的都可以看到数据
	if !reqDTO.Operator.IsAdmin {
		_, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.UnauthorizedError()
		}
	}
	tpls, err := notifymd.ListAllTplByTeamId(ctx, reqDTO.TeamId, []string{"id", "name"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tpls, func(t notifymd.Tpl) (SimpleTplDTO, error) {
		return SimpleTplDTO{
			Id:   t.Id,
			Name: t.Name,
		}, nil
	})
}

func checkManageNotifyTplPermByTeamId(ctx context.Context, operator apisession.UserInfo, teamId int64) (teammd.Team, error) {
	team, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Team{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return team, util.InternalError(err)
	}
	if !b {
		return team, util.UnauthorizedError()
	}
	if p.IsAdmin || p.PermDetail.TeamPerm.CanManageNotifyTpl {
		return team, nil
	}
	return team, util.UnauthorizedError()
}

func checkManageNotifyTplPermByTplId(ctx context.Context, operator apisession.UserInfo, tplId int64) (notifymd.Tpl, teammd.Team, error) {
	tpl, b, err := notifymd.GetTplById(ctx, tplId, []string{"team_id"})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return notifymd.Tpl{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return notifymd.Tpl{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkManageNotifyTplPermByTeamId(ctx, operator, tpl.TeamId)
	return tpl, team, err
}

func notifyEvent(team teammd.Team, tpl notifymd.Tpl, cfg notify.Cfg, operator apisession.UserInfo, action event.NotifyTplEventAction) {
	initPsub()
	psub.Publish(event.NotifyTplTopic, event.NotifyTplEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action: action,
		Name:   tpl.Name,
		Type:   cfg.NotifyType,
	})
}
