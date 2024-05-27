package store

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/logger"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	protocolPattern = regexp.MustCompile(`^[0-9a-zA-Z]+=[0-9a-zA-Z]+(:[0-9a-zA-Z]+=[0-9a-zA-Z]+)*$`)
)

const refLimit = 10

type storeImpl struct{}

func NewStore() Store {
	return new(storeImpl)
}

// InitRepo 初始化仓库
func (s *storeImpl) InitRepo(ctx context.Context, req reqvo.InitRepoReq) (int64, error) {
	err := git.InitRepository(ctx, git.InitRepoOpts{
		Owner: git.User{
			Account: req.UserAccount,
			Email:   req.UserEmail,
		},
		RepoName:      req.RepoName,
		RepoPath:      filepath.Join(git.RepoDir(), req.RepoPath),
		AddReadme:     req.AddReadme,
		GitIgnoreName: req.GitIgnoreName,
		DefaultBranch: req.DefaultBranch,
	})
	// 仓库已存在
	if err == git.RepoExistsErr || err == nil {
		gitSize, _ := getGitSize(req.RepoPath)
		return gitSize, nil
	}
	logger.Logger.WithContext(ctx).Error(err)
	return 0, util.InternalError(err)
}

// DeleteRepo 删除仓库
func (s *storeImpl) DeleteRepo(ctx context.Context, req reqvo.DeleteRepoReq) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	util.RemoveAll(repoPath)
	lfsPath := filepath.Join(git.LfsDir(), req.RepoPath)
	util.RemoveAll(lfsPath)
	return nil
}

