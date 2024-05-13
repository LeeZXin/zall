package cmd

import (
	"github.com/urfave/cli/v2"
	"runtime"
)

var (
	cmdList = []*cli.Command{
		Run,
		RepoServer,
		GitHook,
		WorkflowAgentServer,
		Timer,
		PropServer,
		TcpDetect,
		FileServer,
		Probe,
		DeployServer,
		PromAgent,
	}
)

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.HideHelp = true
	app.DefaultCommand = Run.Name
	app.Commands = append(app.Commands, cmdList...)
	app.Name = "zall"
	app.Usage = "A zall server with zsf"
	app.Description = "by default, it will start the zall server"
	app.Version = formatBuiltWith()
	return app
}

func formatBuiltWith() string {
	return " built with " + runtime.Version()
}
