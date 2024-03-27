package cmd

import (
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var ActionAgent = &cli.Command{
	Name:        "actionAgent",
	Usage:       "This command starts actions server",
	Description: "actionAgent provides ssh action agent to execute command",
	Action:      runActionAgent,
}

func runActionAgent(*cli.Context) error {
	zsf.Run(
		zsf.WithLifeCycles(
			action.NewAgentServer(),
		),
	)
	return nil
}
