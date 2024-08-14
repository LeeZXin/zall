package repoapi

import (
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func InitApi() {
	cfgsrv.InitInner()
	cfgsrv.Inner.InitGitCfg()
	reposrv.Init()
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/gitRepo", apisession.CheckLogin)
		{
			// 仓库信息+权限
			group.GET("/get/:repoId", getRepo)
			// 基本信息
			group.GET("/simpleInfo/:repoId", getSimpleInfo)
			// 详细信息
			group.GET("/detail/:repoId", getDetailInfo)
			// 获取模版列表
			group.GET("/allGitIgnoreTemplateList", allGitIgnoreTemplateList)
			// 初始化仓库
			group.POST("/create", createRepo)
			// 编辑仓库
			group.POST("/update", updateRepo)
			// 删除仓库
			group.DELETE("/delete/:repoId", deleteRepo)
			// 永久删除仓库
			group.DELETE("/deletePermanently/:repoId", deletePermanentlyRepo)
			// 从回收站恢复仓库
			group.PUT("/recoverFromRecycle/:repoId", recoverFromRecycle)
			// 展示仓库列表
			group.GET("/list/:teamId", listRepo)
			// 展示已删除仓库列表
			group.GET("/listDeleted/:teamId", listDeletedRepo)
			// 展示仓库主页
			group.GET("/index", indexRepo)
			// 展示更多文件列表
			group.GET("/entries", entriesRepo)
			// 展示单个文件内容
			group.GET("/catFile", catFile)
			// 展示仓库所有分支
			group.GET("/allBranches/:repoId", allBranches)
			// 分页展示分支+提交
			group.GET("/pageBranchCommits", pageBranchCommits)
			// 分页展示分支+提交
			group.GET("/pageTagCommits", pageTagCommits)
			// 展示仓库所有tag
			group.GET("/allTags/:repoId", allTags)
			// gc
			group.PUT("/gc/:repoId", gc)
			// 提交差异
			group.GET("/diffRefs", diffRefs)
			// 提交差异
			group.GET("/diffCommits", diffCommits)
			// 展示提交文件差异
			group.GET("/diffFile", diffFile)
			// 历史提交
			group.GET("/historyCommits", historyCommits)
			// 迁移团队
			group.PUT("/transferTeam", transferTeam)
			// 获取每一行提交信息
			group.GET("/blame", blame)
			// 删除分支
			group.DELETE("/deleteBranch", deleteBranch)
			// 下载代码压缩包
			group.GET("/archive", createArchive)
			// 删除tag
			group.DELETE("/deleteTag", deleteTag)
			// 归档仓库
			group.PUT("/setArchived/:repoId", setArchived)
			// 归档仓库变为正常仓库
			group.PUT("/setUnArchived/:repoId", setUnArchived)
			// 管理员查看仓库列表
			group.GET("/listByAdmin/:teamId", listRepoByAdmin)
		}
	})
}

