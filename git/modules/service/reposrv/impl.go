package reposrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/gpgkeysrv"
	"github.com/LeeZXin/zall/git/modules/service/sshkeysrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/strutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/patrickmn/go-cache"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const (
	accessRepo = iota
	updateToken
	accessToken
)

const (
	createRepo = iota
	deleteRepo
)

type innerImpl struct {
	pathCache *cache.Cache
	idCache   *cache.Cache
}

// GetByRepoPath 通过相对路径获取仓库信息
func (s *innerImpl) GetByRepoPath(ctx context.Context, path string) (repomd.RepoInfo, bool) {
	v, b := s.pathCache.Get(path)
	if b {
		r := v.(repomd.RepoInfo)
		return r, r.Id != 0
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	r, b, err := repomd.GetByPath(ctx, path)
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		s.pathCache.Set(path, r, time.Second)
	} else {
		s.pathCache.Set(path, r, time.Minute)
	}
	return r, b
}

// GetByRepoId 通过id获取仓库信息
func (s *innerImpl) GetByRepoId(ctx context.Context, id int64) (repomd.RepoInfo, bool) {
	key := strconv.FormatInt(id, 10)
	v, b := s.idCache.Get(key)
	if b {
		r := v.(repomd.RepoInfo)
		return r, r.Id != 0
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	r, b, err := repomd.GetByRepoId(ctx, id)
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		s.idCache.Set(key, r, time.Second)
	} else {
		s.idCache.Set(key, r, time.Minute)
	}
	return r, b
}

func (*innerImpl) CheckRepoToken(ctx context.Context, reqDTO CheckRepoTokenReqDTO) bool {
	if err := reqDTO.IsValid(); err != nil {
		return false
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	token, b, err := repomd.GetRepoToken(ctx, repomd.GetRepoTokenReqDTO{
		RepoId:  reqDTO.Id,
		Account: reqDTO.Account,
	})
	if err != nil || !b {
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
		return false
	}
	return token.Token == reqDTO.Token
}

type outerImpl struct {
}

func (s *outerImpl) EntriesRepo(ctx context.Context, reqDTO EntriesRepoReqDTO) (TreeDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TreeDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return TreeDTO{}, util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return TreeDTO{}, err
	}
	resp, err := client.EntriesRepo(ctx, reqvo.EntriesRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
		Offset:   reqDTO.Offset,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TreeDTO{}, util.InternalError(err)
	}
	return tree2Dto(resp), nil
}

// ListRepo 展示仓库列表
func (*outerImpl) ListRepo(ctx context.Context, reqDTO ListRepoReqDTO) ([]repomd.Repo, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		repoList []repomd.Repo
		err      error
	)
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if !b {
		return nil, util.UnauthorizedError()
	}
	// 没有任何访问仓库权限
	if len(p.PermDetail.RepoPermList) == 0 && !p.PermDetail.DefaultRepoPerm.CanAccessRepo {
		return nil, nil
	}
	// 访问部分仓库
	if len(p.PermDetail.RepoPermList) > 0 {
		repoPermList, _ := listutil.Filter(p.PermDetail.RepoPermList, func(p perm.RepoPermWithId) (bool, error) {
			return p.CanAccessRepo, nil
		})
		// 存在部分可读仓库权限
		if len(repoPermList) > 0 {
			idList, _ := listutil.Map(repoPermList, func(t perm.RepoPermWithId) (int64, error) {
				return t.RepoId, nil
			})
			repoList, err = repomd.GetRepoByIdList(ctx, idList)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return nil, util.InternalError(err)
			}
		}
	} else {
		repoList, err = repomd.GetRepoListByTeamId(ctx, reqDTO.TeamId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
	}
	return repoList, nil
}

// CatFile 展示文件内容
func (s *outerImpl) CatFile(ctx context.Context, reqDTO CatFileReqDTO) (CatFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CatFileRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return CatFileRespDTO{}, util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return CatFileRespDTO{}, err
	}
	resp, err := client.CatFile(ctx, reqvo.CatFileReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
		FileName: reqDTO.FileName,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return CatFileRespDTO{}, util.InternalError(err)
	}
	return CatFileRespDTO{
		FileMode: resp.FileMode,
		ModeName: resp.ModeName,
		Content:  resp.Content,
	}, nil
}

