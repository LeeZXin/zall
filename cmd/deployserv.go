package cmd

import (
	"github.com/LeeZXin/zall/services/modules/api/deployapi"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var DeployServer = &cli.Command{
	Name:        "deployServer",
	Usage:       "This command starts deploy server",
	Description: "this provides deploy service",
	Action:      runDeployServer,
}

func runDeployServer(*cli.Context) error {
	deployapi.InitDeploy()
	zsf.Run(
		zsf.WithLifeCycles(
			httpserver.NewServer(),
		),
	)
	return nil
}
