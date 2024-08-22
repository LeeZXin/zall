package cmd

import (
	"github.com/LeeZXin/zall/timer/modules/api/timerapi"
	"github.com/LeeZXin/zsf/http/httpserver"
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
	timerapi.InitApi()
	zsf.Run(
		zsf.WithLifeCycles(httpserver.NewServer()),
	)
	return nil
}
