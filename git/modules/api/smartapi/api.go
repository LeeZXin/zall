package smartapi

import (
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/git/modules/service/smartsrv"
	"github.com/LeeZXin/zall/git/modules/service/workflowsrv"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	loginUser = "loginUser"
)

func InitApi() {
	smartsrv.Init()
	reposrv.Init()
	workflowsrv.Init()
	usersrv.Init()
	cfgsrv.Init()
	// smart http协议 不实现dumb协议
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/:corpId/:repoName", handleGoGet, packRepoPath, auth)
		{
			group.POST("/git-upload-pack", uploadPack)
			group.POST("/git-receive-pack", receivePack)
			group.GET("/info/refs", infoRefs)
		}
	})
}

func handleGoGet(c *gin.Context) {
	if c.Query("go-get") == "1" {
		split := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
		if len(split) != 2 {
			c.Next()
			return
		}
		cfg, b := cfgsrv.Inner.GetGitCfg(c)
		if !b {
			c.String(http.StatusInternalServerError, "")
			return
		}
		h, _ := url.Parse(cfg.HttpUrl)
		t := "/" + url.PathEscape(split[0]) + "/" + url.PathEscape(split[1])
		ret := fmt.Sprintf(
			`<meta name="go-import" content="%s">`,
			html.EscapeString(fmt.Sprintf(
				"%s git %s",
				h.Host+t,
				cfg.HttpUrl+t,
			),
			))
		c.String(http.StatusOK, ret)
		c.Abort()
	} else {
		c.Next()
	}
}
func packRepoPath(c *gin.Context) {
	corpId := c.Param("corpId")
	repoName := c.Param("repoName")
	repoPath := filepath.Join(corpId, repoName)
	repo, b := reposrv.Inner.GetByRepoPath(c, repoPath)
	if !b {
		c.String(http.StatusNotFound, "not found")
		c.Abort()
		return
	}
	c.Set("repo", repo)
}

func auth(c *gin.Context) {
	account, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Header("WWW-Authenticate", "Basic realm=\".\"")
		c.String(http.StatusUnauthorized, "wrong authorization")
		c.Abort()
		return
	}
	repo := getRepo(c)
	var (
		userInfo usermd.UserInfo
		b        bool
	)
	if password == "" {
		// 检查是否是工作流的git token
		userInfo, b = workflowsrv.Inner.CheckWorkflowToken(c, repo.Id, account)
		if !b {
			c.Header("WWW-Authenticate", "Basic realm=\".\"")
			c.String(http.StatusUnauthorized, "")
			c.Abort()
			return
		}
	} else {
		userInfo, b = usersrv.Inner.CheckAccountAndPassword(c, usersrv.CheckAccountAndPasswordReqDTO{
			Account:  account,
			Password: password,
		})
		if !b {
			c.Header("WWW-Authenticate", "Basic realm=\".\"")
			c.String(http.StatusUnauthorized, "")
			c.Abort()
			return
		}
	}
	c.Set(loginUser, userInfo)
	c.Next()
}

func uploadPack(c *gin.Context) {
	err := smartsrv.Outer.UploadPack(c, smartsrv.UploadPackReqDTO{
		Repo:     getRepo(c),
		Operator: getUserInfo(c),
		C:        c,
	})
	if err != nil {
		util.HandleApiErr(err, c)
	}
}

func receivePack(c *gin.Context) {
	err := smartsrv.Outer.ReceivePack(c, smartsrv.ReceivePackReqDTO{
		Repo:     getRepo(c),
		Operator: getUserInfo(c),
		C:        c,
	})
	if err != nil {
		util.HandleApiErr(err, c)
	}
}

func infoRefs(c *gin.Context) {
	err := smartsrv.Outer.InfoRefs(c, smartsrv.InfoRefsReqDTO{
		Repo:     getRepo(c),
		Operator: getUserInfo(c),
		C:        c,
	})
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
	}
}

func getRepo(c *gin.Context) repomd.Repo {
	return c.MustGet("repo").(repomd.Repo)
}

func getUserInfo(c *gin.Context) usermd.UserInfo {
	return c.MustGet(loginUser).(usermd.UserInfo)
}
