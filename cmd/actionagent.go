package cmd

import (
	zssh "github.com/LeeZXin/zall/pkg/ssh"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var WorkflowAgentServer = &cli.Command{
	Name:        "sshAgent",
	Usage:       "This command starts ssh server",
	Description: "actionAgent provides ssh agent to execute command",
	Action:      runActionAgent,
}

func runActionAgent(*cli.Context) error {
	zsf.Run(
		zsf.WithLifeCycles(
			zssh.NewAgentServer(),
		),
	)
	return nil
}
