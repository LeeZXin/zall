package cmd

import (
	"github.com/LeeZXin/zall/genid/modules/api/idapi"
	"github.com/LeeZXin/zall/git/modules/api/actionapi"
	"github.com/LeeZXin/zall/git/modules/api/branchapi"
	"github.com/LeeZXin/zall/git/modules/api/gpgkeyapi"
	"github.com/LeeZXin/zall/git/modules/api/lfsapi"
	"github.com/LeeZXin/zall/git/modules/api/pullrequestapi"
	"github.com/LeeZXin/zall/git/modules/api/repoapi"
	"github.com/LeeZXin/zall/git/modules/api/smartapi"
	"github.com/LeeZXin/zall/git/modules/api/sshkeyapi"
	"github.com/LeeZXin/zall/git/modules/api/webhookapi"
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
	"github.com/LeeZXin/zall/git/modules/sshproxy"
	reposerver "github.com/LeeZXin/zall/git/repo/server"
	"github.com/LeeZXin/zall/meta/modules/api/cfgapi"
	"github.com/LeeZXin/zall/meta/modules/api/teamapi"
	"github.com/LeeZXin/zall/meta/modules/api/userapi"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/timer/modules/api/taskapi"
	"github.com/LeeZXin/zall/timer/modules/service/tasksrv"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
)

var Run = &cli.Command{
	Name:        "run",
	Usage:       "This command starts zall server",
	Description: "this server provides zall service",
	Action:      runZall,
}

func runZall(*cli.Context) error {
	lifeCycles := make([]zsf.LifeCycle, 0)
	// for meta
	{
		userapi.InitApi()
		teamapi.InitApi()
		cfgapi.InitApi()
	}
	// for git
	{
		git.Init()
		lfsapi.InitApi()
		branchapi.InitApi()
		pullrequestapi.InitApi()
		smartapi.InitApi()
		repoapi.InitApi()
		sshkeyapi.InitApi()
		gpgkeyapi.InitApi()
		webhookapi.InitApi()
		lifeCycles = append(lifeCycles, sshproxy.InitProxy())
		if static.GetBool("git.repo.server.enabled") {
			logger.Logger.Info("git repo server enabled")
			reposerver.InitHttpApi()
			lifeCycles = append(lifeCycles, reposerver.InitSshServer())
		}
		if static.GetBool("actions.enabled") {
			logger.Logger.Info("git actions server enabled")
			actionapi.InitApi()
			actionsrv.InitSrv()
		}
	}
	// for timer
	{
		taskapi.InitApi()
		if static.GetBool("timer.enabled") {
			logger.Logger.Info("timer executor enabled")
			tasksrv.InitTask()
		}
	}
	// for idserver
	{
		if static.GetBool("idserver.enabled") {
			logger.Logger.Info("id server enabled")
			idapi.InitApi()
		}
	}
	lifeCycles = append(lifeCycles, httpserver.NewServer())
	zsf.Run(
		zsf.WithLifeCycles(lifeCycles...),
	)
	return nil
}
