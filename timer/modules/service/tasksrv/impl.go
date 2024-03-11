package tasksrv

import (
	"context"
	"encoding/json"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/timer/modules/model/taskmd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"time"
)

type outerImpl struct{}

// InsertTask 新增任务
func (o *outerImpl) InsertTask(ctx context.Context, reqDTO InsertTaskReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TimerTaskSrvKeysVO.InsertTask),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err = checkPerm(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
		return
	}
	obj := TaskObj{
		TaskType: reqDTO.TaskType,
	}
	switch reqDTO.TaskType {
	case HttpTaskType:
		m, _ := json.Marshal(reqDTO.HttpTask)
		obj.Content = string(m)
	}
	objJson, _ := json.Marshal(obj)
	if err = taskmd.InsertTask(ctx, taskmd.InsertTaskReqDTO{
		Name:       reqDTO.Name,
		CronExp:    reqDTO.CronExp,
		Content:    string(objJson),
		NextTime:   0,
		TaskStatus: taskmd.Closed,
		TeamId:     reqDTO.TeamId,
	}); err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

// ListTask 展示任务列表
func (o *outerImpl) ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]TaskDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkPerm(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
		return nil, 0, err
	}
	tasks, err := taskmd.ListTask(ctx, taskmd.ListTaskReqDTO{
		TeamId: reqDTO.TeamId,
		Name:   reqDTO.Name,
		Cursor: reqDTO.Cursor,
		Limit:  reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(tasks, func(t taskmd.Task) (TaskDTO, error) {
		var obj TaskObj
		_ = json.Unmarshal([]byte(t.Content), &obj)
		dto := TaskDTO{
			Id:         t.Id,
			Name:       t.Name,
			CronExp:    t.CronExp,
			TaskType:   obj.TaskType,
			TeamId:     t.TeamId,
			NextTime:   t.NextTime,
			TaskStatus: t.TaskStatus,
		}
		switch obj.TaskType {
		case HttpTaskType:
			var h HttpTask
			_ = json.Unmarshal([]byte(obj.Content), &h)
			dto.HttpTask = h
		}
		return dto, nil
	})
	if len(tasks) == reqDTO.Limit {
		return ret, tasks[len(tasks)-1].Id, nil
	}
	return ret, 0, nil
}

// EnableTask 启动任务
func (o *outerImpl) EnableTask(ctx context.Context, reqDTO EnableTaskReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TimerTaskSrvKeysVO.EnableTask),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		task taskmd.Task
		b    bool
	)
	task, err = checkPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return
	}
	cron, err := parseCron(task.CronExp)
	if err != nil {
		return util.ThereHasBugErr()
	}
	next := cron.Next(time.Now())
	for i := 0; i < 10; i++ {
		b, err = taskmd.UpdateTaskNextTimeAndStatus(ctx, reqDTO.Id, taskmd.Pending, next.UnixMilli(), task.Version)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if b {
			return
		}
		task, b, err = taskmd.GetTaskById(ctx, reqDTO.Id)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if !b {
			err = util.OperationFailedError()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	err = util.OperationFailedError()
	return
}

// DisableTask 关闭任务
func (o *outerImpl) DisableTask(ctx context.Context, reqDTO DisableTaskReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TimerTaskSrvKeysVO.DisableTask),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		task taskmd.Task
		b    bool
	)
	task, err = checkPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		b, err = taskmd.UpdateTaskStatus(ctx, reqDTO.Id, taskmd.Closed, task.Version)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if b {
			return
		}
		task, b, err = taskmd.GetTaskById(ctx, reqDTO.Id)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if !b {
			err = util.OperationFailedError()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	err = util.OperationFailedError()
	return
}

// DeleteTask 删除任务
func (o *outerImpl) DeleteTask(ctx context.Context, reqDTO DeleteTaskReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TimerTaskSrvKeysVO.DeleteTask),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err = checkPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return
	}
	err = taskmd.DeleteTask(ctx, reqDTO.Id)
	if err != nil {
		err = util.InternalError(err)
		return
	}
	return
}

// TriggerTask 手动执行任务
func (o *outerImpl) TriggerTask(ctx context.Context, reqDTO TriggerTaskReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TimerTaskSrvKeysVO.TriggerTask),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		task taskmd.Task
	)
	task, err = checkPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return
	}
	triggerTask(&task, reqDTO.Operator.Account)
	return nil
}

// ListTaskLog 获取执行历史
func (o *outerImpl) ListTaskLog(ctx context.Context, reqDTO ListTaskLogReqDTO) ([]TaskLogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return nil, 0, err
	}
	logs, err := taskmd.ListTaskLog(ctx, taskmd.ListTaskLogReqDTO{
		Id:     reqDTO.Id,
		Cursor: reqDTO.Cursor,
		Limit:  reqDTO.Limit,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(logs, func(t taskmd.TaskLog) (TaskLogDTO, error) {
		var obj TaskObj
		_ = json.Unmarshal([]byte(t.TaskContent), &obj)
		dto := TaskLogDTO{
			TaskType:    obj.TaskType,
			LogContent:  t.LogContent,
			TriggerType: t.TriggerType,
			TriggerBy:   t.TriggerBy,
			TaskStatus:  t.TaskStatus,
			Created:     t.Created,
		}
		switch obj.TaskType {
		case HttpTaskType:
			var h HttpTask
			_ = json.Unmarshal([]byte(obj.Content), &h)
			dto.HttpTask = h
		}
		return dto, nil
	})
	if len(logs) == reqDTO.Limit {
		return ret, logs[len(logs)-1].Id, nil
	}
	return ret, 0, nil
}

// UpdateTask 更新任务
func (o *outerImpl) UpdateTask(ctx context.Context, reqDTO UpdateTaskReqDTO) (err error) {
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.TimerTaskSrvKeysVO.UpdateTask),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		task taskmd.Task
		b    bool
	)
	task, err = checkPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return
	}
	obj := TaskObj{
		TaskType: reqDTO.TaskType,
	}
	switch reqDTO.TaskType {
	case HttpTaskType:
		m, _ := json.Marshal(reqDTO.HttpTask)
		obj.Content = string(m)
	}
	objJson, _ := json.Marshal(obj)
	for i := 0; i < 10; i++ {
		b, err = taskmd.UpdateTask(ctx, taskmd.UpdateTaskReqDTO{
			Id:         reqDTO.Id,
			Name:       reqDTO.Name,
			CronExp:    reqDTO.CronExp,
			Content:    string(objJson),
			TaskStatus: taskmd.Closed,
			NextTime:   0,
			Version:    task.Version,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if b {
			return
		}
		task, b, err = taskmd.GetTaskById(ctx, reqDTO.Id)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			err = util.InternalError(err)
			return
		}
		if !b {
			err = util.OperationFailedError()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	err = util.OperationFailedError()
	return
}

func checkPerm(ctx context.Context, operator apisession.UserInfo, teamId int64) error {
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.TeamPerm.CanHandleTimer {
		return util.UnauthorizedError()
	}
	return nil
}

func checkPermByTaskId(ctx context.Context, operator apisession.UserInfo, taskId int64) (taskmd.Task, error) {
	task, b, err := taskmd.GetTaskById(ctx, taskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return taskmd.Task{}, util.InternalError(err)
	}
	if !b {
		return taskmd.Task{}, util.InvalidArgsError()
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, task.TeamId, operator.Account)
	if !b {
		return task, util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.TeamPerm.CanHandleTimer {
		return task, util.UnauthorizedError()
	}
	return task, nil
}
