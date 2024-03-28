package cmd

import (
	"github.com/LeeZXin/zall/services/modules/service/deploysrv"
	"github.com/LeeZXin/zsf/property/static"
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
	deploysrv.InitProbeTask(static.GetString("probe.env"))
	zsf.Run()
	return nil
}
