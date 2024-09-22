package crontask

import (
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	c        *cron.Cron
	initOnce = sync.Once{}
)

func initCron() {
	initOnce.Do(func() {
		c = cron.New(
			cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)),
		)
		c.Start()
	})
}

func AddFunc(spec string, fn func()) error {
	initCron()
	_, err := c.AddFunc(spec, fn)
	return err
}
