package timersrv

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/timer/modules/model/timermd"
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
		psub.Subscribe(event.TimerTopic, func(data any) {
			req, ok := data.(event.TimerEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.TimerCreateAction:
							return cfg.Timer.Create
						case event.TimerUpdateAction:
							return cfg.Timer.Update
						case event.TimerDeleteAction:
							return cfg.Timer.Delete
						case event.TimerEnableAction:
							return cfg.Timer.Enable
						case event.TimerDisableAction:
							return cfg.Timer.Disable
						case event.TimerManualTriggerAction:
							return cfg.Timer.ManuallyTrigger
						default:
							return false
						}
					}
					return false
				})
			}
		})
		psub.Subscribe(event.TimerTaskTopic, func(data any) {
			req, ok := data.(event.TimerTaskEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					if events.EnvRelated == nil {
						return false
					}
					cfg, ok := events.EnvRelated[req.Env]
					if ok {
						switch req.Action {
						case event.TimerTaskFailAction:
							return cfg.TimerTask.Fail
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

// CreateTimer 新增任务
func CreateTimer(ctx context.Context, reqDTO CreateTimerReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	team, err := checkManageTimerPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return err
	}
	var timer timermd.Timer
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		timer, err2 = timermd.InsertTimer(ctx, timermd.InsertTimerReqDTO{
			Name:      reqDTO.Name,
			CronExp:   reqDTO.CronExp,
			Content:   reqDTO.Task,
			TeamId:    reqDTO.TeamId,
			Env:       reqDTO.Env,
			IsEnabled: false,
			Creator:   reqDTO.Operator.Account,
		})
		if err2 != nil {
			return err2
		}
		return timermd.InsertExecute(ctx, timermd.InsertExecuteReqDTO{
			TimerId:   timer.Id,
			IsEnabled: false,
			Env:       reqDTO.Env,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTimerEvent(
		reqDTO.Operator,
		team,
		timer,
		event.TimerCreateAction,
	)
	return nil
}

// ListTimer 展示任务列表
func ListTimer(ctx context.Context, reqDTO ListTimerReqDTO) ([]TimerDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageTimerPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId)
	if err != nil {
		return nil, 0, err
	}
	timers, total, err := timermd.ListTimer(ctx, timermd.ListTimerReqDTO{
		TeamId:   reqDTO.TeamId,
		Name:     reqDTO.Name,
		PageNum:  reqDTO.PageNum,
		PageSize: 10,
		Env:      reqDTO.Env,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret := listutil.MapNe(timers, func(t timermd.Timer) TimerDTO {
		return TimerDTO{
			Id:        t.Id,
			Name:      t.Name,
			CronExp:   t.CronExp,
			Task:      t.GetContent(),
			TeamId:    t.TeamId,
			IsEnabled: t.IsEnabled,
			Env:       t.Env,
			Creator:   t.Creator,
		}
	})
	return ret, total, nil
}

// EnableTimer 启动任务
func EnableTimer(ctx context.Context, reqDTO EnableTimerReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	timer, team, err := checkManageTimerPermByTimerId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	cron, err := ParseCron(timer.CronExp)
	if err != nil {
		return util.ThereHasBugErr()
	}
	now := time.Now()
	nextTime := cron.Next(now)
	if nextTime.Before(now) {
		return util.NewBizErr(apicode.OperationFailedErrCode, i18n.CronExpError)
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err2 := timermd.EnableExecute(ctx, reqDTO.Id, nextTime.UnixMilli())
		if err2 != nil {
			return err2
		}
		if b {
			b, err2 = timermd.EnableTimer(ctx, reqDTO.Id)
			if err2 != nil {
				return err2
			}
			if !b {
				return errors.New("failed")
			}
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTimerEvent(
		reqDTO.Operator,
		team,
		timer,
		event.TimerEnableAction,
	)
	return nil
}

// DisableTimer 关闭任务
func DisableTimer(ctx context.Context, reqDTO DisableTimerReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	timer, team, err := checkManageTimerPermByTimerId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		b, err2 := timermd.DisableExecute(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		if b {
			b, err2 = timermd.DisableTimer(ctx, reqDTO.Id)
			if err2 != nil {
				return err2
			}
			if !b {
				return errors.New("failed")
			}
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTimerEvent(
		reqDTO.Operator,
		team,
		timer,
		event.TimerDisableAction,
	)
	return nil
}

// DeleteTimer 删除任务
func DeleteTimer(ctx context.Context, reqDTO DeleteTimerReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	timer, team, err := checkManageTimerPermByTimerId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err2 := timermd.DeleteTimer(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		_, err2 = timermd.DeleteExecuteByTimerId(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		err2 = timermd.DeleteLogByTimerId(ctx, reqDTO.Id)
		return err2
	})
	if err != nil {
		return util.InternalError(err)
	}
	notifyTimerEvent(
		reqDTO.Operator,
		team,
		timer,
		event.TimerDeleteAction,
	)
	return nil
}

// TriggerTask 手动执行任务
func TriggerTask(ctx context.Context, reqDTO TriggerTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		timer timermd.Timer
	)
	timer, team, err := checkManageTimerPermByTimerId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	// 手动触发通知
	notifyTimerEvent(
		reqDTO.Operator,
		team,
		timer,
		event.TimerManualTriggerAction,
	)
	err = triggerTask(&timer, reqDTO.Operator.Account, reqDTO.Operator.Name)
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

// ListLog 获取执行历史
func ListLog(ctx context.Context, reqDTO ListLogReqDTO) ([]LogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, _, err := checkManageTimerPermByTimerId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return nil, 0, err
	}
	// 当月第一天
	beginTime := reqDTO.monthTime.AddDate(0, 0, -reqDTO.monthTime.Day()+1)
	// 当月最后一天
	endTime := reqDTO.monthTime.AddDate(0, 1, -1)
	logs, total, err := timermd.ListLog(ctx, timermd.ListLogReqDTO{
		TimerId:   reqDTO.Id,
		PageNum:   reqDTO.PageNum,
		PageSize:  10,
		BeginTime: beginTime,
		EndTime:   time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, endTime.Location()),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret := listutil.MapNe(logs, func(t timermd.Log) LogDTO {
		return LogDTO{
			Task:        t.GetTaskContent(),
			ErrLog:      t.ErrLog,
			TriggerType: t.TriggerType,
			TriggerBy:   t.TriggerBy,
			IsSuccess:   t.IsSuccess,
			Created:     t.Created,
		}
	})
	return ret, total, nil
}

// UpdateTimer 更新任务
func UpdateTimer(ctx context.Context, reqDTO UpdateTimerReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	timer, team, err := checkManageTimerPermByTimerId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = timermd.UpdateTimer(ctx, timermd.UpdateTimerReqDTO{
		Id:      reqDTO.Id,
		Name:    reqDTO.Name,
		CronExp: reqDTO.CronExp,
		Content: reqDTO.Task,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyTimerEvent(
		reqDTO.Operator,
		team,
		timer,
		event.TimerUpdateAction,
	)
	return nil
}

func checkManageTimerPermByTeamId(ctx context.Context, operator apisession.UserInfo, teamId int64) (teammd.Team, error) {
	team, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Team{}, util.UnauthorizedError()
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
	if !p.IsAdmin && !p.PermDetail.TeamPerm.CanManageTimer {
		return team, util.UnauthorizedError()
	}
	return team, nil
}

func checkManageTimerPermByTimerId(ctx context.Context, operator apisession.UserInfo, timerId int64) (timermd.Timer, teammd.Team, error) {
	timer, b, err := timermd.GetTimerById(ctx, timerId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return timermd.Timer{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return timermd.Timer{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkManageTimerPermByTeamId(ctx, operator, timer.TeamId)
	return timer, team, err
}

func notifyTimerEvent(operator apisession.UserInfo, team teammd.Team, timer timermd.Timer, action event.TimerEventAction) {
	initPsub()
	psub.Publish(event.TimerTopic, event.TimerEvent{
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
		BaseTimer: event.BaseTimer{
			TimerId:   timer.Id,
			TimerName: timer.Name,
			Env:       timer.Env,
		},
		Action: action,
	})
}

func notifyTimerTaskEvent(operatorAccount, operatorName string, team teammd.Team, timer timermd.Timer, triggerType, taskStatus string, action event.TimerTaskEventAction) {
	initPsub()
	psub.Publish(event.TimerTaskTopic, event.TimerTaskEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operatorAccount,
			OperatorName: operatorName,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		BaseTimer: event.BaseTimer{
			TimerId:   timer.Id,
			TimerName: timer.Name,
			Env:       timer.Env,
		},
		Action:      action,
		TriggerType: triggerType,
		TaskStatus:  taskStatus,
	})
}
