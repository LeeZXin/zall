package alertsrv

import (
	"errors"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/patrickmn/go-cache"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	SnapshotPath = "alert-limiter.snapshot"
)

type Limiter struct {
	store        *cache.Cache
	lock         sync.RWMutex
	snapshotTask *taskutil.PeriodicalTask
}

func NewLimiter() *Limiter {
	ret := new(Limiter)
	ret.store = util.NewGoCache()
	err := ret.store.LoadFile(SnapshotPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Logger.Errorf("load limiter snapshot: %v failed with err: %v", SnapshotPath, err)
	}
	ret.snapshotTask, _ = taskutil.NewPeriodicalTask(30*time.Second, func() {
		err := ret.store.SaveFile(SnapshotPath)
		if err != nil {
			logger.Logger.Errorf("save limiter snapshot: %v failed with err: %v", SnapshotPath, err)
		}
	})
	ret.snapshotTask.Start()
	quit.AddShutdownHook(ret.snapshotTask.Stop, true)
	return ret
}

func (l *Limiter) TryPass(id int64, expired time.Duration) bool {
	key := strconv.FormatInt(id, 10)
	l.lock.RLock()
	_, b := l.store.Get(key)
	l.lock.RUnlock()
	if b {
		return false
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	_, b = l.store.Get(key)
	if b {
		return false
	}
	l.store.Set(key, true, expired)
	return true
}