// TreeRepo 代码基本数据
func (s *outerImpl) TreeRepo(ctx context.Context, reqDTO TreeRepoReqDTO) (TreeRepoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return TreeRepoRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return TreeRepoRespDTO{}, util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return TreeRepoRespDTO{}, err
	}
	resp, err := client.TreeRepo(ctx, reqvo.TreeRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return TreeRepoRespDTO{}, util.InternalError(err)
	}
	return TreeRepoRespDTO{
		ReadmeText:   resp.ReadmeText,
		HasReadme:    resp.HasReadme,
		LatestCommit: commit2Dto(resp.LatestCommit),
		Tree:         tree2Dto(resp.Tree),
	}, nil
}

func tree2Dto(vo reqvo.TreeVO) TreeDTO {
	ret := TreeDTO{
		Limit:   vo.Limit,
		Offset:  vo.Offset,
		HasMore: vo.HasMore,
	}
	ret.Files, _ = listutil.Map(vo.Files, func(t reqvo.FileVO) (FileDTO, error) {
		return FileDTO{
			Mode:    t.Mode,
			RawPath: t.RawPath,
			Path:    t.Path,
			Commit:  commit2Dto(t.Commit),
		}, nil
	})
	return ret
}

// CreateRepo 初始化仓库
func (s *outerImpl) CreateRepo(ctx context.Context, reqDTO CreateRepoReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.CreateRepo),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err = checkTeamPerm(ctx, reqDTO.TeamId, reqDTO.Operator, createRepo)
	if err != nil {
		return
	}
	var b bool
	// 相对路径
	relativePath := filepath.Join("zgit", reqDTO.Name+".git")
	_, b, err = repomd.GetByPath(ctx, relativePath)
	// 数据库异常
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 仓库已存在 不能添加
	if b {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.RepoAlreadyExists)
		return
	}
	// 添加数据
	insertReq := repomd.InsertRepoReqDTO{
		Name:          reqDTO.Name,
		Path:          relativePath,
		Author:        reqDTO.Operator.Account,
		TeamId:        reqDTO.TeamId,
		RepoDesc:      reqDTO.Desc,
		DefaultBranch: reqDTO.DefaultBranch,
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 插入数据库
		_, err := repomd.InsertRepo(ctx, insertReq)
		if err != nil {
			return err
		}
		// 调用store
		err = client.InitRepo(ctx, reqvo.InitRepoReq{
			UserAccount:   reqDTO.Operator.Account,
			UserEmail:     reqDTO.Operator.Email,
			RepoName:      reqDTO.Operator.Name,
			RepoPath:      relativePath,
			AddReadme:     reqDTO.AddReadme,
			GitIgnoreName: reqDTO.GitIgnoreName,
			DefaultBranch: reqDTO.DefaultBranch,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
	}
	return
}

// AllGitIgnoreTemplateList 所有gitignore模版名称
func (*outerImpl) AllGitIgnoreTemplateList() []string {
	return gitignoreSet.AllKeys()
}

// DeleteRepo 删除仓库
func (s *outerImpl) DeleteRepo(ctx context.Context, reqDTO DeleteRepoReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.DeleteRepo),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		err = util.InvalidArgsError()
		return
	}
	err = checkTeamPerm(ctx, repo.TeamId, reqDTO.Operator, deleteRepo)
	if err != nil {
		return
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		err := client.DeleteRepo(ctx, reqvo.DeleteRepoReq{
			RepoPath: repo.Path,
		})
		if err != nil {
			return err
		}
		_, err = repomd.DeleteRepo(ctx, reqDTO.Id)
		return err
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// AllBranches 仓库所有分支
func (s *outerImpl) AllBranches(ctx context.Context, reqDTO AllBranchesReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	branches, err := client.GetAllBranches(ctx, reqvo.GetAllBranchesReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return branches, nil
}

// AllTags 仓库所有tag
func (s *outerImpl) AllTags(ctx context.Context, reqDTO AllTagsReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	tags, err := client.GetAllTags(ctx, reqvo.GetAllTagsReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return tags, nil
}

// Gc git gc
func (s *outerImpl) Gc(ctx context.Context, reqDTO GcReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	err = client.Gc(ctx, reqvo.GcReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

func (s *outerImpl) DiffCommits(ctx context.Context, reqDTO DiffCommitsReqDTO) (DiffCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffCommitsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return DiffCommitsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffCommitsRespDTO{}, err
	}
	refs, err := client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath: repo.Path,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffCommitsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffCommitsRespDTO{}, util.InternalError(err)
	}
	ret := DiffCommitsRespDTO{
		Target:       refs.Target,
		Head:         refs.Head,
		TargetCommit: commit2Dto(refs.TargetCommit),
		HeadCommit:   commit2Dto(refs.HeadCommit),
		NumFiles:     refs.NumFiles,
		DiffNumsStats: DiffNumsStatInfoDTO{
			FileChangeNums: refs.DiffNumsStats.FileChangeNums,
			InsertNums:     refs.DiffNumsStats.InsertNums,
			DeleteNums:     refs.DiffNumsStats.DeleteNums,
		},
		ConflictFiles: refs.ConflictFiles,
	}
	ret.DiffNumsStats.Stats, _ = listutil.Map(refs.DiffNumsStats.Stats, func(t reqvo.DiffNumsStatVO) (DiffNumsStatDTO, error) {
		return DiffNumsStatDTO{
			RawPath:    t.Path,
			Path:       path.Base(t.Path),
			TotalNums:  t.TotalNums,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	ret.Commits, _ = listutil.Map(refs.Commits, func(t reqvo.CommitVO) (CommitDTO, error) {
		return commit2Dto(t), nil
	})
	ret.CanMerge = len(ret.Commits) > 0 && len(ret.ConflictFiles) == 0
	return ret, nil
}

func (s *outerImpl) DiffFile(ctx context.Context, reqDTO DiffFileReqDTO) (DiffFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffFileRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return DiffFileRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffFileRespDTO{}, err
	}
	resp, err := client.DiffFile(ctx, reqvo.DiffFileReq{
		RepoPath: repo.Path,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		FileName: reqDTO.FileName,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffFileRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffFileRespDTO{}, util.InternalError(err)
	}
	ret := DiffFileRespDTO{
		FilePath:    resp.FilePath,
		OldMode:     resp.OldMode,
		Mode:        resp.Mode,
		IsSubModule: resp.IsSubModule,
		FileType:    resp.FileType,
		IsBinary:    resp.IsBinary,
		RenameFrom:  resp.RenameFrom,
		RenameTo:    resp.RenameTo,
		CopyFrom:    resp.CopyFrom,
		CopyTo:      resp.CopyTo,
	}
	ret.Lines, _ = listutil.Map(resp.Lines, func(t reqvo.DiffLineVO) (DiffLineDTO, error) {
		return DiffLineDTO{
			Index:   t.Index,
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}, nil
	})
	return ret, nil
}

func commit2Dto(commit reqvo.CommitVO) CommitDTO {
	return CommitDTO{
		Author: UserDTO{
			Account: commit.Author.Account,
			Email:   commit.Author.Email,
		},
		Committer: UserDTO{
			Account: commit.Committer.Account,
			Email:   commit.Committer.Email,
		},
		AuthoredTime:  commit.AuthoredTime,
		CommittedTime: commit.CommittedTime,
		CommitMsg:     commit.CommitMsg,
		CommitId:      commit.CommitId,
		ShortId:       util.LongCommitId2ShortId(commit.CommitId),
	}
}

func (s *outerImpl) ShowDiffTextContent(ctx context.Context, reqDTO ShowDiffTextContentReqDTO) ([]DiffLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	var startLine int
	if reqDTO.Direction == UpDirection {
		if reqDTO.Limit < 0 {
			startLine = 0
		} else {
			startLine = reqDTO.Offset - reqDTO.Limit
		}
	} else {
		startLine = reqDTO.Offset
	}
	if startLine < 0 {
		startLine = 0
	}
	lines, err := client.ShowDiffTextContent(ctx, reqvo.ShowDiffTextContentReq{
		RepoPath:  repo.Path,
		CommitId:  reqDTO.CommitId,
		FileName:  reqDTO.FileName,
		StartLine: startLine,
		Limit:     reqDTO.Limit,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return nil, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(lines, func(t reqvo.DiffLineVO) (DiffLineDTO, error) {
		return DiffLineDTO{
			Index:   t.Index,
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}, nil
	})
}

func (s *outerImpl) HistoryCommits(ctx context.Context, reqDTO HistoryCommitsReqDTO) (HistoryCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return HistoryCommitsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.Id)
	if !b {
		return HistoryCommitsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return HistoryCommitsRespDTO{}, err
	}
	resp, err := client.HistoryCommits(ctx, reqvo.HistoryCommitsReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Offset:   reqDTO.Cursor,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return HistoryCommitsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return HistoryCommitsRespDTO{}, err
	}
	// 缓存 用于校验签名
	gpgMap := make(map[string][]openpgp.EntityList)
	sshMap := make(map[string][]string)
	ret := HistoryCommitsRespDTO{
		Cursor: resp.Cursor,
	}
	ret.Data, _ = listutil.Map(resp.Data, func(t reqvo.CommitVO) (CommitDTO, error) {
		r := commit2Dto(t)
		if t.CommitSig != "" {
			sig := signature.CommitSig(t.CommitSig)
			if sig.IsSSHSig() {
				sshKeys, b := sshMap[t.Committer.Account]
				if !b {
					verified, err := sshkeysrv.Inner.GetVerifiedByAccount(ctx, t.Committer.Account)
					if err != nil {
						verified = []string{}
						sshMap[t.Committer.Account] = verified
					}
					sshKeys = verified
				}
				for _, key := range sshKeys {
					if e := signature.VerifySshSignature(t.CommitSig, t.Payload, key); e == nil {
						r.Verified = true
						break
					}
				}
			} else if sig.IsGPGSig() {
				gpgKeys, b := gpgMap[t.Committer.Account]
				if !b {
					verified, err := gpgkeysrv.Inner.GetVerifiedByAccount(ctx, t.Committer.Account)
					if err != nil {
						verified = []openpgp.EntityList{}
						gpgMap[t.Committer.Account] = verified
					}
					gpgKeys = verified
				}
				for _, keys := range gpgKeys {
					if _, err := signature.CheckArmoredDetachedSignature(keys, t.Payload, t.CommitSig); err == nil {
						r.Verified = true
						break
					}
				}
			}
		}
		return r, nil
	})
	return ret, nil
}

func (*outerImpl) InsertRepoToken(ctx context.Context, reqDTO InsertRepoTokenReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.InsertRepoToken),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err = checkPermByRepoId(ctx, reqDTO.Id, reqDTO.Operator, updateToken)
	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		_, err = repomd.InsertRepoToken(ctx, repomd.InsertRepoTokenReqDTO{
			RepoId:  reqDTO.Id,
			Account: strutil.RandomStr(16),
			Token:   strutil.RandomStr(16),
		})
		if err != nil && xormutil.IsDuplicatedEntryError(err) {
			continue
		}
		break
	}
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteRepoToken(ctx context.Context, reqDTO DeleteRepoTokenReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.DeleteRepoToken),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	token, b, err := repomd.GetByTokenId(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 校验权限
	err = checkPermByRepoId(ctx, token.RepoId, reqDTO.Operator, updateToken)
	if err != nil {
		return
	}
	_, err = repomd.DeleteRepoToken(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListRepoToken(ctx context.Context, reqDTO ListRepoTokenReqDTO) ([]RepoTokenDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	err := checkPermByRepoId(ctx, reqDTO.Id, reqDTO.Operator, accessToken)
	if err != nil {
		return nil, err
	}
	tokens, err := repomd.ListRepoToken(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tokens, func(t repomd.RepoToken) (RepoTokenDTO, error) {
		return RepoTokenDTO{
			Id:      t.Id,
			Account: t.Account,
			Token:   t.Token,
			Created: t.Created,
		}, nil
	})
}

func checkPermByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo, permCode int) error {
	repo, b := Inner.GetByRepoId(ctx, repoId)
	if !b {
		return util.InvalidArgsError()
	}
	return checkPermByRepo(ctx, repo, operator, permCode)
}

func checkPermByRepo(ctx context.Context, repo repomd.RepoInfo, operator apisession.UserInfo, permCode int) error {
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case accessRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo
	case accessToken:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanAccessToken
	case updateToken:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanUpdateToken
	}
	if pass {
		return nil
	}
	return util.UnauthorizedError()
}

func checkTeamPerm(ctx context.Context, teamId int64, operator apisession.UserInfo, permCode int) error {
	_, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return nil
	}
	p, b := teamsrv.Inner.GetUserPermDetail(ctx, teamId, operator.Account)
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case createRepo:
		pass = p.PermDetail.TeamPerm.CanCreateRepo
	case deleteRepo:
		pass = p.PermDetail.TeamPerm.CanDeleteRepo
	}
	if pass {
		return nil
	}
	return util.UnauthorizedError()
}

func (s *outerImpl) RefreshAllGitHooks(ctx context.Context, reqDTO RefreshAllGitHooksReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.RefreshAllGitHooks),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	// 没有权限
	if !reqDTO.Operator.IsAdmin {
		return util.UnauthorizedError()
	}
	go func() {
		ctx, closer := xormstore.Context(context.Background())
		defer closer.Close()
		err := repomd.IterateRepo(ctx, func(repo *repomd.Repo) error {
			return client.InitRepoHook(context.Background(), reqvo.InitRepoHookReq{
				RepoPath: repo.Path,
			})
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	}()
	return
}

func (*outerImpl) TransferTeam(ctx context.Context, reqDTO TransferTeamReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.TransferTeam),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	// 只有系统管理员才有权限
	if !reqDTO.Operator.IsAdmin {
		err = util.UnauthorizedError()
		return
	}
	_, err = repomd.TransferTeam(ctx, reqDTO.Id, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}
