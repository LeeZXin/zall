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
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func UpdateTask(ctx context.Context, reqDTO UpdateTaskReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		And("version = ?", reqDTO.Version).
		Cols("name", "cron_exp", "content", "next_time", "task_status", "version").
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

func UpdateTaskStatus(ctx context.Context, taskId int64, newStatus TaskStatus, version int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", taskId).
		And("version = ?", version).
		Cols("task_status", "version").
		Update(&Task{
			TaskStatus: newStatus,
			Version:    version + 1,
		})
	return rows == 1, err
}

func UpdateTaskNextTimeAndStatus(ctx context.Context, taskId int64, status TaskStatus, nextTime, version int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", taskId).
		And("version = ?", version).
		Cols("next_time", "version", "task_status").
		Update(&Task{
			TaskStatus: status,
			NextTime:   nextTime,
			Version:    version + 1,
		})
	return rows == 1, err
}

func IterateTask(ctx context.Context, nextTime int64, status TaskStatus, fn func(*Task) error) error {
	return xormutil.MustGetXormSession(ctx).
		Where("next_time <= ?", nextTime).
		And("task_status = ?", status).
		Iterate(new(Task), func(idx int, obj interface{}) error {
			if err := fn(obj.(*Task)); err != nil {
				return err
			}
			return nil
		})
}

func ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]Task, error) {
	session := xormutil.MustGetXormSession(ctx).Where("team_id = ?", reqDTO.TeamId)
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
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return err
}

func ListTaskLog(ctx context.Context, reqDTO ListTaskLogReqDTO) ([]TaskLog, error) {
	ret := make([]TaskLog, 0)
	session := xormutil.MustGetXormSession(ctx).Where("task_id = ?", reqDTO.Id)
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func UpdateLogContentAndStatus(ctx context.Context, reqDTO UpdateLogContentAndStatusReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("log_content", "task_status").
		Update(&TaskLog{
			LogContent: reqDTO.Content,
			TaskStatus: reqDTO.Status,
		})
	return err
}

func InsertInstance(ctx context.Context, instanceId string) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Instance{
		InstanceId:    instanceId,
		HeartbeatTime: time.Now().UnixMilli(),
	})
	return err
}

func UpdateHeartbeatTime(ctx context.Context, instanceId string, heartbeatTime int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Cols("heartbeat_time").
		Update(&Instance{
			HeartbeatTime: heartbeatTime,
		})
	return rows == 1, err
}

func DeleteInstance(ctx context.Context, instanceId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("instance_id = ?", instanceId).
		Limit(1).
		Delete(new(Instance))
	return err
}

func GetValidInstances(ctx context.Context, after int64) ([]Instance, error) {
	ret := make([]Instance, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("heartbeat_time>= ?", after).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func DeleteInValidInstances(ctx context.Context, before int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("heartbeat_time < ?", before).
		Delete(new(Instance))
	return err
}

func GetTaskById(ctx context.Context, id int64) (Task, bool, error) {
	var ret Task
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func DeleteTask(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Task))
	return err
}
