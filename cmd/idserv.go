package cmd

import (
	"github.com/LeeZXin/zall/genid/modules/api/idapi"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/services/registry"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var IdServer = &cli.Command{
	Name:        "idServer",
	Usage:       "This command starts id server",
	Description: "this provides id-generating service",
	Action:      runIdServer,
}

func runIdServer(*cli.Context) error {
	{
		idapi.InitApi()
	}
	zsf.Run(
		zsf.WithLifeCycles(
			httpserver.NewServer(
				httpserver.WithRegistryAction(
					registry.NewDefaultHttpAction(registry.NewDefaultEtcdRegistry()),
				),
			),
		),
	)
	return nil
}
