package tasksrv

import "github.com/robfig/cron/v3"

var (
	parser cron.Parser
)

func ParseCron(spec string) (cron.Schedule, error) {
	return parser.Parse(spec)
}
