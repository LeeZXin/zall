package cmd

import (
	"github.com/LeeZXin/zall/git/modules/api/actionapi"
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
	"github.com/LeeZXin/zsf/starter"
	"github.com/urfave/cli/v2"
)

var Actions = &cli.Command{
	Name:        "actions",
	Usage:       "This command starts actions server",
	Description: "zgit actions provides cicd actions",
	Action:      runActions,
}

func runActions(*cli.Context) error {
	// action
	actionapi.InitApi()
	actionsrv.InitSrv()
	starter.Run()
	return nil
}
