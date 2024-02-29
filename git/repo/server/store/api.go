package store

import (
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

var (
	storeSrv = NewStore()
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/v1/git/store")
		{
			group.POST("/initRepo", initRepo)
			group.POST("/delRepo", delRepo)
			group.POST("/getAllBranches", getAllBranches)
			group.POST("/getAllTags", getAllTags)
			group.POST("/gc", gc)
			group.POST("/diffRefs", diffRefs)
			group.POST("/diffFile", diffFile)
			group.POST("/getRepoSize", getRepoSize)
			group.POST("/showDiffTextContent", showDiffTextContent)
			group.POST("/historyCommits", historyCommits)
			group.POST("/initRepoHook", initRepoHook)
			group.POST("/entriesRepo", entriesRepo)
			group.POST("/catFile", catFile)
			group.POST("/treeRepo", treeRepo)
			group.POST("/merge", merge)
		}
		group = e.Group("/api/v1/git/smart/:corpId/:repoName", packRepoPath)
		{
			group.POST("/git-upload-pack", uploadPack)
			group.POST("/git-receive-pack", receivePack)
			group.GET("/info/refs", infoRefs)
		}
	})
}

func packRepoPath(c *gin.Context) {
	c.Set("repoPath", filepath.Join(c.Param("corpId"), c.Param("repoName")))
}

func initRepo(c *gin.Context) {
	var req reqvo.InitRepoReq
	if util.ShouldBindJSON(&req, c) {
		err := storeSrv.InitRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func delRepo(c *gin.Context) {
	var req reqvo.DeleteRepoReq
	if util.ShouldBindJSON(&req, c) {
		err := storeSrv.DeleteRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getAllBranches(c *gin.Context) {
	var req reqvo.GetAllBranchesReq
	if util.ShouldBindJSON(&req, c) {
		ret, err := storeSrv.GetAllBranches(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}

func getAllTags(c *gin.Context) {
	var req reqvo.GetAllTagsReq
	if util.ShouldBindJSON(&req, c) {
		ret, err := storeSrv.GetAllTags(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}

func gc(c *gin.Context) {
	var req reqvo.GcReq
	if util.ShouldBindJSON(&req, c) {
		err := storeSrv.Gc(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func diffRefs(c *gin.Context) {
	var req reqvo.DiffRefsReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.DiffRefs(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.DiffRefsResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func diffFile(c *gin.Context) {
	var req reqvo.DiffFileReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.DiffFile(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.DiffFileResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func getRepoSize(c *gin.Context) {
	var req reqvo.GetRepoSizeReq
	if util.ShouldBindJSON(&req, c) {
		size, err := storeSrv.GetRepoSize(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[int64]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     size,
		})
	}
}

func showDiffTextContent(c *gin.Context) {
	var req reqvo.ShowDiffTextContentReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.ShowDiffTextContent(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]reqvo.DiffLineVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func historyCommits(c *gin.Context) {
	var req reqvo.HistoryCommitsReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.HistoryCommits(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.HistoryCommitsResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func initRepoHook(c *gin.Context) {
	var req reqvo.InitRepoHookReq
	if util.ShouldBindJSON(&req, c) {
		err := storeSrv.InitRepoHook(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func entriesRepo(c *gin.Context) {
	var req reqvo.EntriesRepoReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.EntriesRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.TreeVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func catFile(c *gin.Context) {
	var req reqvo.CatFileReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.CatFile(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.CatFileResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func treeRepo(c *gin.Context) {
	var req reqvo.TreeRepoReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := storeSrv.TreeRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.TreeRepoResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func uploadPack(c *gin.Context) {
	storeSrv.UploadPack(reqvo.UploadPackReq{
		RepoPath: c.GetString("repoPath"),
		C:        c,
	})
}

func receivePack(c *gin.Context) {
	storeSrv.ReceivePack(reqvo.ReceivePackReq{
		RepoPath: c.GetString("repoPath"),
		C:        c,
	})
}

func infoRefs(c *gin.Context) {
	storeSrv.InfoRefs(c.Request.Context(), reqvo.InfoRefsReq{
		RepoPath: c.GetString("repoPath"),
		C:        c,
	})
}

func merge(c *gin.Context) {
	var req reqvo.MergeReq
	if util.ShouldBindJSON(&req, c) {
		err := storeSrv.Merge(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
