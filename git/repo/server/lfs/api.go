package lfs

import (
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/git/repo/server/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"regexp"
)

var (
	oidPattern = regexp.MustCompile(`^[a-f\d]{64}$`)

	lfsSrv Lfs
)

func InitApi() {
	lfsSrv = NewLfs()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/v1/lfs/prop", apisession.CheckToken)
		{
			group.POST("/stat", lfsStat)
			group.POST("/exists", lfsExists)
			group.POST("/batchExists", lfsBatchExists)
		}
		group = e.Group("/api/v1/lfs/file/:corpId/:repoName/:oid", apisession.CheckToken, packRepoPath, packOid)
		{
			group.PUT("/upload", lfsUpload)
			group.GET("/download", lfsDownload)
		}
	})
}

func packOid(c *gin.Context) {
	oid := c.Param("oid")
	if !oidPattern.MatchString(oid) {
		c.String(http.StatusBadRequest, "")
		c.Abort()
	}
	c.Set("oid", c.Param("oid"))
}

func lfsUpload(c *gin.Context) {
	lfsSrv.Upload(c, reqvo.LfsUploadReq{
		RepoPath: c.GetString("repoPath"),
		Oid:      c.GetString("oid"),
		C:        c,
	})
}

func lfsExists(c *gin.Context) {
	var req reqvo.LfsExistsReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := lfsSrv.Exists(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[bool]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func lfsBatchExists(c *gin.Context) {
	var req reqvo.LfsBatchExistsReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := lfsSrv.BatchExists(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[map[string]bool]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func lfsStat(c *gin.Context) {
	var req reqvo.LfsStatReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := lfsSrv.Stat(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.LfsStatResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func lfsDownload(c *gin.Context) {
	lfsSrv.Download(c, reqvo.LfsDownloadReq{
		RepoPath: c.GetString("repoPath"),
		Oid:      c.GetString("oid"),
		C:        c,
	})
}

func packRepoPath(c *gin.Context) {
	c.Set("repoPath", filepath.Join(c.Param("corpId"), c.Param("repoName")))
}
