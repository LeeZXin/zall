package repoapi

import (
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitApi() {
	cfgsrv.Inner.InitGitCfg()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/gitRepo", apisession.CheckLogin)
		{
			// 获取模版列表
			group.GET("/allGitIgnoreTemplateList", allGitIgnoreTemplateList)
			// 初始化仓库
			group.POST("/create", createRepo)
			// 删除仓库
			group.POST("/delete", deleteRepo)
			// 展示仓库列表
			group.POST("/list", listRepo)
			// 展示仓库主页
			group.POST("/tree", treeRepo)
			// 展示更多文件列表
			group.POST("/entries", entriesRepo)
			// 展示单个文件内容
			group.POST("/catFile", catFile)
			// 展示仓库所有分支
			group.POST("/allBranches", allBranches)
			// 展示仓库所有tag
			group.POST("/allTags", allTags)
			// gc
			group.POST("/gc", gc)
			// 提交差异
			group.POST("/diffCommits", diffCommits)
			// 展示提交文件差异
			group.POST("/diffFile", diffFile)
			// 展示文件内容
			group.POST("/showDiffTextContent", showDiffTextContent)
			// 历史提交
			group.POST("/historyCommits", historyCommits)
			// 获取令牌列表
			group.POST("/listRepoToken", listRepoToken)
			// 删除访问令牌
			group.POST("/deleteRepoToken", deleteRepoToken)
			// 创建访问令牌
			group.POST("/insertRepoToken", insertRepoToken)
			// 刷新hook
			group.Any("/refreshAllGitHooks", refreshAllGitHooks)
			// 迁移项目组
			group.POST("/transferTeam", transferTeam)
		}
	})
}

func allBranches(c *gin.Context) {
	var req AllBranchesReqVO
	if util.ShouldBindJSON(&req, c) {
		branches, err := reposrv.Outer.AllBranches(c, reposrv.AllBranchesReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     branches,
		})
	}
}

func allTags(c *gin.Context) {
	var req AllTagsReqVO
	if util.ShouldBindJSON(&req, c) {
		branches, err := reposrv.Outer.AllTags(c, reposrv.AllTagsReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     branches,
		})
	}
}

func gc(c *gin.Context) {
	var req GcReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.Gc(c, reposrv.GcReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

// allGitIgnoreTemplateList 获取模版列表
func allGitIgnoreTemplateList(c *gin.Context) {
	templateList := reposrv.Outer.AllGitIgnoreTemplateList()
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     templateList,
	})
}

// treeRepo 代码详情页
func treeRepo(c *gin.Context) {
	var req TreeRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := reposrv.Outer.TreeRepo(c, reposrv.TreeRepoReqDTO{
			Id:       req.Id,
			Ref:      req.Ref,
			Dir:      req.Dir,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, TreeRepoRespVO{
			BaseResp:     ginutil.DefaultSuccessResp,
			ReadmeText:   respDTO.ReadmeText,
			LatestCommit: commitDto2Vo(respDTO.LatestCommit),
			Tree: TreeVO{
				Offset:  respDTO.Tree.Offset,
				Files:   fileDto2Vo(respDTO.Tree.Files),
				Limit:   respDTO.Tree.Limit,
				HasMore: respDTO.Tree.HasMore,
			},
		})
	}
}

// entriesRepo 展示文件列表
func entriesRepo(c *gin.Context) {
	var req EntriesRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		repoRespDTO, err := reposrv.Outer.EntriesRepo(c, reposrv.EntriesRepoReqDTO{
			Id:       req.Id,
			Ref:      req.Ref,
			Dir:      req.Dir,
			Offset:   req.Offset,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, TreeVO{
			Offset:  repoRespDTO.Offset,
			Files:   fileDto2Vo(repoRespDTO.Files),
			Limit:   repoRespDTO.Limit,
			HasMore: repoRespDTO.HasMore,
		})
	}
}

