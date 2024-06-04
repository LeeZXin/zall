package alertsrv

import (
	"context"
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
	store                *cache.Cache
	lock                 sync.RWMutex
	snapshotTaskStopFunc taskutil.StopFunc
}

func NewLimiter() *Limiter {
	ret := new(Limiter)
	ret.store = util.NewGoCache()
	err := ret.store.LoadFile(SnapshotPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Logger.Errorf("load limiter snapshot: %v failed with err: %v", SnapshotPath, err)
	}
	ret.snapshotTaskStopFunc, _ = taskutil.RunPeriodicalTask(30*time.Second, 30*time.Second, func(context.Context) {
		err := ret.store.SaveFile(SnapshotPath)
		if err != nil {
			logger.Logger.Errorf("save limiter snapshot: %v failed with err: %v", SnapshotPath, err)
		}
	})
	quit.AddShutdownHook(quit.ShutdownHook(ret.snapshotTaskStopFunc), true)
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
