package crontask

import (
	"github.com/LeeZXin/zsf/logger"
	"github.com/robfig/cron/v3"
	"strings"
	"sync"
	"time"
)

var (
	c        *cron.Cron
	initOnce = sync.Once{}
)

type cronLogger struct {
}

// Info logs routine messages about cron's operation.
func (*cronLogger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	logger.Logger.Infof(
		formatString(len(keysAndValues)),
		append([]interface{}{msg}, keysAndValues...)...,
	)
}

// Error logs an error condition.
func (*cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	logger.Logger.Errorf(formatString(len(keysAndValues)+2),
		append([]interface{}{msg, "error", err}, keysAndValues...)...)
}

func initCron() {
	initOnce.Do(func() {
		c = cron.New(
			cron.WithLogger(new(cronLogger)),
			cron.WithParser(cron.NewParser(cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		)
		c.Start()
	})
}

func AddFunc(spec string, fn func()) error {
	initCron()
	_, err := c.AddFunc(spec, fn)
	return err
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []any) []interface{} {
	var formattedArgs []any
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
