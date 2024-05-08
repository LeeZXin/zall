package limiter

import "sync/atomic"

type Limiter interface {
	Borrow() bool
	Return()
}

// CountLimiter 简单的次数限制
type CountLimiter struct {
	limit     *atomic.Int64
	initLimit int64
}

func (l *CountLimiter) Borrow() bool {
	ret := l.limit.Add(-1)
	if ret < 0 {
		l.limit.Add(1)
		return false
	}
	return true
}

func (l *CountLimiter) Return() {
	ret := l.limit.Add(1)
	if ret > l.initLimit {
		l.limit.Add(-1)
	}
}

func NewCountLimiter(limit int64) Limiter {
	if limit < 0 {
		limit = 0
	}
	l := &atomic.Int64{}
	l.Store(limit)
	return &CountLimiter{
		limit:     l,
		initLimit: limit,
	}
}
