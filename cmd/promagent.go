package cmd

import (
	"github.com/LeeZXin/zall/promagent/agent"
	"github.com/urfave/cli/v2"
)

var PromAgent = &cli.Command{
	Name:        "promAgent",
	Usage:       "This command starts prom agent",
	Description: "promAgent provides etcd discovery for prometheus",
	Action:      runPromAgent,
}

func runPromAgent(*cli.Context) error {
	agent.StartAgent()
	return nil
}
