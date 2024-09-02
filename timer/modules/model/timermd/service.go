package timermd

import (
	"context"
	"github.com/LeeZXin/zall/pkg/timer"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsTimerNameValid(name string) bool {
	return len(name) <= 64 && len(name) > 0
}

func InsertTimer(ctx context.Context, reqDTO InsertTimerReqDTO) (Timer, error) {
	ret := Timer{
		Name:    reqDTO.Name,
		CronExp: reqDTO.CronExp,
		Content: &xormutil.Conversion[timer.Task]{
			Data: reqDTO.Content,
		},
		TeamId:    reqDTO.TeamId,
		Env:       reqDTO.Env,
		IsEnabled: reqDTO.IsEnabled,
		Creator:   reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTimer(ctx context.Context, reqDTO UpdateTimerReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("name", "cron_exp", "content").
		Update(&Timer{
			Name:    reqDTO.Name,
			CronExp: reqDTO.CronExp,
			Content: &xormutil.Conversion[timer.Task]{
				Data: reqDTO.Content,
			},
		})
	return rows == 1, err
}

func UpdateExecuteNextTime(ctx context.Context, timerId, nextTime, runVersion int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("timer_id = ?", timerId).
		And("run_version = ?", runVersion).
		Cols("next_time", "run_version").
		Update(&Execute{
			NextTime:   nextTime,
			RunVersion: runVersion + 1,
		})
	return rows == 1, err
}

func EnableExecute(ctx context.Context, timerId int64, nextTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("timer_id = ?", timerId).
		And("is_enabled = 0").
		Cols("is_enabled", "next_time").
		Update(&Execute{
			IsEnabled: true,
			NextTime:  nextTime,
		})
	return rows == 1, err
}

func EnableTimer(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("is_enabled = 0").
		Cols("is_enabled").
		Update(&Timer{
			IsEnabled: true,
		})
	return rows == 1, err
}

func DisableExecute(ctx context.Context, timerId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("timer_id = ?", timerId).
		And("is_enabled = 1").
		Cols("is_enabled", "next_time").
		Update(&Execute{
			IsEnabled: false,
			NextTime:  0,
		})
	return rows == 1, err
}

func DisableTimer(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("is_enabled = 1").
		Cols("is_enabled").
		Update(&Timer{
			IsEnabled: false,
		})
	return rows == 1, err
}

func IterateExecute(ctx context.Context, nextTime int64, env string, fn func(*Execute) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("next_time <= ?", nextTime).
		And("env = ?", env).
		And("is_enabled = 1").
		Iterate(new(Execute), func(idx int, obj interface{}) error {
			return fn(obj.(*Execute))
		})
}

func ListTimer(ctx context.Context, reqDTO ListTimerReqDTO) ([]Timer, int64, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId).
		And("env = ?", reqDTO.Env)
	if reqDTO.Name != "" {
		session.And("name like ?", reqDTO.Name+"%")
	}
	ret := make([]Timer, 0)
	total, err := session.
		Desc("id").
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}

func InsertLog(ctx context.Context, reqDTO InsertLogReqDTO) error {
	ret := Log{
		TimerId: reqDTO.TimerId,
		TaskContent: &xormutil.Conversion[timer.Task]{
			Data: reqDTO.TaskContent,
		},
		ErrLog:      reqDTO.ErrLog,
		TriggerType: reqDTO.TriggerType,
		TriggerBy:   reqDTO.TriggerBy,
		IsSuccess:   reqDTO.IsSuccess,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func ListLog(ctx context.Context, reqDTO ListLogReqDTO) ([]Log, int64, error) {
	ret := make([]Log, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("timer_id = ?", reqDTO.TimerId).
		And("created between ? and ?", reqDTO.BeginTime.Format(time.DateTime), reqDTO.EndTime.Format(time.DateTime)).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		Desc("id").
		FindAndCount(&ret)
	return ret, total, err
}

func GetTimerById(ctx context.Context, id int64) (Timer, bool, error) {
	var ret Timer
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func DeleteTimer(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Timer))
	return err
}

func InsertExecute(ctx context.Context, reqDTO InsertExecuteReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Execute{
			TimerId:   reqDTO.TimerId,
			IsEnabled: reqDTO.IsEnabled,
			NextTime:  reqDTO.NextTime,
			Env:       reqDTO.Env,
		})
	return err
}

func DeleteExecuteByTimerId(ctx context.Context, timerId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("timer_id = ?", timerId).
		Delete(new(Execute))
	return rows == 1, err
}

func CountTimerByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Count(new(Timer))
}

func DeleteLogByTimerId(ctx context.Context, timerId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("timer_id = ?", timerId).
		Delete(new(Log))
	return err
}