func commitDto2Vo(dto reposrv.CommitDTO) CommitVO {
	return CommitVO{
		Author: UserVO{
			Account: dto.Author.Account,
			Email:   dto.Author.Email,
		},
		Committer: UserVO{
			Account: dto.Committer.Account,
			Email:   dto.Committer.Email,
		},
		AuthoredTime:  util.ReadableTimeComparingNow(time.UnixMilli(dto.AuthoredTime)),
		CommittedTime: util.ReadableTimeComparingNow(time.UnixMilli(dto.CommittedTime)),
		CommitMsg:     dto.CommitMsg,
		CommitId:      dto.CommitId,
		ShortId:       dto.ShortId,
		Verified:      dto.Verified,
	}
}

func fileDto2Vo(list []reposrv.FileDTO) []FileVO {
	ret, _ := listutil.Map(list, func(t reposrv.FileDTO) (FileVO, error) {
		return FileVO{
			Mode:    t.Mode,
			RawPath: t.RawPath,
			Path:    t.Path,
			Commit:  commitDto2Vo(t.Commit),
		}, nil
	})
	return ret
}

func createRepo(c *gin.Context) {
	var req CreateRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.CreateRepo(c, reposrv.CreateRepoReqDTO{
			Operator:      apisession.MustGetLoginUser(c),
			Name:          req.Name,
			Desc:          req.Desc,
			AddReadme:     req.AddReadme,
			GitIgnoreName: req.GitIgnoreName,
			DefaultBranch: req.DefaultBranch,
			TeamId:        req.TeamId,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteRepo(c *gin.Context) {
	var req DeleteRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.DeleteRepo(c, reposrv.DeleteRepoReqDTO{
			Operator: apisession.MustGetLoginUser(c),
			Id:       req.Id,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func listRepo(c *gin.Context) {
	var req ListRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		repoList, err := reposrv.Outer.ListRepo(c, reposrv.ListRepoReqDTO{
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(repoList, func(t repomd.Repo) (RepoVO, error) {
			return RepoVO{
				Id:       t.Id,
				Name:     t.Name,
				Path:     t.Path,
				Author:   t.Author,
				RepoDesc: t.RepoDesc,
				GitSize:  t.GitSize,
				LfsSize:  t.LfsSize,
				Created:  t.Created.Format(time.DateTime),
				Updated:  t.Updated.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]RepoVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func catFile(c *gin.Context) {
	var req CatFileReqVO
	if util.ShouldBindJSON(&req, c) {
		resp, err := reposrv.Outer.CatFile(c, reposrv.CatFileReqDTO{
			Id:       req.Id,
			Ref:      req.Ref,
			Dir:      req.Dir,
			FileName: req.FileName,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, CatFileRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Mode:     resp.ModeName,
			Content:  resp.Content,
		})
	}
}

func diffFile(c *gin.Context) {
	var req DiffFileReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := reposrv.Outer.DiffFile(c, reposrv.DiffFileReqDTO{
			Id:       req.Id,
			Target:   req.Target,
			Head:     req.Head,
			FileName: req.FileName,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := DiffFileRespVO{
			FilePath:    respDTO.FilePath,
			OldMode:     respDTO.OldMode,
			Mode:        respDTO.Mode,
			IsSubModule: respDTO.IsSubModule,
			FileType:    respDTO.FileType,
			IsBinary:    respDTO.IsBinary,
			RenameFrom:  respDTO.RenameFrom,
			RenameTo:    respDTO.RenameTo,
			CopyFrom:    respDTO.CopyFrom,
			CopyTo:      respDTO.CopyTo,
		}
		ret.Lines, _ = listutil.Map(respDTO.Lines, func(t reposrv.DiffLineDTO) (DiffLineVO, error) {
			return DiffLineVO{
				Index:   t.Index,
				LeftNo:  t.LeftNo,
				Prefix:  t.Prefix,
				RightNo: t.RightNo,
				Text:    t.Text,
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func diffCommits(c *gin.Context) {
	var req PrepareMergeReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := reposrv.Outer.DiffCommits(c, reposrv.DiffCommitsReqDTO{
			Id:       req.Id,
			Target:   req.Target,
			Head:     req.Head,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		respVO := PrepareMergeRespVO{
			BaseResp:     ginutil.DefaultSuccessResp,
			Target:       respDTO.Target,
			Head:         respDTO.Head,
			TargetCommit: commitDto2Vo(respDTO.TargetCommit),
			HeadCommit:   commitDto2Vo(respDTO.HeadCommit),
			NumFiles:     respDTO.NumFiles,
			DiffNumsStats: DiffNumsStatInfoVO{
				FileChangeNums: respDTO.DiffNumsStats.FileChangeNums,
				InsertNums:     respDTO.DiffNumsStats.InsertNums,
				DeleteNums:     respDTO.DiffNumsStats.DeleteNums,
			},
			ConflictFiles: respDTO.ConflictFiles,
			CanMerge:      respDTO.CanMerge,
		}
		respVO.Commits, _ = listutil.Map(respDTO.Commits, func(t reposrv.CommitDTO) (CommitVO, error) {
			return commitDto2Vo(t), nil
		})
		respVO.DiffNumsStats.Stats, _ = listutil.Map(respDTO.DiffNumsStats.Stats, func(t reposrv.DiffNumsStatDTO) (DiffNumsStatVO, error) {
			return DiffNumsStatVO{
				RawPath:    t.RawPath,
				Path:       t.Path,
				TotalNums:  t.TotalNums,
				InsertNums: t.InsertNums,
				DeleteNums: t.DeleteNums,
			}, nil
		})
		c.JSON(http.StatusOK, respVO)
	}
}

func showDiffTextContent(c *gin.Context) {
	var req ShowDiffTextContentReqVO
	if util.ShouldBindJSON(&req, c) {
		lines, err := reposrv.Outer.ShowDiffTextContent(c, reposrv.ShowDiffTextContentReqDTO{
			Id:        req.Id,
			CommitId:  req.CommitId,
			FileName:  req.FileName,
			Offset:    req.Offset,
			Limit:     req.Limit,
			Direction: req.Direction,
			Operator:  apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(lines, func(t reposrv.DiffLineDTO) (DiffLineVO, error) {
			return DiffLineVO{
				Index:   t.Index,
				LeftNo:  t.LeftNo,
				Prefix:  t.Prefix,
				RightNo: t.RightNo,
				Text:    t.Text,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]DiffLineVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func historyCommits(c *gin.Context) {
	var req HistoryCommitsReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := reposrv.Outer.HistoryCommits(c, reposrv.HistoryCommitsReqDTO{
			Id:       req.Id,
			Ref:      req.Ref,
			Cursor:   req.Cursor,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := HistoryCommitsRespVO{
			Cursor: respDTO.Cursor,
		}
		ret.Data, _ = listutil.Map(respDTO.Data, func(t reposrv.CommitDTO) (CommitVO, error) {
			return commitDto2Vo(t), nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func listRepoToken(c *gin.Context) {
	var req ListRepoTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		tokens, err := reposrv.Outer.ListRepoToken(c, reposrv.ListRepoTokenReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(tokens, func(t reposrv.RepoTokenDTO) (RepoTokenVO, error) {
			return RepoTokenVO{
				Id:      t.Id,
				Account: t.Account,
				Token:   t.Token,
				Created: t.Created.Format(time.DateTime),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]RepoTokenVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func insertRepoToken(c *gin.Context) {
	var req CreateRepoTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.InsertRepoToken(c, reposrv.InsertRepoTokenReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteRepoToken(c *gin.Context) {
	var req DeleteRepoTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.DeleteRepoToken(c, reposrv.DeleteRepoTokenReqDTO{
			Id:       req.Id,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func refreshAllGitHooks(c *gin.Context) {
	go func() {
		reposrv.Outer.RefreshAllGitHooks(c, reposrv.RefreshAllGitHooksReqDTO{
			Operator: apisession.MustGetLoginUser(c),
		})
	}()
	util.DefaultOkResponse(c)
}

func transferTeam(c *gin.Context) {
	var req TransferTeam
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.TransferTeam(c, reposrv.TransferTeamReqDTO{
			Id:       req.Id,
			TeamId:   req.TeamId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}
