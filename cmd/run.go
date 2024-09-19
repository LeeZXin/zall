package cmd

import (
	"github.com/LeeZXin/zall/alert/modules/api/alertapi"
	"github.com/LeeZXin/zall/alert/modules/service/alertsrv"
	"github.com/LeeZXin/zall/dbaudit/modules/api/mysqldbapi"
	"github.com/LeeZXin/zall/deploy/modules/api/deployapi"
	"github.com/LeeZXin/zall/discovery/modules/api/discoveryapi"
	"github.com/LeeZXin/zall/fileserv/modules/api/fileapi"
	"github.com/LeeZXin/zall/git/modules/api/branchapi"
	"github.com/LeeZXin/zall/git/modules/api/gpgkeyapi"
	"github.com/LeeZXin/zall/git/modules/api/lfsapi"
	"github.com/LeeZXin/zall/git/modules/api/pullrequestapi"
	"github.com/LeeZXin/zall/git/modules/api/repoapi"
	"github.com/LeeZXin/zall/git/modules/api/smartapi"
	"github.com/LeeZXin/zall/git/modules/api/sshkeyapi"
	"github.com/LeeZXin/zall/git/modules/api/webhookapi"
	"github.com/LeeZXin/zall/git/modules/api/workflowapi"
	"github.com/LeeZXin/zall/git/modules/sshproxy"
	reposerver "github.com/LeeZXin/zall/git/repo/server"
	"github.com/LeeZXin/zall/gw"
	"github.com/LeeZXin/zall/meta/modules/api/appapi"
	"github.com/LeeZXin/zall/meta/modules/api/cfgapi"
	"github.com/LeeZXin/zall/meta/modules/api/teamapi"
	"github.com/LeeZXin/zall/meta/modules/api/userapi"
	"github.com/LeeZXin/zall/meta/modules/api/zalletapi"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/notify/modules/api/notifyapi"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/promagent/agent"
	"github.com/LeeZXin/zall/promagent/modules/api/promapi"
	"github.com/LeeZXin/zall/property/modules/api/propertyapi"
	"github.com/LeeZXin/zall/teamhook/modules/api/teamhookapi"
	"github.com/LeeZXin/zall/thirdpart/modules/api/tpfeishuapi"
	"github.com/LeeZXin/zall/thirdpart/modules/api/tpweworkapi"
	"github.com/LeeZXin/zall/thirdpart/modules/service/tpfeishusrv"
	"github.com/LeeZXin/zall/thirdpart/modules/service/tpweworksrv"
	"github.com/LeeZXin/zall/timer/modules/api/timerapi"
	"github.com/LeeZXin/zall/timer/modules/service/timersrv"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/services/registry"
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
	// 前端路由
	gw.Init()
	// for apisession
	{
		apisession.Init()
	}
	// for envs
	{
		cfgsrv.InitEnvCfg()
	}
	lifeCycles := make([]zsf.LifeCycle, 0)
	// for meta
	{
		userapi.InitApi()
		teamapi.InitApi()
		cfgapi.InitApi()
		appapi.InitApi()
	}
	// for cfg
	{
		// 初始化全局配置
		cfgsrv.InitSysCfg()
		cfgsrv.InitLoginCfg()
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
	}
	// for workflow
	{
		workflowapi.InitApi()
	}
	// for timer
	{
		timerapi.InitApi()
		timersrv.InitTask()
	}
	// for property
	{
		propertyapi.InitApi()
	}
	// for files server
	{
		if static.GetBool("files.enabled") {
			logger.Logger.Info("file server enabled")
			fileapi.InitApi()
		}
	}
	// for deploy
	{
		deployapi.InitApi()
	}
	// prom
	{
		promapi.InitApi()
		if static.GetBool("prom.agent.enabled") {
			logger.Logger.Info("prom agent enabled")
			agent.StartAgent()
		}
	}
	// for db
	{
		mysqldbapi.InitApi()
	}
	// for alert
	{
		alertapi.InitApi()
		if static.GetBool("alert.enabled") {
			alertsrv.InitTask()
		}
	}
	// for zallet
	{
		zalletapi.InitApi()
	}
	// for discovery
	{
		discoveryapi.InitApi()
	}
	// for notify
	{
		notifyapi.InitApi()
	}
	// for team hook
	{
		teamhookapi.InitApi()
	}
	// for wework access token
	{
		tpweworkapi.InitApi()
		if static.GetBool("wework.accessToken.task.enabled") {
			logger.Logger.Info("wework access token task enabled")
			tpweworksrv.InitGetAccessTokenTask()
		}
	}
	// for feishu access token
	{
		tpfeishuapi.InitApi()
		if static.GetBool("feishu.accessToken.task.enabled") {
			logger.Logger.Info("feishu access token task enabled")
			tpfeishusrv.InitGetAccessTokenTask()
		}
	}
	lifeCycles = append(lifeCycles,
		httpserver.NewServer(
			httpserver.WithRegistry(
				registry.NewDefaultEtcdRegistry(),
			),
			httpserver.WithEnablePromApi(true),
		),
	)
	zsf.Run(
		zsf.WithLifeCycles(lifeCycles...),
	)
	return nil
}
