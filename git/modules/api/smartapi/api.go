package smartapi

import (
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/git/modules/service/smartsrv"
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
	fromAccessToken = "fromAccessToken"
	loginUser       = "loginUser"
)

func InitApi() {
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
		h, _ := url.Parse(cfg.AppUrl)
		t := "/" + url.PathEscape(split[0]) + "/" + url.PathEscape(split[1])
		ret := fmt.Sprintf(
			`<meta name="go-import" content="%s">`,
			html.EscapeString(fmt.Sprintf(
				"%s git %s",
				h.Host+t,
				cfg.AppUrl+t,
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
	repo, b := reposrv.Inner.GetByRepoPath(c.Request.Context(), repoPath)
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
	userInfo, b := usersrv.Inner.CheckAccountAndPassword(c, usersrv.CheckAccountAndPasswordReqDTO{
		Account:  account,
		Password: password,
	})
	if !b {
		repo := getRepo(c)
		// 常规账号密码不存在就检查访问令牌
		b = reposrv.Inner.CheckAccessToken(c.Request.Context(), reposrv.CheckAccessTokenReqDTO{
			Id:      repo.Id,
			Account: account,
			Token:   password,
		})
		if !b {
			c.Header("WWW-Authenticate", "Basic realm=\".\"")
			c.String(http.StatusUnauthorized, "")
			c.Abort()
			return
		}
		userInfo = usermd.UserInfo{
			Account: account,
			Name:    fmt.Sprintf("%s's accessToken", repo.Name),
			Email:   "zgit@noreply.fake",
		}
		c.Set(fromAccessToken, true)
	} else {
		c.Set(fromAccessToken, false)
	}
	c.Set(loginUser, userInfo)
	c.Next()
}

func uploadPack(c *gin.Context) {
	err := smartsrv.Outer.UploadPack(c.Request.Context(), smartsrv.UploadPackReqDTO{
		Repo:            getRepo(c),
		Operator:        getUserInfo(c),
		C:               c,
		FromAccessToken: getFromAccessToken(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
	}
}

func receivePack(c *gin.Context) {
	err := smartsrv.Outer.ReceivePack(c.Request.Context(), smartsrv.ReceivePackReqDTO{
		Repo:     getRepo(c),
		Operator: getUserInfo(c),
		C:        c,
	})
	if err != nil {
		util.HandleApiErr(err, c)
	}
}

func infoRefs(c *gin.Context) {
	err := smartsrv.Outer.InfoRefs(c.Request.Context(), smartsrv.InfoRefsReqDTO{
		Repo:            getRepo(c),
		Operator:        getUserInfo(c),
		FromAccessToken: getFromAccessToken(c),
		C:               c,
	})
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
	}
}

func getRepo(c *gin.Context) repomd.RepoInfo {
	return c.MustGet("repo").(repomd.RepoInfo)
}

func getFromAccessToken(c *gin.Context) bool {
	return c.GetBool(fromAccessToken)
}

func getUserInfo(c *gin.Context) usermd.UserInfo {
	return c.MustGet(loginUser).(usermd.UserInfo)
}
