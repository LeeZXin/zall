package cmd

import (
	"github.com/LeeZXin/zall/prop/modules/api/propapi"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var PropServer = &cli.Command{
	Name:        "propServer",
	Usage:       "This command starts prop server",
	Description: "this provides prop service",
	Action:      runPropServer,
}

func runPropServer(*cli.Context) error {
	propapi.InitApi()
	zsf.Run(
		zsf.WithLifeCycles(
			httpserver.NewServer(),
		),
	)
	return nil
}
