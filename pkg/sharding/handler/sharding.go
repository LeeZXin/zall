package handler

import (
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"time"
)

type HeartbeatHandler func()
type DeleteInstanceHandler func()
type TaskHandler func()

type ShardingPeriodicalHandler struct {
	cfg           *Config
	heartbeatTask *taskutil.PeriodicalTask
	executeTask   *taskutil.PeriodicalTask
}

type Config struct {
	HeartbeatInterval     time.Duration
	HeartbeatHandler      HeartbeatHandler
	DeleteInstanceHandler DeleteInstanceHandler

	TaskInterval time.Duration
	TaskHandler  TaskHandler
}

func NewShardingPeriodicalHandler(cfg *Config) (*ShardingPeriodicalHandler, error) {
	var err error
	ret := new(ShardingPeriodicalHandler)
	ret.cfg = cfg
	ret.heartbeatTask, err = taskutil.NewPeriodicalTask(cfg.HeartbeatInterval, cfg.HeartbeatHandler)
	if err != nil {
		return nil, err
	}
	ret.executeTask, err = taskutil.NewPeriodicalTask(cfg.TaskInterval, ret.cfg.TaskHandler)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (h *ShardingPeriodicalHandler) Start() {
	h.cfg.HeartbeatHandler()
	h.heartbeatTask.Start()
	h.cfg.TaskHandler()
	h.executeTask.Start()
	quit.AddShutdownHook(h.Stop, true)
}

func (h *ShardingPeriodicalHandler) Stop() {
	if h.heartbeatTask != nil {
		h.heartbeatTask.Stop()
	}
	if h.executeTask != nil {
		h.executeTask.Stop()
	}
	if h.cfg.DeleteInstanceHandler != nil {
		h.cfg.DeleteInstanceHandler()
	}
}