// GetAllBranches 获取所有分支
func (s *storeImpl) GetAllBranches(ctx context.Context, req reqvo.GetAllBranchesReq) ([]reqvo.RefVO, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ret, err := git.GetAllBranchList(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(ret, func(t git.Ref) (reqvo.RefVO, error) {
		return reqvo.RefVO{
			LastCommitId: t.Sha,
			Name:         t.Name,
		}, nil
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PageBranchAndLastCommit 分页获取分支+最后提交信息
func (s *storeImpl) PageBranchAndLastCommit(ctx context.Context, req reqvo.PageRefCommitsReq) ([]reqvo.RefCommitVO, int64, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ret, err := git.GetAllBranchList(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	offset := (req.PageNum - 1) * refLimit
	totalCount := len(ret)
	if offset >= totalCount {
		return []reqvo.RefCommitVO{}, int64(totalCount), nil
	}
	ret = ret[offset:min(offset+refLimit, totalCount)]
	data, err := listutil.Map(ret, func(t git.Ref) (reqvo.RefCommitVO, error) {
		commit, err := git.GetCommitByCommitId(ctx, repoPath, t.Sha)
		if err != nil {
			return reqvo.RefCommitVO{}, err
		}
		return reqvo.RefCommitVO{
			Commit: commit2Vo(commit),
			Name:   t.Name,
		}, nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	return data, int64(totalCount), nil
}

// PageTagAndCommit 分页获取tag+提交信息
func (s *storeImpl) PageTagAndCommit(ctx context.Context, req reqvo.PageRefCommitsReq) ([]reqvo.RefCommitVO, int64, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ret, err := git.GetAllTagList(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	offset := (req.PageNum - 1) * refLimit
	totalCount := len(ret)
	if offset >= totalCount {
		return []reqvo.RefCommitVO{}, int64(totalCount), nil
	}
	ret = ret[offset:min(offset+refLimit, totalCount)]
	data, err := listutil.Map(ret, func(t git.Ref) (reqvo.RefCommitVO, error) {
		commit, err := git.GetCommitByTag(ctx, repoPath, t.Name)
		if err != nil {
			return reqvo.RefCommitVO{}, err
		}
		return reqvo.RefCommitVO{
			Commit: commit2Vo(commit),
			Name:   t.Name,
		}, nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	return data, int64(totalCount), nil
}

// DeleteBranch 删除分支
func (s *storeImpl) DeleteBranch(ctx context.Context, req reqvo.DeleteBranchReq) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	err := git.DeleteBranch(ctx, repoPath, req.Branch, true)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// GetAllTags 获取所有tag
func (s *storeImpl) GetAllTags(ctx context.Context, req reqvo.GetAllTagsReq) ([]reqvo.RefVO, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ret, err := git.GetAllTagList(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(ret, func(t git.Ref) (reqvo.RefVO, error) {
		return reqvo.RefVO{
			LastCommitId: t.Sha,
			Name:         t.Name,
		}, nil
	})
}

// Gc 触发仓库gc
func (s *storeImpl) Gc(ctx context.Context, req reqvo.GcReq) (int64, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	err := git.Gc(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return 0, util.InternalError(err)
	}
	gitSize, err := getGitSize(req.RepoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return 0, util.InternalError(err)
	}
	return gitSize, nil
}

// DiffRefs 对比两个ref差异
func (s *storeImpl) DiffRefs(ctx context.Context, req reqvo.DiffRefsReq) (reqvo.DiffRefsResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	head := req.HeadType.PackRef(req.Head)
	target := req.TargetType.PackRef(req.Target)
	if !git.CheckExists(ctx, repoPath, head) {
		return reqvo.DiffRefsResp{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, repoPath, target) {
		return reqvo.DiffRefsResp{}, util.InvalidArgsError()
	}
	info, err := git.GetDiffRefsInfo(ctx, repoPath, req.Target, req.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.DiffRefsResp{}, util.InternalError(err)
	}
	return diffResToResp(info), nil
}

func diffResToResp(info git.DiffRefsInfo) reqvo.DiffRefsResp {
	ret := reqvo.DiffRefsResp{
		Target:       info.Target,
		Head:         info.Head,
		TargetCommit: commit2Vo(info.TargetCommit),
		HeadCommit:   commit2Vo(info.HeadCommit),
		NumFiles:     info.NumFiles,
		DiffNumsStats: reqvo.DiffNumsStatInfoVO{
			FileChangeNums: info.DiffNumsStats.FileChangeNums,
			InsertNums:     info.DiffNumsStats.InsertNums,
			DeleteNums:     info.DiffNumsStats.DeleteNums,
		},
		ConflictFiles: info.ConflictFiles,
	}
	ret.DiffNumsStats.Stats, _ = listutil.Map(info.DiffNumsStats.Stats, func(t git.DiffNumsStat) (reqvo.DiffNumsStatVO, error) {
		return reqvo.DiffNumsStatVO{
			RawPath:    t.Path,
			Path:       path.Base(t.Path),
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	ret.Commits, _ = listutil.Map(info.Commits, func(t git.Commit) (reqvo.CommitVO, error) {
		return commit2Vo(t), nil
	})
	ret.CanMerge = info.IsMergeAble()
	return ret
}

// DiffCommits 对比两个提交差异
func (s *storeImpl) DiffCommits(ctx context.Context, req reqvo.DiffCommitsReq) (reqvo.DiffCommitsResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if !git.CheckExists(ctx, repoPath, req.CommitId) {
		return reqvo.DiffCommitsResp{}, util.InvalidArgsError()
	}
	info, err := git.GetDiffCommitsInfo(ctx, repoPath, req.CommitId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.DiffCommitsResp{}, util.InternalError(err)
	}
	ret := reqvo.DiffCommitsResp{
		Commit:   commit2Vo(info.Commit),
		NumFiles: info.NumFiles,
		DiffNumsStats: reqvo.DiffNumsStatInfoVO{
			FileChangeNums: info.DiffNumsStats.FileChangeNums,
			InsertNums:     info.DiffNumsStats.InsertNums,
			DeleteNums:     info.DiffNumsStats.DeleteNums,
		},
	}
	ret.DiffNumsStats.Stats, _ = listutil.Map(info.DiffNumsStats.Stats, func(t git.DiffNumsStat) (reqvo.DiffNumsStatVO, error) {
		return reqvo.DiffNumsStatVO{
			RawPath:    t.Path,
			Path:       path.Base(t.Path),
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	return ret, nil
}

// CanMerge 两个ref是否可以合并
func (s *storeImpl) CanMerge(ctx context.Context, req reqvo.CanMergeReq) (bool, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	head := req.HeadType.PackRef(req.Head)
	target := req.TargetType.PackRef(req.Target)
	if !git.CheckExists(ctx, repoPath, head) {
		return false, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, repoPath, target) {
		return false, util.InvalidArgsError()
	}
	ret, err := git.CanMerge(ctx, repoPath, target, head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return false, util.InternalError(err)
	}
	return ret, nil
}

// DiffFile 对比两个分支单个文件差异
func (s *storeImpl) DiffFile(ctx context.Context, req reqvo.DiffFileReq) (reqvo.DiffFileResp, error) {
	if req.Target == "" {
		return reqvo.DiffFileResp{}, util.InvalidArgsError()
	}
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if req.Head == "" {
		req.Head = git.EmptyTreeSHA
	}
	d, err := git.GetDiffFileDetail(ctx, repoPath, req.Target, req.Head, req.FilePath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.DiffFileResp{}, err
	}
	ret := reqvo.DiffFileResp{
		FilePath:    d.FilePath,
		OldMode:     d.OldMode,
		Mode:        d.Mode,
		IsSubModule: d.IsSubModule,
		FileType:    d.FileType.String(),
		IsBinary:    d.IsBinary,
		RenameFrom:  d.RenameFrom,
		RenameTo:    d.RenameTo,
		CopyFrom:    d.CopyFrom,
		CopyTo:      d.CopyTo,
	}
	ret.Lines, _ = listutil.Map(d.Lines, func(t git.DiffLine) (reqvo.DiffLineVO, error) {
		return reqvo.DiffLineVO{
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}, nil
	})
	return ret, nil
}

// GetRepoSize 获取仓库大小
func (s *storeImpl) GetRepoSize(ctx context.Context, req reqvo.GetRepoSizeReq) (int64, int64, error) {
	return getGitAndLfsSize(ctx, req.RepoPath)
}

func getGitAndLfsSize(ctx context.Context, repoPath string) (int64, int64, error) {
	gitSize, err := getGitSize(repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return 0, 0, util.InternalError(err)
	}
	lfsSize, err := git.GetDirSize(filepath.Join(git.LfsDir(), repoPath))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return 0, 0, util.InternalError(err)
	}
	return gitSize, lfsSize, nil
}

func getGitSize(repoPath string) (int64, error) {
	return git.GetDirSize(filepath.Join(git.RepoDir(), repoPath))
}

// ShowDiffTextContent 获取某个commitId文件内容
func (s *storeImpl) ShowDiffTextContent(ctx context.Context, req reqvo.ShowDiffTextContentReq) ([]reqvo.DiffLineVO, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if !git.CheckRefIsCommit(ctx, repoPath, req.CommitId) {
		return nil, util.InvalidArgsError()
	}
	lineList, err := git.ShowFileTextContentByCommitId(ctx, repoPath, req.CommitId, req.FileName, req.StartLine, req.Limit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	ret := make([]reqvo.DiffLineVO, 0, len(lineList))
	for i, line := range lineList {
		n := req.StartLine + i
		ret = append(ret, reqvo.DiffLineVO{
			LeftNo:  n,
			Prefix:  " ",
			RightNo: n,
			Text:    line,
		})
	}
	return ret, nil
}

// HistoryCommits 获取历史提交记录
func (s *storeImpl) HistoryCommits(ctx context.Context, req reqvo.HistoryCommitsReq) (reqvo.HistoryCommitsResp, error) {
	limit := 10
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if !git.CheckExists(ctx, repoPath, req.Ref) {
		return reqvo.HistoryCommitsResp{}, util.InvalidArgsError()
	}
	commitList, err := git.GetGitLogCommitList(ctx, repoPath, req.Ref, req.Offset, limit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.HistoryCommitsResp{}, util.InternalError(err)
	}
	ret := reqvo.HistoryCommitsResp{
		Cursor: req.Offset + limit,
	}
	ret.Data, _ = listutil.Map(commitList, func(t git.Commit) (reqvo.CommitVO, error) {
		return commit2Vo(t), nil
	})
	return ret, nil
}

// InitRepoHook 重建仓库hook
func (s *storeImpl) InitRepoHook(ctx context.Context, req reqvo.InitRepoHookReq) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	err := git.InitRepoHook(repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// EntriesRepo 仓库文件列表
func (s *storeImpl) EntriesRepo(ctx context.Context, req reqvo.EntriesRepoReq) ([]reqvo.BlobVO, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if req.Dir == "" {
		req.Dir = "."
	}
	ref := req.RefType.PackRef(req.Ref)
	blobs, err := git.LsTreeBlob(ctx, repoPath, ref, req.Dir)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(blobs, func(t git.LsTreeRet) (reqvo.BlobVO, error) {
		return reqvo.BlobVO{
			Mode:    t.Mode.Readable(),
			RawPath: t.Path,
			Path:    path.Base(t.Path),
		}, nil
	})
}

// CatFile 展示文件内容
func (s *storeImpl) CatFile(ctx context.Context, req reqvo.CatFileReq) (reqvo.CatFileResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ref := req.RefType.PackRef(req.Ref)
	commit, err := git.GetFileLatestCommit(ctx, repoPath, ref, req.FilePath)
	if err != nil {
		return reqvo.CatFileResp{}, err
	}
	fileMode, content, size, _, err := git.GetFileTextContentByRef(ctx, repoPath, ref, req.FilePath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.CatFileResp{}, util.InternalError(err)
	}
	return reqvo.CatFileResp{
		FileMode: fileMode.String(),
		ModeName: fileMode.Readable(),
		Content:  content,
		Size:     size,
		Commit:   commit2Vo(commit),
	}, nil
}

// IndexRepo 仓库首页
func (s *storeImpl) IndexRepo(ctx context.Context, req reqvo.IndexRepoReq) (reqvo.IndexRepoResp, error) {
	ref := req.RefType.PackRef(req.Ref)
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	dir := "."
	latestCommit, err := git.GetFileLatestCommit(ctx, repoPath, ref, dir)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.IndexRepoResp{}, util.InternalError(err)
	}
	commits, err := git.LsTreeCommit(ctx, repoPath, ref, dir)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.IndexRepoResp{}, util.InternalError(err)
	}
	_, readme, _, hasReadme, err := git.GetFileTextContentByRef(ctx, repoPath, ref, filepath.Join(dir, "README.md"))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	if err == nil && !hasReadme {
		_, readme, _, hasReadme, err = git.GetFileTextContentByRef(ctx, repoPath, ref, filepath.Join(dir, "readme.md"))
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	}
	return reqvo.IndexRepoResp{
		ReadmeText:   readme,
		HasReadme:    hasReadme,
		LatestCommit: commit2Vo(latestCommit),
		Tree:         lsRet2TreeVO(commits),
	}, nil
}

// UploadPack git-upload-pack
func (s *storeImpl) UploadPack(req reqvo.UploadPackReq) {
	// 校验content-type
	if req.C.GetHeader("Content-Type") != "application/x-git-upload-pack-request" {
		req.C.String(http.StatusForbidden, "bad content type")
		return
	}
	reqBody := req.C.Request.Body
	var err error
	if req.C.GetHeader("Content-Encoding") == "gzip" {
		reqBody, err = gzip.NewReader(reqBody)
		if err != nil {
			logger.Logger.WithContext(req.C).Error(err)
			req.C.String(http.StatusInternalServerError, "internal error")
			return
		}
	}
	defer reqBody.Close()
	// 不缓存任何东西
	req.C.Writer.WriteHeader(http.StatusOK)
	req.C.Header("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	req.C.Header("Pragma", "no-cache")
	req.C.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
	req.C.Header("Content-Type", "application/x-git-upload-pack-result")
	req.C.Writer.WriteHeaderNow()
	env := make([]string, 0)
	protocol := req.C.GetHeader("Git-Protocol")
	if protocolPattern.MatchString(protocol) {
		env = append(env, "GIT_PROTOCOL="+protocol)
	}
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	env = append(env, util.JoinFields(
		gitenv.EnvRepoId, req.C.GetHeader("Repo-PrId"),
		gitenv.EnvPusherAccount, req.C.GetHeader("Pusher-Account"),
		gitenv.EnvPusherEmail, req.C.GetHeader("Pusher-Email"),
		gitenv.EnvAppUrl, req.C.GetHeader("App-Url"),
		gitenv.EnvHookToken, git.HookToken(),
	)...)
	err = git.UploadPack(req.C, repoPath, reqBody, req.C.Writer, env)
	if err != nil {
		logger.Logger.WithContext(req.C).Error(err)
		req.C.String(http.StatusInternalServerError, "internal error")
	}
}

// ReceivePack git-receive-pack
func (s *storeImpl) ReceivePack(req reqvo.ReceivePackReq) {
	// 校验content-type
	if req.C.GetHeader("Content-Type") != "application/x-git-receive-pack-request" {
		req.C.String(http.StatusForbidden, "bad content type")
		return
	}
	reqBody := req.C.Request.Body
	var err error
	if req.C.GetHeader("Content-Encoding") == "gzip" {
		reqBody, err = gzip.NewReader(reqBody)
		if err != nil {
			req.C.String(http.StatusInternalServerError, "internal error")
			return
		}
	}
	defer reqBody.Close()
	// 不缓存任何东西
	req.C.Writer.WriteHeader(http.StatusOK)
	req.C.Header("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	req.C.Header("Pragma", "no-cache")
	req.C.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
	req.C.Header("Content-Type", "application/x-git-receive-pack-result")
	req.C.Writer.WriteHeaderNow()
	env := make([]string, 0)
	protocol := req.C.GetHeader("Git-Protocol")
	if protocolPattern.MatchString(protocol) {
		env = append(env, "GIT_PROTOCOL="+protocol)
	}
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	env = append(env, util.JoinFields(
		gitenv.EnvHookUrl, fmt.Sprintf("http://127.0.0.1:%d", common.HttpServerPort()),
		gitenv.EnvRepoId, req.C.GetHeader("Repo-PrId"),
		gitenv.EnvPusherAccount, req.C.GetHeader("Pusher-Account"),
		gitenv.EnvPusherEmail, req.C.GetHeader("Pusher-Email"),
		gitenv.EnvAppUrl, req.C.GetHeader("App-Url"),
		gitenv.EnvHookToken, git.HookToken(),
		"--stateless-rpc", repoPath,
	)...)
	err = git.ReceivePack(req.C, repoPath, reqBody, req.C.Writer, env)
	if err != nil {
		logger.Logger.WithContext(req.C).Error(err)
		req.C.String(http.StatusInternalServerError, "internal error")
	}
}

// InfoRefs smart http infoRefs
func (s *storeImpl) InfoRefs(ctx context.Context, req reqvo.InfoRefsReq) {
	serviceParam := req.C.Query("service")
	// 不缓存任何东西
	req.C.Header("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	req.C.Header("Pragma", "no-cache")
	req.C.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
	req.C.Header("Content-Type", fmt.Sprintf("application/x-%s-advertisement", serviceParam))
	req.C.Writer.WriteHeaderNow()
	env := make([]string, 0)
	protocol := req.C.GetHeader("Git-Protocol")
	if protocolPattern.MatchString(protocol) {
		env = append(env, "GIT_PROTOCOL="+protocol)
	}
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	refs, err := git.InfoRefs(ctx, repoPath, strings.TrimPrefix(serviceParam, "git-"), env)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		req.C.String(http.StatusInternalServerError, "internal error")
		return
	}
	_, _ = req.C.Writer.Write(packetWrite("# service=" + serviceParam + "\n"))
	_, _ = req.C.Writer.Write([]byte("0000"))
	_, _ = req.C.Writer.Write(refs)
}

// Merge 合并两个分支
func (s *storeImpl) Merge(ctx context.Context, req reqvo.MergeReq) (reqvo.DiffRefsResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	info, err := git.Merge(ctx, repoPath, req.Target, req.Head, git.MergeRepoOpts{
		RepoId:        req.MergeOpts.RepoId,
		PrId:          req.MergeOpts.PrId,
		PusherAccount: req.MergeOpts.PusherAccount,
		Message:       req.MergeOpts.Message,
		AppUrl:        req.MergeOpts.AppUrl,
	})
	if err != nil {
		if strings.Contains(err.Error(), "can not merge") {
			return diffResToResp(info), util.NewBizErr(apicode.PullRequestCannotMergeCode, i18n.PullRequestCannotMerge)
		}
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.DiffRefsResp{}, util.InternalError(err)
	}
	return diffResToResp(info), nil
}

// Blame git blame获取每一行提交人和时间
func (s *storeImpl) Blame(ctx context.Context, req reqvo.BlameReq) ([]reqvo.BlameLineVO, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ref := req.RefType.PackRef(req.Ref)
	lines, err := git.Blame(ctx, repoPath, ref, req.FilePath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(lines, func(t git.BlameLine) (reqvo.BlameLineVO, error) {
		return reqvo.BlameLineVO{
			Number: t.Number,
			Commit: commit2Vo(*t.Commit),
		}, nil
	})
}

// DeleteTag 删除tag
func (s *storeImpl) DeleteTag(ctx context.Context, req reqvo.DeleteTagReqVO) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if !git.CheckRefIsTag(ctx, repoPath, req.Tag) {
		return util.InvalidArgsError()
	}
	err := git.DeleteTag(ctx, repoPath, req.Tag)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// CreateArchive 下载压缩包
func (s *storeImpl) CreateArchive(ctx context.Context, req reqvo.CreateArchiveReq) {
	var (
		archiveType git.ArchiveType
		ref         string
	)
	if strings.HasSuffix(req.FileName, ".zip") {
		archiveType = git.ZIP
		ref = strings.TrimSuffix(req.FileName, ".zip")
	} else if strings.HasSuffix(req.FileName, ".tar.gz") {
		archiveType = git.TARGZ
		ref = strings.TrimSuffix(req.FileName, ".tar.gz")
	} else {
		util.HandleApiErr(util.InvalidArgsError(), req.C)
		return
	}
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	var (
		commitId string
		err      error
	)
	if git.CheckRefIsTag(ctx, repoPath, ref) {
		var commit git.Commit
		commit, err = git.GetCommitByTag(ctx, repoPath, ref)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			util.HandleApiErr(util.InternalError(err), req.C)
			return
		}
		commitId = commit.Id
	} else if git.CheckRefIsBranch(ctx, repoPath, ref) {
		commitId, err = git.GetBranchCommitId(ctx, repoPath, ref)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			util.HandleApiErr(util.InternalError(err), req.C)
			return
		}
	} else {
		util.HandleApiErr(util.InvalidArgsError(), req.C)
		return
	}
	// 不缓存任何东西
	req.C.Writer.WriteHeader(http.StatusOK)
	req.C.Header("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	req.C.Header("Pragma", "no-cache")
	req.C.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
	req.C.Header("Content-Type", archiveType.HttpContentType())
	req.C.Header("Content-Disposition", "attachment; filename=\""+req.FileName+"\"")
	req.C.Header("Access-Control-Expose-Headers", "Content-Disposition")
	req.C.Writer.WriteHeaderNow()
	// 暂时不搞http缓存
	err = git.CreateArchive(ctx, repoPath, commitId, archiveType, req.C.Writer)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		util.HandleApiErr(util.InternalError(err), req.C)
		return
	}
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

func lsRet2TreeVO(commits []git.FileCommit) reqvo.TreeVO {
	files, _ := listutil.Map(commits, func(t git.FileCommit) (reqvo.FileVO, error) {
		ret := reqvo.FileVO{
			Mode:    t.Mode.Readable(),
			RawPath: t.Path,
			Path:    path.Base(t.Path),
			Commit: reqvo.CommitVO{
				Author: reqvo.UserVO{
					Account: t.Author.Account,
					Email:   t.Author.Email,
				},
				Committer: reqvo.UserVO{
					Account: t.Committer.Account,
					Email:   t.Committer.Email,
				},
				AuthoredTime:  t.AuthorSigTime.UnixMilli(),
				CommittedTime: t.CommitSigTime.UnixMilli(),
				CommitMsg:     t.CommitMsg,
				CommitId:      t.Blob,
				ShortId:       util.LongCommitId2ShortId(t.Blob),
			},
		}
		return ret, nil
	})
	return reqvo.TreeVO{
		Files: files,
	}
}

func commit2Vo(c git.Commit) reqvo.CommitVO {
	ret := reqvo.CommitVO{
		Parent: c.Parent,
		Author: reqvo.UserVO{
			Account: c.Author.Account,
			Email:   c.Author.Email,
		},
		Committer: reqvo.UserVO{
			Account: c.Committer.Account,
			Email:   c.Committer.Email,
		},
		AuthoredTime:  c.AuthorSigTime.UnixMilli(),
		CommittedTime: c.CommitSigTime.UnixMilli(),
		CommitMsg:     c.CommitMsg,
		CommitId:      c.Id,
		ShortId:       util.LongCommitId2ShortId(c.Id),
		CommitSig:     c.CommitSig.String(),
		Payload:       c.Payload,
	}
	if c.Tag != nil {
		ret.Tagger = reqvo.UserVO{
			Account: c.Tag.Tagger.Account,
			Email:   c.Tag.Tagger.Email,
		}
		ret.TaggerTime = c.Tag.TagTime.UnixMilli()
		ret.ShortTagId = util.LongCommitId2ShortId(c.Tag.Id)
		ret.TagCommitMsg = c.Tag.CommitMsg
	}
	return ret
}
