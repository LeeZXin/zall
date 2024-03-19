package cmd

import (
	"github.com/LeeZXin/zall/git/modules/api/actionagentapi"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var ActionAgent = &cli.Command{
	Name:        "actionAgent",
	Usage:       "This command starts actions server",
	Description: "zgit actions provides action agent",
	Action:      runActionAgent,
}

func runActionAgent(*cli.Context) error {
	// action
	actionagentapi.InitApi()
	zsf.Run(
		zsf.WithLifeCycles(
			httpserver.NewServer(),
		),
	)
	return nil
}
