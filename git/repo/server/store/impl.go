package store

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/gitenv"
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

const (
	LsTreeCommitLimit = 25
)

var (
	protocolPattern = regexp.MustCompile(`^[0-9a-zA-Z]+=[0-9a-zA-Z]+(:[0-9a-zA-Z]+=[0-9a-zA-Z]+)*$`)
)

type storeImpl struct{}

func NewStore() Store {
	return new(storeImpl)
}

// InitRepo 初始化仓库
func (s *storeImpl) InitRepo(ctx context.Context, req reqvo.InitRepoReq) error {
	err := git.InitRepository(ctx, git.InitRepoOpts{
		Owner: git.User{
			Account: req.UserAccount,
			Email:   req.UserEmail,
		},
		RepoName:      req.RepoName,
		RepoPath:      filepath.Join(git.RepoDir(), req.RepoPath),
		CreateReadme:  req.CreateReadme,
		GitIgnoreName: req.GitIgnoreName,
		DefaultBranch: req.DefaultBranch,
	})
	// 仓库已存在
	if err == git.RepoExistsErr {
		return nil
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteRepo 删除仓库
func (s *storeImpl) DeleteRepo(ctx context.Context, req reqvo.DeleteRepoReq) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	util.RemoveAll(repoPath)
	return nil
}

// GetAllBranches 获取所有分支
func (s *storeImpl) GetAllBranches(ctx context.Context, req reqvo.GetAllBranchesReq) ([]string, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ret, err := git.GetAllBranchList(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return ret, nil
}

// GetAllTags 获取所有tag
func (s *storeImpl) GetAllTags(ctx context.Context, req reqvo.GetAllTagsReq) ([]string, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	ret, err := git.GetAllTagList(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return ret, nil
}

// Gc 触发仓库gc
func (s *storeImpl) Gc(ctx context.Context, req reqvo.GcReq) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	err := git.Gc(ctx, repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DiffRefs 对比两个ref差异
func (s *storeImpl) DiffRefs(ctx context.Context, req reqvo.DiffRefsReq) (reqvo.DiffRefsResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if !git.CheckExists(ctx, repoPath, req.Head) {
		return reqvo.DiffRefsResp{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, repoPath, req.Target) {
		return reqvo.DiffRefsResp{}, util.InvalidArgsError()
	}
	// 读取ssh key和gpg key
	info, err := git.GetDiffRefsInfo(ctx, repoPath, req.Target, req.Head)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.DiffRefsResp{}, util.InternalError(err)
	}
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
			TotalNums:  t.TotalNums,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	ret.Commits, _ = listutil.Map(info.Commits, func(t git.Commit) (reqvo.CommitVO, error) {
		return commit2Vo(t), nil
	})
	ret.CanMerge = len(ret.Commits) > 0 && len(ret.ConflictFiles) == 0
	return ret, nil
}

// DiffFile 对比两个分支单个文件差异
func (s *storeImpl) DiffFile(ctx context.Context, req reqvo.DiffFileReq) (reqvo.DiffFileResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if !git.CheckExists(ctx, repoPath, req.Target) {
		return reqvo.DiffFileResp{}, util.InvalidArgsError()
	}
	if !git.CheckExists(ctx, repoPath, req.Head) {
		return reqvo.DiffFileResp{}, util.InvalidArgsError()
	}
	d, err := git.GetDiffFileDetail(ctx, repoPath, req.Target, req.Head, req.FileName)
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
			Index:   t.Index,
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}, nil
	})
	return ret, nil
}

// GetRepoSize 获取仓库大小
func (s *storeImpl) GetRepoSize(ctx context.Context, req reqvo.GetRepoSizeReq) (int64, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	size, err := git.GetRepoSize(repoPath)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return 0, util.InternalError(err)
	}
	return size, nil
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
			Index:   i,
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
func (s *storeImpl) EntriesRepo(ctx context.Context, req reqvo.EntriesRepoReq) (reqvo.TreeVO, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	if req.Dir == "" {
		req.Dir = "."
	}
	commits, err := git.LsTreeCommit(ctx, repoPath, req.Ref, req.Dir, req.Offset, LsTreeCommitLimit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.TreeVO{}, util.InternalError(err)
	}
	return lsRet2TreeDTO(commits, req.Offset, LsTreeCommitLimit), nil
}

// CatFile 展示文件内容
func (s *storeImpl) CatFile(ctx context.Context, req reqvo.CatFileReq) (reqvo.CatFileResp, error) {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	fileMode, content, _, err := git.GetFileContentByRef(ctx, repoPath, req.Ref, req.FileName)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.CatFileResp{}, util.InternalError(err)
	}
	return reqvo.CatFileResp{
		FileMode: fileMode.String(),
		ModeName: fileMode.Readable(),
		Content:  content,
	}, nil
}

// TreeRepo 仓库首页
func (s *storeImpl) TreeRepo(ctx context.Context, req reqvo.TreeRepoReq) (reqvo.TreeRepoResp, error) {
	if req.Dir == "" {
		req.Dir = "."
	}
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	latestCommit, err := git.GetFileLatestCommit(ctx, repoPath, req.Ref, req.Dir)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.TreeRepoResp{}, util.InternalError(err)
	}
	commits, err := git.LsTreeCommit(ctx, repoPath, req.Ref, req.Dir, 0, LsTreeCommitLimit)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return reqvo.TreeRepoResp{}, util.InternalError(err)
	}
	_, readme, hasReadme, err := git.GetFileContentByRef(ctx, repoPath, req.Ref, filepath.Join(req.Dir, "readme.md"))
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	if err == nil && !hasReadme {
		_, readme, hasReadme, err = git.GetFileContentByRef(ctx, repoPath, req.Ref, filepath.Join(req.Dir, "README.md"))
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	}
	return reqvo.TreeRepoResp{
		ReadmeText:   readme,
		HasReadme:    hasReadme,
		LatestCommit: commit2Vo(latestCommit),
		Tree:         lsRet2TreeDTO(commits, 0, LsTreeCommitLimit),
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
		gitenv.EnvRepoId, req.C.GetHeader("Repo-Id"),
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
		gitenv.EnvRepoId, req.C.GetHeader("Repo-Id"),
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
func (s *storeImpl) Merge(ctx context.Context, req reqvo.MergeReq) error {
	repoPath := filepath.Join(git.RepoDir(), req.RepoPath)
	err := git.Merge(ctx, repoPath, req.Target, req.Head, git.MergeRepoOpts{
		RepoId:        req.MergeOpts.RepoId,
		PrId:          req.MergeOpts.PrId,
		PusherAccount: req.MergeOpts.PusherAccount,
		Message:       req.MergeOpts.Message,
		AppUrl:        req.MergeOpts.AppUrl,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

func lsRet2TreeDTO(commits []git.FileCommit, offset, limit int) reqvo.TreeVO {
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
		Files:   files,
		Limit:   limit,
		Offset:  offset,
		HasMore: len(commits) == limit,
	}
}

func commit2Vo(c git.Commit) reqvo.CommitVO {
	return reqvo.CommitVO{
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
}
