package cmd

import (
	"github.com/LeeZXin/zall/fileserv/modules/api/fileapi"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var FileServer = &cli.Command{
	Name:        "fileServer",
	Usage:       "This command starts file server",
	Description: "this provides file service",
	Action:      runFileServer,
}

func runFileServer(*cli.Context) error {
	fileapi.InitApi()
	zsf.Run(
		zsf.WithLifeCycles(
			httpserver.NewServer(),
		),
	)
	return nil
}
