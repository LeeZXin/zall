package cmd

import (
	"github.com/LeeZXin/zall/pkg/sshagent"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var WorkflowAgentServer = &cli.Command{
	Name:        "workflowAgent",
	Usage:       "This command starts workflow agent server",
	Description: "actionAgent provides ssh agent to execute command or workflow",
	Action:      runActionAgent,
}

func runActionAgent(*cli.Context) error {
	zsf.Run(
		zsf.WithLifeCycles(
			sshagent.NewAgentServer(),
		),
	)
	return nil
}
