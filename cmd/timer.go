package cmd

import (
	"github.com/LeeZXin/zall/timer/modules/service/tasksrv"
	"github.com/LeeZXin/zsf/starter"
	"github.com/urfave/cli/v2"
)

var Timer = &cli.Command{
	Name:        "timer",
	Usage:       "This command starts timer server",
	Description: "this provides timer service",
	Action:      runTimer,
}

func runTimer(*cli.Context) error {
	{
		tasksrv.InitTask()
	}
	starter.Run()
	return nil
}
