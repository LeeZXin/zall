package cmd

import (
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var Actions = &cli.Command{
	Name:        "actions",
	Usage:       "This command starts actions server",
	Description: "zgit actions provides cicd actions",
	Action:      runActions,
}

func runActions(*cli.Context) error {
	// action
	actionagentapi.InitApi()
	actionsrv.InitSrv()
	zsf.Run(
		zsf.WithLifeCycles(
			httpserver.NewServer(),
		),
	)
	return nil
}
