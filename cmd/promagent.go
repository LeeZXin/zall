package cmd

import (
	"github.com/LeeZXin/zall/promagent/agent"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"net/http"
)

var PromAgent = &cli.Command{
	Name:        "promAgent",
	Usage:       "This command starts prom agent",
	Description: "promAgent provides etcd discovery for prometheus",
	Action:      runPromAgent,
}

func runPromAgent(*cli.Context) error {
	agent.StartAgent()
	// 探针
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		e.Any("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "")
		})
	})
	zsf.Run(
		zsf.WithLifeCycles(httpserver.NewServer()),
	)
	return nil
}