func listRepoByAdmin(c *gin.Context) {
	repos, err := reposrv.Outer.ListRepoByAdmin(c, reposrv.ListRepoByAdminReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(repos, func(t reposrv.SimpleRepoDTO) (SimpleRepoVO, error) {
		return SimpleRepoVO{
			RepoId: t.RepoId,
			Name:   t.Name,
			TeamId: t.TeamId,
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]SimpleRepoVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func setArchived(c *gin.Context) {
	err := reposrv.Outer.SetRepoArchivedStatus(c, reposrv.SetRepoArchivedStatusReqDTO{
		RepoId:     getRepoId(c),
		IsArchived: true,
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func setUnArchived(c *gin.Context) {
	err := reposrv.Outer.SetRepoArchivedStatus(c, reposrv.SetRepoArchivedStatusReqDTO{
		RepoId:     getRepoId(c),
		IsArchived: false,
		Operator:   apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func getRepo(c *gin.Context) {
	repo, perm, err := reposrv.Outer.GetRepoAndPerm(c, reposrv.GetRepoAndPermReqDTO{
		RepoId:   getRepoId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[RepoWithPermVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: RepoWithPermVO{
			SimpleRepoVO: SimpleRepoVO{
				RepoId: repo.RepoId,
				Name:   repo.Name,
				TeamId: repo.TeamId,
			},
			Perm: perm,
		},
	})
}

func updateRepo(c *gin.Context) {
	var req UpdateRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.UpdateRepo(c, reposrv.UpdateRepoReqDTO{
			RepoId:       req.RepoId,
			Desc:         req.Desc,
			DisableLfs:   req.DisableLfs,
			LfsLimitSize: req.LfsLimitSize,
			GitLimitSize: req.GitLimitSize,
			Operator:     apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func getSimpleInfo(c *gin.Context) {
	info, err := reposrv.Outer.GetSimpleInfo(c, reposrv.GetSimpleInfoReqDTO{
		RepoId:   getRepoId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[SimpleInfoVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data: SimpleInfoVO{
			Branches:     info.Branches,
			Tags:         info.Tags,
			CloneHttpUrl: info.CloneHttpUrl,
			CloneSshUrl:  info.CloneSshUrl,
		},
	})
}

func getDetailInfo(c *gin.Context) {
	info, err := reposrv.Outer.GetDetailInfo(c, reposrv.GetDetailInfoReqDTO{
		RepoId:   getRepoId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, ginutil.DataResp[RepoVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     repo2VO(info),
	})
}

func deleteBranch(c *gin.Context) {
	var req DeleteBranchReqVO
	if util.ShouldBindQuery(&req, c) {
		err := reposrv.Outer.DeleteBranch(c, reposrv.DeleteBranchReqDTO{
			RepoId:   req.RepoId,
			Branch:   req.Branch,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func deleteTag(c *gin.Context) {
	var req DeleteTagReqVO
	if util.ShouldBindQuery(&req, c) {
		err := reposrv.Outer.DeleteTag(c, reposrv.DeleteTagReqDTO{
			RepoId:   req.RepoId,
			Tag:      req.Tag,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		util.DefaultOkResponse(c)
	}
}

func allBranches(c *gin.Context) {
	branches, err := reposrv.Outer.AllBranches(c, reposrv.AllBranchesReqDTO{
		RepoId:   getRepoId(c),
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

func pageBranchCommits(c *gin.Context) {
	var req PageRefCommitsReqVO
	if util.ShouldBindQuery(&req, c) {
		branches, total, err := reposrv.Outer.PageBranchCommits(c, reposrv.PageRefCommitsReqDTO{
			RepoId:   req.RepoId,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(branches, func(t reposrv.BranchCommitDTO) (BranchCommitVO, error) {
			ret := BranchCommitVO{
				Name:              t.Name,
				IsProtectedBranch: t.IsProtectedBranch,
				LastCommit:        commitDto2Vo(t.LastCommit),
			}
			if t.LastPullRequest != nil {
				ret.LastPullRequest = &PullRequestVO{
					Id:       t.LastPullRequest.Id,
					PrStatus: t.LastPullRequest.PrStatus,
					PrTitle:  t.LastPullRequest.PrTitle,
					Created:  t.LastPullRequest.Created.Format(time.DateTime),
				}
			}
			return ret, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[BranchCommitVO]{
			DataResp: ginutil.DataResp[[]BranchCommitVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func pageTagCommits(c *gin.Context) {
	var req PageRefCommitsReqVO
	if util.ShouldBindQuery(&req, c) {
		tags, total, err := reposrv.Outer.PageTagCommits(c, reposrv.PageRefCommitsReqDTO{
			RepoId:   req.RepoId,
			PageNum:  req.PageNum,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(tags, func(t reposrv.TagCommitDTO) (TagCommitVO, error) {
			return TagCommitVO{
				Name:   t.Name,
				Commit: commitDto2Vo(t.Commit),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.Page2Resp[TagCommitVO]{
			DataResp: ginutil.DataResp[[]TagCommitVO]{
				BaseResp: ginutil.DefaultSuccessResp,
				Data:     data,
			},
			PageNum:    req.PageNum,
			TotalCount: total,
		})
	}
}

func allTags(c *gin.Context) {
	branches, err := reposrv.Outer.AllTags(c, reposrv.AllTagsReqDTO{
		RepoId:   getRepoId(c),
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

func gc(c *gin.Context) {
	err := reposrv.Outer.Gc(c, reposrv.GcReqDTO{
		RepoId:   getRepoId(c),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

// allGitIgnoreTemplateList 获取模版列表
func allGitIgnoreTemplateList(c *gin.Context) {
	templateList := reposrv.Outer.AllGitIgnoreTemplateList()
	c.JSON(http.StatusOK, ginutil.DataResp[[]string]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     templateList,
	})
}

// indexRepo 代码详情页
func indexRepo(c *gin.Context) {
	var req IndexRepoReqVO
	if util.ShouldBindQuery(&req, c) {
		respDTO, err := reposrv.Outer.IndexRepo(c, reposrv.IndexRepoReqDTO{
			RepoId:   req.RepoId,
			Ref:      req.Ref,
			RefType:  req.RefType,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, IndexRepoRespVO{
			BaseResp:     ginutil.DefaultSuccessResp,
			ReadmeText:   respDTO.ReadmeText,
			HasReadme:    respDTO.HasReadme,
			LatestCommit: commitDto2Vo(respDTO.LatestCommit),
			Tree: TreeVO{
				Files: fileDto2Vo(respDTO.Tree.Files),
			},
		})
	}
}

// entriesRepo 展示文件列表
func entriesRepo(c *gin.Context) {
	var req EntriesRepoReqVO
	if util.ShouldBindQuery(&req, c) {
		blobs, err := reposrv.Outer.EntriesRepo(c, reposrv.EntriesRepoReqDTO{
			RepoId:   req.RepoId,
			Ref:      req.Ref,
			Dir:      req.Dir,
			RefType:  req.RefType,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(blobs, func(t reposrv.BlobDTO) (BlobVO, error) {
			return BlobVO{
				Mode:    t.Mode,
				RawPath: t.RawPath,
				Path:    t.Path,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]BlobVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func commitDto2Vo(d reposrv.CommitDTO) CommitVO {
	ret := CommitVO{
		Parent: d.Parent,
		Author: UserVO{
			Account: d.Author.Account,
			Email:   d.Author.Email,
		},
		Committer: UserVO{
			Account: d.Committer.Account,
			Email:   d.Committer.Email,
		},
		AuthoredTime:  time.UnixMilli(d.AuthoredTime).Format(time.DateTime),
		CommittedTime: time.UnixMilli(d.CommittedTime).Format(time.DateTime),
		CommitMsg:     d.CommitMsg,
		CommitId:      d.CommitId,
		ShortId:       d.ShortId,
		Verified:      d.Verified,
	}
	if d.TaggerTime > 0 {
		t := time.UnixMilli(d.TaggerTime).Format(time.DateTime)
		ret.TaggerTime = &t
		ret.Tagger = &UserVO{
			Account: d.Tagger.Account,
			Email:   d.Tagger.Email,
		}
		ret.ShortTagId = &d.ShortTagId
		ret.TagCommitMsg = &d.TagCommitMsg
	}
	if d.Signer.Account != "" {
		ret.Signer = &SignerVO{
			Account:   d.Signer.Account,
			AvatarUrl: d.Signer.AvatarUrl,
			Name:      d.Signer.Name,
			Key:       d.Signer.Key,
			Type:      d.Signer.Type,
		}
	}
	return ret
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
	err := reposrv.Outer.DeleteRepo(c, reposrv.DeleteRepoReqDTO{
		Operator: apisession.MustGetLoginUser(c),
		RepoId:   getRepoId(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func deletePermanentlyRepo(c *gin.Context) {
	err := reposrv.Outer.DeleteRepoPermanently(c, reposrv.DeleteRepoReqDTO{
		Operator: apisession.MustGetLoginUser(c),
		RepoId:   getRepoId(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func recoverFromRecycle(c *gin.Context) {
	err := reposrv.Outer.RecoverFromRecycle(c, reposrv.RecoverFromRecycleReqDTO{
		Operator: apisession.MustGetLoginUser(c),
		RepoId:   getRepoId(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	util.DefaultOkResponse(c)
}

func listRepo(c *gin.Context) {
	repoList, err := reposrv.Outer.ListRepo(c, reposrv.ListRepoReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(repoList, func(t reposrv.RepoDTO) (RepoVO, error) {
		return repo2VO(t), nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]RepoVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func listDeletedRepo(c *gin.Context) {
	repoList, err := reposrv.Outer.ListDeletedRepo(c, reposrv.ListDeletedRepoReqDTO{
		TeamId:   cast.ToInt64(c.Param("teamId")),
		Operator: apisession.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	data, _ := listutil.Map(repoList, func(t reposrv.DeletedRepoDTO) (DeletedRepoVO, error) {
		return DeletedRepoVO{
			RepoVO:  repo2VO(t.RepoDTO),
			Deleted: t.Deleted.Format(time.DateTime),
		}, nil
	})
	c.JSON(http.StatusOK, ginutil.DataResp[[]DeletedRepoVO]{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
	})
}

func repo2VO(t reposrv.RepoDTO) RepoVO {
	return RepoVO{
		RepoId:       t.Id,
		Name:         t.Name,
		Path:         t.Path,
		RepoDesc:     t.RepoDesc,
		GitSize:      t.GitSize,
		LfsSize:      t.LfsSize,
		Created:      t.Created.Format(time.DateTime),
		TeamId:       t.TeamId,
		LastOperated: t.LastOperated.Format(time.DateTime),
		DisableLfs:   t.DisableLfs,
		LfsLimitSize: t.LfsLimitSize,
		GitLimitSize: t.GitLimitSize,
		IsArchived:   t.IsArchived,
	}
}

func catFile(c *gin.Context) {
	var req CatFileReqVO
	if util.ShouldBindQuery(&req, c) {
		resp, err := reposrv.Outer.CatFile(c, reposrv.CatFileReqDTO{
			RepoId:   req.RepoId,
			Ref:      req.Ref,
			FilePath: req.FilePath,
			RefType:  req.RefType,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DataResp[CatFileVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data: CatFileVO{
				FileMode: resp.ModeName,
				Content:  resp.Content,
				Size:     util.VolumeReadable(resp.Size),
				Commit:   commitDto2Vo(resp.Commit),
			},
		})
	}
}

func blame(c *gin.Context) {
	var req BlameReqVO
	if util.ShouldBindQuery(&req, c) {
		lines, err := reposrv.Outer.Blame(c, reposrv.BlameReqDTO{
			RepoId:   req.RepoId,
			Ref:      req.Ref,
			FilePath: req.FilePath,
			RefType:  req.RefType,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		data, _ := listutil.Map(lines, func(t reposrv.BlameLineDTO) (BlameLineVO, error) {
			return BlameLineVO{
				Number: t.Number,
				Commit: commitDto2Vo(t.Commit),
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[[]BlameLineVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     data,
		})
	}
}

func diffFile(c *gin.Context) {
	var req DiffFileReqVO
	if util.ShouldBindQuery(&req, c) {
		respDTO, err := reposrv.Outer.DiffFile(c, reposrv.DiffFileReqDTO{
			RepoId:   req.RepoId,
			Target:   req.Target,
			Head:     req.Head,
			FilePath: req.FilePath,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := DiffFileVO{
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
				LeftNo:  t.LeftNo,
				Prefix:  t.Prefix,
				RightNo: t.RightNo,
				Text:    t.Text,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[DiffFileVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}

func createArchive(c *gin.Context) {
	var req CreateArchiveReqVO
	if util.ShouldBindQuery(&req, c) {
		err := reposrv.Outer.CreateArchive(c, reposrv.CreateArchiveReqDTO{
			RepoId:   req.RepoId,
			FileName: req.FileName,
			C:        c,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
	}
}

func diffRefs(c *gin.Context) {
	var req DiffRefsReqVO
	if util.ShouldBindQuery(&req, c) {
		respDTO, err := reposrv.Outer.DiffRefs(c, reposrv.DiffRefsReqDTO{
			RepoId:     req.RepoId,
			Target:     req.Target,
			TargetType: req.TargetType,
			Head:       req.Head,
			HeadType:   req.HeadType,
			Operator:   apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		respVO := DiffRefsVO{
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
				InsertNums: t.InsertNums,
				DeleteNums: t.DeleteNums,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[DiffRefsVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     respVO,
		})
	}
}

func diffCommits(c *gin.Context) {
	var req DiffCommitsReqVO
	if util.ShouldBindQuery(&req, c) {
		respDTO, err := reposrv.Outer.DiffCommits(c, reposrv.DiffCommitsReqDTO{
			RepoId:   req.RepoId,
			CommitId: req.CommitId,
			Operator: apisession.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		respVO := DiffCommitsVO{
			Commit:   commitDto2Vo(respDTO.Commit),
			NumFiles: respDTO.NumFiles,
			DiffNumsStats: DiffNumsStatInfoVO{
				FileChangeNums: respDTO.DiffNumsStats.FileChangeNums,
				InsertNums:     respDTO.DiffNumsStats.InsertNums,
				DeleteNums:     respDTO.DiffNumsStats.DeleteNums,
			},
		}
		respVO.DiffNumsStats.Stats, _ = listutil.Map(respDTO.DiffNumsStats.Stats, func(t reposrv.DiffNumsStatDTO) (DiffNumsStatVO, error) {
			return DiffNumsStatVO{
				RawPath:    t.RawPath,
				Path:       t.Path,
				InsertNums: t.InsertNums,
				DeleteNums: t.DeleteNums,
			}, nil
		})
		c.JSON(http.StatusOK, ginutil.DataResp[DiffCommitsVO]{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     respVO,
		})
	}
}

func historyCommits(c *gin.Context) {
	var req HistoryCommitsReqVO
	if util.ShouldBindQuery(&req, c) {
		respDTO, err := reposrv.Outer.HistoryCommits(c, reposrv.HistoryCommitsReqDTO{
			RepoId:   req.RepoId,
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

func transferTeam(c *gin.Context) {
	var req TransferTeamReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Outer.TransferTeam(c, reposrv.TransferTeamReqDTO{
			RepoId:   req.RepoId,
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

func getRepoId(c *gin.Context) int64 {
	return cast.ToInt64(c.Param("repoId"))
}
