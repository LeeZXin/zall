package tasksrv

import "github.com/robfig/cron/v3"

var (
	parser cron.Parser
)

func init() {
	parser = cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)
}

func ParseCron(spec string) (cron.Schedule, error) {
	return parser.Parse(spec)
}
