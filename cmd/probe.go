package cmd

import (
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var Probe = &cli.Command{
	Name:        "probe",
	Usage:       "This command starts service probe",
	Description: "this provides service probe",
	Action:      runProbe,
}

func runProbe(*cli.Context) error {
	//deploysrv.InitProbeTask()
	zsf.Run()
	return nil
}
