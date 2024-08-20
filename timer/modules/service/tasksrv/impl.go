package tasksrv

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/pkg/apicode"
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

// CreateTask 新增任务
func (o *outerImpl) CreateTask(ctx context.Context, reqDTO CreateTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkManageTimerPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
		return err
	}
	err := xormstore.WithTx(ctx, func(ctx context.Context) error {
		task, err2 := taskmd.InsertTask(ctx, taskmd.InsertTaskReqDTO{
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
		return taskmd.InsertExecute(ctx, taskmd.InsertExecuteReqDTO{
			TaskId:    task.Id,
			IsEnabled: false,
			Env:       reqDTO.Env,
		})
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// ListTask 展示任务列表
func (o *outerImpl) ListTask(ctx context.Context, reqDTO ListTaskReqDTO) ([]TaskDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	if err := checkManageTimerPermByTeamId(ctx, reqDTO.Operator, reqDTO.TeamId); err != nil {
		return nil, 0, err
	}
	tasks, total, err := taskmd.PageTask(ctx, taskmd.PageTaskReqDTO{
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
	ret, _ := listutil.Map(tasks, func(t taskmd.Task) (TaskDTO, error) {
		return TaskDTO{
			Id:        t.Id,
			Name:      t.Name,
			CronExp:   t.CronExp,
			Task:      t.GetContent(),
			TeamId:    t.TeamId,
			IsEnabled: t.IsEnabled,
			Env:       t.Env,
			Creator:   t.Creator,
		}, nil
	})
	return ret, total, nil
}

// EnableTask 启动任务
func (o *outerImpl) EnableTask(ctx context.Context, reqDTO EnableTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	task, err := checkManageTimerPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	cron, err := ParseCron(task.CronExp)
	if err != nil {
		return util.ThereHasBugErr()
	}
	now := time.Now()
	nextTime := cron.Next(now)
	if nextTime.Before(now) {
		return util.NewBizErr(apicode.OperationFailedErrCode, i18n.CronExpError)
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := taskmd.EnableExecute(ctx,
			reqDTO.Id,
			nextTime.UnixMilli(),
		)
		if err2 != nil {
			return err2
		}
		_, err2 = taskmd.EnableTask(ctx, reqDTO.Id)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DisableTask 关闭任务
func (o *outerImpl) DisableTask(ctx context.Context, reqDTO DisableTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageTimerPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		_, err2 := taskmd.DisableExecute(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		_, err2 = taskmd.DisableTask(ctx, reqDTO.Id)
		return err2
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteTask 删除任务
func (o *outerImpl) DeleteTask(ctx context.Context, reqDTO DeleteTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageTimerPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err2 := taskmd.DeleteTask(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		_, err2 = taskmd.DeleteExecuteByTaskId(ctx, reqDTO.Id)
		if err2 != nil {
			return err2
		}
		err2 = taskmd.DeleteLogByTaskId(ctx, reqDTO.Id)
		return err2
	})
	if err != nil {
		return util.InternalError(err)
	}
	return nil
}

// TriggerTask 手动执行任务
func (o *outerImpl) TriggerTask(ctx context.Context, reqDTO TriggerTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		task taskmd.Task
	)
	task, err := checkManageTimerPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	err = triggerTask(&task, reqDTO.Operator.Account)
	if err != nil {
		return util.OperationFailedError()
	}
	return nil
}

// ListTaskLog 获取执行历史
func (o *outerImpl) ListTaskLog(ctx context.Context, reqDTO ListTaskLogReqDTO) ([]TaskLogDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, err := checkManageTimerPermByTaskId(ctx, reqDTO.Operator, reqDTO.TaskId)
	if err != nil {
		return nil, 0, err
	}
	// 当月第一天
	beginTime := reqDTO.monthTime.AddDate(0, 0, -reqDTO.monthTime.Day()+1)
	// 当月最后一天
	endTime := reqDTO.monthTime.AddDate(0, 1, -1)
	logs, total, err := taskmd.ListTaskLog(ctx, taskmd.ListTaskLogReqDTO{
		TaskId:    reqDTO.TaskId,
		PageNum:   reqDTO.PageNum,
		PageSize:  10,
		BeginTime: beginTime,
		EndTime:   time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, endTime.Location()),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	ret, _ := listutil.Map(logs, func(t taskmd.TaskLog) (TaskLogDTO, error) {
		return TaskLogDTO{
			Task:        t.GetTaskContent(),
			ErrLog:      t.ErrLog,
			TriggerType: t.TriggerType,
			TriggerBy:   t.TriggerBy,
			IsSuccess:   t.IsSuccess,
			Created:     t.Created,
		}, nil
	})
	return ret, total, nil
}

// UpdateTask 更新任务
func (o *outerImpl) UpdateTask(ctx context.Context, reqDTO UpdateTaskReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkManageTimerPermByTaskId(ctx, reqDTO.Operator, reqDTO.Id)
	if err != nil {
		return err
	}
	_, err = taskmd.UpdateTask(ctx, taskmd.UpdateTaskReqDTO{
		Id:      reqDTO.Id,
		Name:    reqDTO.Name,
		CronExp: reqDTO.CronExp,
		Content: reqDTO.Task,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func checkManageTimerPermByTeamId(ctx context.Context, operator apisession.UserInfo, teamId int64) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if !p.IsAdmin && !p.PermDetail.TeamPerm.CanManageTimer {
		return util.UnauthorizedError()
	}
	return nil
}

func checkManageTimerPermByTaskId(ctx context.Context, operator apisession.UserInfo, taskId int64) (taskmd.Task, error) {
	task, b, err := taskmd.GetTaskById(ctx, taskId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return taskmd.Task{}, util.InternalError(err)
	}
	if !b {
		return taskmd.Task{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return task, nil
	}
	return task, checkManageTimerPermByTeamId(ctx, operator, task.TeamId)
}
