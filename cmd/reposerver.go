package cmd

import (
	"github.com/LeeZXin/zall/git/repo/server"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var RepoServer = &cli.Command{
	Name:        "repoServer",
	Usage:       "This command starts git repo server",
	Description: "this provides repo service, should call rpc from zall server",
	Action:      runRepoServer,
}

func runRepoServer(*cli.Context) error {
	{
		git.Init()
		server.InitHttpApi()
	}
	zsf.Run(
		zsf.WithLifeCycles(
			server.InitSshServer(),
			httpserver.NewServer(),
		),
	)
	return nil
}
