package taskmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsTaskNameValid(name string) bool {
	return len(name) <= 64 && len(name) > 0
}

func InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) (Task, error) {
	ret := Task{
		Name:      reqDTO.Name,
		CronExp:   reqDTO.CronExp,
		Content:   &reqDTO.Content,
		TeamId:    reqDTO.TeamId,
		Env:       reqDTO.Env,
		IsEnabled: reqDTO.IsEnabled,
		Creator:   reqDTO.Creator,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func UpdateTask(ctx context.Context, reqDTO UpdateTaskReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.TaskId).
		Cols("name", "cron_exp", "content").
		Update(&Task{
			Name:    reqDTO.Name,
			CronExp: reqDTO.CronExp,
			Content: &reqDTO.Content,
		})
	return rows == 1, err
}

func UpdateExecuteNextTime(ctx context.Context, taskId, nextTime, runVersion int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		And("run_version = ?", runVersion).
		Cols("next_time", "run_version").
		Update(&Execute{
			NextTime:   nextTime,
			RunVersion: runVersion + 1,
		})
	return rows == 1, err
}

func EnableExecute(ctx context.Context, taskId int64, nextTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		And("is_enabled = 0").
		Cols("is_enabled", "next_time").
		Update(&Execute{
			IsEnabled: true,
			NextTime:  nextTime,
		})
	return rows == 1, err
}

func EnableTask(ctx context.Context, taskId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", taskId).
		And("is_enabled = 0").
		Cols("is_enabled").
		Update(&Task{
			IsEnabled: true,
		})
	return rows == 1, err
}

func DisableExecute(ctx context.Context, taskId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		And("is_enabled = 1").
		Cols("is_enabled", "next_time").
		Update(&Execute{
			IsEnabled: false,
			NextTime:  0,
		})
	return rows == 1, err
}

func DisableTask(ctx context.Context, taskId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", taskId).
		And("is_enabled = 1").
		Cols("is_enabled").
		Update(&Task{
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
			if err := fn(obj.(*Execute)); err != nil {
				return err
			}
			return nil
		})
}

func PageTask(ctx context.Context, reqDTO PageTaskReqDTO) ([]Task, int64, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId).
		And("env = ?", reqDTO.Env)
	if reqDTO.Name != "" {
		session.And("name like ?", reqDTO.Name+"%")
	}
	ret := make([]Task, 0)
	total, err := session.
		Desc("id").
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}

func InsertTaskLog(ctx context.Context, reqDTO InsertTaskLogReqDTO) error {
	ret := TaskLog{
		TaskId:      reqDTO.TaskId,
		TaskContent: reqDTO.TaskContent,
		ErrLog:      reqDTO.ErrLog,
		TriggerType: reqDTO.TriggerType,
		TriggerBy:   reqDTO.TriggerBy,
		IsSuccess:   reqDTO.IsSuccess,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func ListTaskLog(ctx context.Context, reqDTO ListTaskLogReqDTO) ([]TaskLog, int64, error) {
	ret := make([]TaskLog, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", reqDTO.TaskId).
		And("created between ? and ?", reqDTO.BeginTime.Format(time.DateTime), reqDTO.EndTime.Format(time.DateTime)).
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		OrderBy("id desc").
		FindAndCount(&ret)
	return ret, total, err
}

func GetTaskById(ctx context.Context, id int64) (Task, bool, error) {
	var ret Task
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func DeleteTask(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(Task))
	return err
}

func InsertExecute(ctx context.Context, reqDTO InsertExecuteReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Execute{
			TaskId:    reqDTO.TaskId,
			IsEnabled: reqDTO.IsEnabled,
			NextTime:  reqDTO.NextTime,
			Env:       reqDTO.Env,
		})
	return err
}

func DeleteExecuteByTaskId(ctx context.Context, taskId int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		Delete(new(Execute))
	return rows == 1, err
}

func CountTaskByTeamId(ctx context.Context, teamId int64) (int64, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", teamId).
		Count(new(Task))
}

func DeleteLogByTaskId(ctx context.Context, taskId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", taskId).
		Delete(new(TaskLog))
	return err
}
