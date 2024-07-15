package tasksrv

import (
	"context"
	"github.com/robfig/cron/v3"
)

var (
	Outer OuterService
)

func Init() {
	if Outer == nil {
		Outer = new(outerImpl)
		initTask()
		parser = cron.NewParser(
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
		)
	}
}

type OuterService interface {
	// CreateTask 新增任务
	CreateTask(context.Context, CreateTaskReqDTO) error
	// ListTask 展示任务列表
	ListTask(context.Context, ListTaskReqDTO) ([]TaskDTO, int64, error)
	// EnableTask 启动任务
	EnableTask(context.Context, EnableTaskReqDTO) error
	// DisableTask 关闭任务
	DisableTask(context.Context, DisableTaskReqDTO) error
	// DeleteTask 删除任务
	DeleteTask(context.Context, DeleteTaskReqDTO) error
	// PageTaskLog 获取执行历史
	PageTaskLog(context.Context, PageTaskLogReqDTO) ([]TaskLogDTO, int64, error)
	// TriggerTask 手动执行任务
	TriggerTask(context.Context, TriggerTaskReqDTO) error
	// UpdateTask 更新任务
	UpdateTask(context.Context, UpdateTaskReqDTO) error
}
