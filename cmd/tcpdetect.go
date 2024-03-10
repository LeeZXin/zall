package cmd

import (
	"github.com/LeeZXin/zall/tcpdetect/modules/service/detectsrv"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var TcpDetect = &cli.Command{
	Name:        "tcpDetect",
	Usage:       "This command starts tcp detect server",
	Description: "this provides tcp detect service",
	Action:      runTcpDetect,
}

func runTcpDetect(*cli.Context) error {
	{
		detectsrv.InitDetect()
	}
	zsf.Run()
	return nil
}
