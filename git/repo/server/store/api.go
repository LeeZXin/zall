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

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/v1/git/store")
		{
			group.POST("/initRepo", initRepo)
			group.POST("/deleteRepo", deleteRepo)
			group.POST("/getAllBranches", getAllBranches)
			group.POST("/listBranchAndLastCommit", listBranchAndLastCommit)
			group.POST("/listTagAndCommit", listTagAndCommit)
			group.POST("/deleteBranch", deleteBranch)
			group.POST("/getAllTags", getAllTags)
			group.POST("/gc", gc)
			group.POST("/diffRefs", diffRefs)
			group.POST("/diffCommits", diffCommits)
			group.POST("/diffFile", diffFile)
			group.POST("/showDiffTextContent", showDiffTextContent)
			group.POST("/historyCommits", historyCommits)
			group.POST("/initRepoHook", initRepoHook)
			group.POST("/entriesRepo", entriesRepo)
			group.POST("/catFile", catFile)
			group.POST("/indexRepo", indexRepo)
			group.POST("/merge", merge)
			group.POST("/blame", blame)
			group.POST("/canMerge", canMerge)
			group.GET("/archive/:corpId/:repoName/:fileName", packRepoPath, createArchive)
			group.POST("/deleteTag", deleteTag)
		}
		group = e.Group("/api/v1/git/smart/:corpId/:repoName", packRepoPath)
		{
			group.POST("/git-upload-pack", uploadPack)
			group.POST("/git-receive-pack", receivePack)
			group.GET("/info/refs", infoRefs)
		}
	})
}

func createArchive(c *gin.Context) {
	CreateArchive(c, reqvo.CreateArchiveReq{
		RepoPath: c.GetString("repoPath"),
		FileName: c.Param("fileName"),
		C:        c,
	})
}

func packRepoPath(c *gin.Context) {
	c.Set("repoPath", filepath.Join(c.Param("corpId"), c.Param("repoName")))
}

func initRepo(c *gin.Context) {
	var req reqvo.InitRepoReq
	if util.ShouldBindJSON(&req, c) {
		gitSize, err := InitRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[int64]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     gitSize,
		})
	}
}

func deleteRepo(c *gin.Context) {
	var req reqvo.DeleteRepoReq
	if util.ShouldBindJSON(&req, c) {
		err := DeleteRepo(c, req)
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
		ret, err := GetAllBranches(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]reqvo.RefVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}

func deleteBranch(c *gin.Context) {
	var req reqvo.DeleteBranchReq
	if util.ShouldBindJSON(&req, c) {
		err := DeleteBranch(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listBranchAndLastCommit(c *gin.Context) {
	var req reqvo.ListRefCommitsReq
	if util.ShouldBindJSON(&req, c) {
		ret, total, err := ListBranchAndLastCommit(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.Page2Resp[reqvo.RefCommitVO]{
			DataResp: ginutil.DataResp[[]reqvo.RefCommitVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     ret,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func listTagAndCommit(c *gin.Context) {
	var req reqvo.ListRefCommitsReq
	if util.ShouldBindJSON(&req, c) {
		ret, total, err := ListTagAndCommit(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.Page2Resp[reqvo.RefCommitVO]{
			DataResp: ginutil.DataResp[[]reqvo.RefCommitVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     ret,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func getAllTags(c *gin.Context) {
	var req reqvo.GetAllTagsReq
	if util.ShouldBindJSON(&req, c) {
		ret, err := GetAllTags(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]reqvo.RefVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}

func gc(c *gin.Context) {
	var req reqvo.GcReq
	if util.ShouldBindJSON(&req, c) {
		gitSize, err := Gc(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.GcResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data: reqvo.GcResp{
				GitSize: gitSize,
			},
		})
	}
}

func diffRefs(c *gin.Context) {
	var req reqvo.DiffRefsReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := DiffRefs(c, req)
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

func diffCommits(c *gin.Context) {
	var req reqvo.DiffCommitsReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := DiffCommits(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.DiffCommitsResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func canMerge(c *gin.Context) {
	var req reqvo.CanMergeReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := CanMerge(c, req)
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

func diffFile(c *gin.Context) {
	var req reqvo.DiffFileReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := DiffFile(c, req)
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

func showDiffTextContent(c *gin.Context) {
	var req reqvo.ShowDiffTextContentReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := ShowDiffTextContent(c, req)
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
		resp, err := HistoryCommits(c, req)
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
		err := InitRepoHook(c, req)
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
		resp, err := EntriesRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]reqvo.BlobVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func catFile(c *gin.Context) {
	var req reqvo.CatFileReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := CatFile(c, req)
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

func indexRepo(c *gin.Context) {
	var req reqvo.IndexRepoReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := IndexRepo(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[reqvo.IndexRepoResp]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     resp,
		})
	}
}

func uploadPack(c *gin.Context) {
	UploadPack(reqvo.UploadPackReq{
		RepoPath: c.GetString("repoPath"),
		C:        c,
	})
}

func receivePack(c *gin.Context) {
	ReceivePack(reqvo.ReceivePackReq{
		RepoPath: c.GetString("repoPath"),
		C:        c,
	})
}

func infoRefs(c *gin.Context) {
	InfoRefs(c.Request.Context(), reqvo.InfoRefsReq{
		RepoPath: c.GetString("repoPath"),
		C:        c,
	})
}

func merge(c *gin.Context) {
	var req reqvo.MergeReq
	if util.ShouldBindJSON(&req, c) {
		resp, err := Merge(c, req)
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

func blame(c *gin.Context) {
	var req reqvo.BlameReq
	if util.ShouldBindJSON(&req, c) {
		lines, err := Blame(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]reqvo.BlameLineVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     lines,
		})
	}
}

func deleteTag(c *gin.Context) {
	var req reqvo.DeleteTagReqVO
	if util.ShouldBindJSON(&req, c) {
		err := DeleteTag(c, req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
