package cmd

import (
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
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/starter"
	"github.com/urfave/cli/v2"
)

var Run = &cli.Command{
	Name:        "run",
	Usage:       "This command starts zall server",
	Description: "this server provides zall service",
	Action:      runZall,
}

func runZall(*cli.Context) error {
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
		sshproxy.InitProxy()
		if static.GetBool("git.repo.server.enabled") {
			logger.Logger.Info("git repo server enabled")
			reposerver.InitHttpApi()
			reposerver.InitSshServer()
		}
		if static.GetBool("actions.enabled") {
			logger.Logger.Info("git actions server enabled")
			actionapi.InitApi()
			actionsrv.InitSrv()
		}
	}
	logger.Logger.Info("start zall server successfully")
	starter.Run()
	return nil
}
