package actionagentapi

import (
	"bytes"
	"github.com/LeeZXin/zall/pkg/git/process"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

var (
	agentToken string
)

func InitApi() {
	agentToken = static.GetString("action.agent.token")
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/actionAgent", checkValid)
		{
			group.POST("/execute", runScript)
		}
	})
}

func runScript(c *gin.Context) {
	var req RunScriptReqVO
	if util.ShouldBindJSON(&req, c) {
		cmd := exec.CommandContext(c, "bash", "-c", req.Script)
		cmd.Env = append(os.Environ(), req.Envs...)
		cmd.Dir = req.Workdir
		stderr := new(bytes.Buffer)
		stdout := new(bytes.Buffer)
		cmd.Stderr = stderr
		cmd.Stdout = stdout
		process.SetSysProcAttribute(cmd)
		cmd.Run()
		c.JSON(http.StatusOK, RunScriptRespVO{
			Stderr: stderr.String(),
			Stdout: stdout.String(),
		})
	}

}

func checkValid(c *gin.Context) {
	if c.GetHeader("Authorization") != agentToken {
		c.String(http.StatusForbidden, "invalid token")
		c.Abort()
		return
	}
	c.Next()
}
