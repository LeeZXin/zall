package cmd

import (
	"github.com/LeeZXin/zall/timer/modules/service/tasksrv"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var Timer = &cli.Command{
	Name:        "timer",
	Usage:       "This command starts timer server",
	Description: "this provides timer service",
	Action:      runTimer,
}

func runTimer(*cli.Context) error {
	tasksrv.InitTask(static.GetString("timer.env"))
	zsf.Run()
	return nil
}
