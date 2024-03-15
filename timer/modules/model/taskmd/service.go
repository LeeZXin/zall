package taskmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
)

func IsTaskNameValid(name string) bool {
	return len(name) <= 64 && len(name) > 0
}

func InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) error {
	ret := Task{
		Name:       reqDTO.Name,
		CronExp:    reqDTO.CronExp,
		Content:    reqDTO.Content,
		NextTime:   reqDTO.NextTime,
		TaskStatus: reqDTO.TaskStatus,
		TeamId:     reqDTO.TeamId,
		Version:    0,
	}
	_, err := xormutil.MustGetXormSession(ctx).Table("ztimer_task_content_" + reqDTO.Env).Insert(&ret)
	return err
}

func UpdateTask(ctx context.Context, reqDTO UpdateTaskReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		And("version = ?", reqDTO.Version).
		Cols("name", "cron_exp", "content", "next_time", "task_status", "version").
		Table("ztimer_task_content_" + reqDTO.Env).
		Update(&Task{
			Name:       reqDTO.Name,
			CronExp:    reqDTO.CronExp,
			Content:    reqDTO.Content,
			NextTime:   reqDTO.NextTime,
			TaskStatus: reqDTO.TaskStatus,
			Version:    reqDTO.Version + 1,
		})
	return rows == 1, err
}

func UpdateTaskStatus(ctx context.Context, reqDTO UpdateTaskStatusReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.TaskId).
		And("version = ?", reqDTO.Version).
		Cols("task_status", "version").
		Table("ztimer_task_content_" + reqDTO.Env).
		Update(&Task{
			TaskStatus: reqDTO.NewStatus,
			Version:    reqDTO.Version + 1,
		})
	return rows == 1, err
}

func UpdateTaskNextTimeAndStatus(ctx context.Context, reqDTO UpdateTaskNextTimeAndStatusReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.TaskId).
		And("version = ?", reqDTO.Version).
		Cols("next_time", "version", "task_status").
		Table("ztimer_task_content_" + reqDTO.Env).
		Update(&Task{
			TaskStatus: reqDTO.Status,
			NextTime:   reqDTO.NextTime,
			Version:    reqDTO.Version + 1,
		})
	return rows == 1, err
}

func IterateTask(ctx context.Context, nextTime int64, status TaskStatus, env string, fn func(*Task) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("next_time <= ?", nextTime).
		And("task_status = ?", status).
		Table("ztimer_task_content_"+env).
		Iterate(new(Task), func(idx int, obj interface{}) error {
			if err := fn(obj.(*Task)); err != nil {
				return err
			}
			return nil
		})
}

func ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]Task, error) {
	session := xormutil.MustGetXormSession(ctx).
		Where("team_id = ?", reqDTO.TeamId).
		Table("ztimer_task_content_" + reqDTO.Env)
	if reqDTO.Cursor > 0 {
		session.And("id > ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	if reqDTO.Name != "" {
		session.And("name like ?", "%"+reqDTO.Name+"%")
	}
	ret := make([]Task, 0)
	err := session.OrderBy("id asc").Find(&ret)
	return ret, err
}

func InsertTaskLog(ctx context.Context, reqDTO InsertTaskLogReqDTO) error {
	ret := TaskLog{
		TaskId:      reqDTO.TaskId,
		TaskContent: reqDTO.TaskContent,
		LogContent:  reqDTO.LogContent,
		TriggerType: reqDTO.TriggerType,
		TriggerBy:   reqDTO.TriggerBy,
		TaskStatus:  reqDTO.TaskStatus,
	}
	_, err := xormutil.MustGetXormSession(ctx).Table("ztimer_task_log_" + reqDTO.Env).Insert(&ret)
	return err
}

func ListTaskLog(ctx context.Context, reqDTO ListTaskLogReqDTO) ([]TaskLog, error) {
	ret := make([]TaskLog, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("task_id = ?", reqDTO.Id).
		Table("ztimer_task_log_" + reqDTO.Env)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func InsertInstance(ctx context.Context, instanceId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("ztimer_instance_" + env).
		Insert(&Instance{
			InstanceId:    instanceId,
			HeartbeatTime: time.Now().UnixMilli(),
		})
	return err
}

func UpdateHeartbeatTime(ctx context.Context, instanceId string, heartbeatTime int64, env string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Cols("heartbeat_time").
		Table("ztimer_instance_" + env).
		Update(&Instance{
			HeartbeatTime: heartbeatTime,
		})
	return rows == 1, err
}

func DeleteInstance(ctx context.Context, instanceId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Table("ztimer_instance_" + env).
		Limit(1).
		Delete(new(Instance))
	return err
}

func GetValidInstances(ctx context.Context, after int64, env string) ([]Instance, error) {
	ret := make([]Instance, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("heartbeat_time>= ?", after).
		Table("ztimer_instance_" + env).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func DeleteInValidInstances(ctx context.Context, before int64, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("heartbeat_time < ?", before).
		Table("ztimer_instance_" + env).
		Delete(new(Instance))
	return err
}

func GetTaskById(ctx context.Context, id int64, env string) (Task, bool, error) {
	var ret Task
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Table("ztimer_task_content_" + env).
		Get(&ret)
	return ret, b, err
}

func DeleteTask(ctx context.Context, id int64, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Table("ztimer_task_content_" + env).
		Delete(new(Task))
	return err
}
