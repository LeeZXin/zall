package reposrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/gpgkeysrv"
	"github.com/LeeZXin/zall/git/modules/service/sshkeysrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/limiter"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/patrickmn/go-cache"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	accessRepo = iota
	updateRepo
	accessToken
	updateToken
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
func (s *innerImpl) GetByRepoPath(ctx context.Context, path string) (repomd.Repo, bool) {
	v, b := s.pathCache.Get(path)
	if b {
		r := v.(repomd.Repo)
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
func (s *innerImpl) GetByRepoId(ctx context.Context, id int64) (repomd.Repo, bool) {
	key := strconv.FormatInt(id, 10)
	v, b := s.idCache.Get(key)
	if b {
		r := v.(repomd.Repo)
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

type outerImpl struct {
	CreateArchiveLimiter limiter.Limiter
}

func newOuterImpl() OuterService {
	limit := static.GetInt64("createArchiveLimit")
	if limit <= 0 {
		limit = 10
	}
	return &outerImpl{
		CreateArchiveLimiter: limiter.NewCountLimiter(limit),
	}
}

// GetRepo 获取仓库信息
func (s *outerImpl) GetRepo(ctx context.Context, reqDTO GetRepoReqDTO) (repomd.Repo, bool, error) {
	if err := reqDTO.IsValid(); err != nil {
		return repomd.Repo{}, false, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, false, util.InternalError(err)
	}
	return repo, b, nil
}

func (s *outerImpl) EntriesRepo(ctx context.Context, reqDTO EntriesRepoReqDTO) ([]BlobDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return nil, util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	resp, err := client.EntriesRepo(ctx, reqvo.EntriesRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(resp, func(t reqvo.BlobVO) (BlobDTO, error) {
		return BlobDTO{
			Mode:    t.Mode,
			RawPath: t.RawPath,
			Path:    t.Path,
		}, nil
	})
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
	sort.SliceStable(repoList, func(i, j int) bool {
		return repoList[i].LastOperated.After(repoList[j].LastOperated)
	})
	return repoList, nil
}

// CatFile 展示文件内容
func (s *outerImpl) CatFile(ctx context.Context, reqDTO CatFileReqDTO) (CatFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CatFileRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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
		FilePath: reqDTO.FilePath,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return CatFileRespDTO{}, util.InternalError(err)
	}
	return CatFileRespDTO{
		FileMode: resp.FileMode,
		ModeName: resp.ModeName,
		Content:  resp.Content,
		Size:     resp.Size,
		Commit:   commit2Dto(resp.Commit),
	}, nil
}

// IndexRepo 代码首页
func (s *outerImpl) IndexRepo(ctx context.Context, reqDTO IndexRepoReqDTO) (IndexRepoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return IndexRepoRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return IndexRepoRespDTO{}, util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return IndexRepoRespDTO{}, err
	}
	resp, err := client.IndexRepo(ctx, reqvo.IndexRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return IndexRepoRespDTO{}, util.InternalError(err)
	}
	return IndexRepoRespDTO{
		ReadmeText:   resp.ReadmeText,
		HasReadme:    resp.HasReadme,
		LatestCommit: commit2Dto(resp.LatestCommit),
		Tree:         tree2Dto(resp.Tree),
	}, nil
}

// SimpleInfo 基本信息
func (s *outerImpl) SimpleInfo(ctx context.Context, reqDTO SimpleInfoReqDTO) (SimpleInfoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return SimpleInfoRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return SimpleInfoRespDTO{}, util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return SimpleInfoRespDTO{}, err
	}
	branches, err := client.GetAllBranches(ctx, reqvo.GetAllBranchesReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SimpleInfoRespDTO{}, util.InternalError(err)
	}
	tags, err := client.GetAllTags(ctx, reqvo.GetAllTagsReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SimpleInfoRespDTO{}, util.InternalError(err)
	}
	ret := SimpleInfoRespDTO{}
	ret.Branches, _ = listutil.Map(branches, func(t reqvo.RefVO) (string, error) {
		return t.Name, nil
	})
	ret.Tags, _ = listutil.Map(tags, func(t reqvo.RefVO) (string, error) {
		return t.Name, nil
	})
	cfg, b := cfgsrv.Inner.GetGitCfg(ctx)
	if b {
		if cfg.HttpUrl != "" {
			ret.CloneHttpUrl = strings.TrimSuffix(cfg.HttpUrl, "/") + "/" + repo.Path
		}
		if cfg.SshUrl != "" {
			ret.CloneSshUrl = strings.TrimSuffix(cfg.SshUrl, "/") + "/" + repo.Path
		}
	}
	return ret, nil
}

func tree2Dto(vo reqvo.TreeVO) TreeDTO {
	ret := TreeDTO{}
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
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 插入数据库
		repo, err2 := repomd.InsertRepo(ctx, repomd.InsertRepoReqDTO{
			Name:          reqDTO.Name,
			Path:          relativePath,
			Author:        reqDTO.Operator.Account,
			TeamId:        reqDTO.TeamId,
			RepoDesc:      reqDTO.Desc,
			DefaultBranch: reqDTO.DefaultBranch,
			RepoStatus:    repomd.OpenRepoStatus,
			LastOperated:  time.Now(),
		})
		if err2 != nil {
			return err2
		}
		// 调用store
		gitSize, err2 := client.InitRepo(ctx, reqvo.InitRepoReq{
			UserAccount:   reqDTO.Operator.Account,
			UserEmail:     reqDTO.Operator.Email,
			RepoName:      reqDTO.Name,
			RepoPath:      relativePath,
			AddReadme:     reqDTO.AddReadme,
			GitIgnoreName: reqDTO.GitIgnoreName,
			DefaultBranch: reqDTO.DefaultBranch,
		})
		if err2 == nil {
			repomd.UpdateGitSize(ctx, repo.Id, gitSize)
		}
		return err2
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
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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
		_, err = repomd.DeleteRepo(ctx, reqDTO.RepoId)
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
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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
	return listutil.Map(branches, func(t reqvo.RefVO) (string, error) {
		return t.Name, nil
	})
}

// DeleteBranch 删除分支
func (s *outerImpl) DeleteBranch(ctx context.Context, reqDTO DeleteBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return err
	}
	branches, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if b, _ = branches.IsProtectedBranch(reqDTO.Branch); b {
		return util.InvalidArgsError()
	}
	err = client.DeleteBranch(ctx, reqvo.DeleteBranchReq{
		RepoPath: repo.Path,
		Branch:   reqDTO.Branch,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// DeleteTag 删除tag
func (s *outerImpl) DeleteTag(ctx context.Context, reqDTO DeleteTagReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return util.InvalidArgsError()
	}
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return err
	}
	err = client.DeleteTag(ctx, reqvo.DeleteTagReqVO{
		RepoPath: repo.Path,
		Tag:      reqDTO.Tag,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	return nil
}

// AllTags 仓库所有tag
func (s *outerImpl) AllTags(ctx context.Context, reqDTO AllTagsReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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
	return listutil.Map(tags, func(t reqvo.RefVO) (string, error) {
		return t.Name, nil
	})
}

// Gc git gc
func (s *outerImpl) Gc(ctx context.Context, reqDTO GcReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	gitSize, err := client.Gc(ctx, reqvo.GcReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	repomd.UpdateGitSize(ctx, reqDTO.RepoId, gitSize)
	return nil
}

// DiffRefs 比较分支或tag的不同
func (s *outerImpl) DiffRefs(ctx context.Context, reqDTO DiffRefsReqDTO) (DiffRefsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffRefsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return DiffRefsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffRefsRespDTO{}, err
	}
	refs, err := client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath:   repo.Path,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Head:       reqDTO.Head,
		HeadType:   reqDTO.HeadType,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffRefsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffRefsRespDTO{}, util.InternalError(err)
	}
	ret := DiffRefsRespDTO{
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
			RawPath:    t.RawPath,
			Path:       t.Path,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	ret.Commits, _ = listutil.Map(refs.Commits, func(t reqvo.CommitVO) (CommitDTO, error) {
		return commit2Dto(t), nil
	})
	ret.CanMerge = refs.CanMerge
	return ret, nil
}

// DiffCommits 比较commits不同
func (s *outerImpl) DiffCommits(ctx context.Context, reqDTO DiffCommitsReqDTO) (DiffCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffCommitsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return DiffCommitsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffCommitsRespDTO{}, err
	}
	refs, err := client.DiffCommits(ctx, reqvo.DiffCommitsReq{
		RepoPath: repo.Path,
		CommitId: reqDTO.CommitId,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffCommitsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffCommitsRespDTO{}, util.InternalError(err)
	}
	ret := DiffCommitsRespDTO{
		Commit:   commit2Dto(refs.Commit),
		NumFiles: refs.NumFiles,
		DiffNumsStats: DiffNumsStatInfoDTO{
			FileChangeNums: refs.DiffNumsStats.FileChangeNums,
			InsertNums:     refs.DiffNumsStats.InsertNums,
			DeleteNums:     refs.DiffNumsStats.DeleteNums,
		},
	}
	ret.DiffNumsStats.Stats, _ = listutil.Map(refs.DiffNumsStats.Stats, func(t reqvo.DiffNumsStatVO) (DiffNumsStatDTO, error) {
		return DiffNumsStatDTO{
			RawPath:    t.RawPath,
			Path:       t.Path,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}, nil
	})
	return ret, nil
}

// Blame 获取每一行提交信息
func (s *outerImpl) Blame(ctx context.Context, reqDTO BlameReqDTO) ([]BlameLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	lines, err := client.Blame(ctx, reqvo.BlameReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		FilePath: reqDTO.FilePath,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return nil, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(lines, func(t reqvo.BlameLineVO) (BlameLineDTO, error) {
		return BlameLineDTO{
			Number: t.Number,
			Commit: commit2Dto(t.Commit),
		}, nil
	})
}

func (s *outerImpl) DiffFile(ctx context.Context, reqDTO DiffFileReqDTO) (DiffFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffFileRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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
		FilePath: reqDTO.FilePath,
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
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}, nil
	})
	return ret, nil
}

func commit2Dto(c reqvo.CommitVO) CommitDTO {
	return CommitDTO{
		Parent: c.Parent,
		Author: UserDTO{
			Account: c.Author.Account,
			Email:   c.Author.Email,
		},
		Committer: UserDTO{
			Account: c.Committer.Account,
			Email:   c.Committer.Email,
		},
		AuthoredTime:  c.AuthoredTime,
		CommittedTime: c.CommittedTime,
		CommitMsg:     c.CommitMsg,
		CommitId:      c.CommitId,
		ShortId:       util.LongCommitId2ShortId(c.CommitId),
		Tagger: UserDTO{
			Account: c.Tagger.Account,
			Email:   c.Tagger.Email,
		},
		TaggerTime:   c.TaggerTime,
		ShortTagId:   c.ShortTagId,
		TagCommitMsg: c.TagCommitMsg,
	}
}

func (s *outerImpl) ShowDiffTextContent(ctx context.Context, reqDTO ShowDiffTextContentReqDTO) ([]DiffLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
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

func checkPermByRepo(ctx context.Context, repo repomd.Repo, operator apisession.UserInfo, permCode int) error {
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
	case updateRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanUpdateRepo
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
	_, err = repomd.TransferTeam(ctx, reqDTO.RepoId, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// PageBranchCommits 分页获取分支+提交信息
func (*outerImpl) PageBranchCommits(ctx context.Context, reqDTO PageRefCommitsReqDTO) ([]BranchCommitDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, 0, err
	}
	branches, totalCount, err := client.PageBranchAndLastCommit(ctx, reqvo.PageRefCommitsReq{
		RepoPath: repo.Path,
		PageNum:  reqDTO.PageNum,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var prMap map[string]pullrequestmd.PullRequest
	if len(branches) > 0 {
		heads, _ := listutil.Map(branches, func(t reqvo.RefCommitVO) (string, error) {
			return t.Name, nil
		})
		pullRequests, err := pullrequestmd.GetLastPullRequestByRepoIdAndHead(ctx, reqDTO.RepoId, heads)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		prMap, _ = listutil.CollectToMap(pullRequests, func(t pullrequestmd.PullRequest) (string, error) {
			return t.Head, nil
		}, func(t pullrequestmd.PullRequest) (pullrequestmd.PullRequest, error) {
			return t, nil
		})
	} else {
		prMap = map[string]pullrequestmd.PullRequest{}
	}
	pbList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(branches, func(t reqvo.RefCommitVO) (BranchCommitDTO, error) {
		pr, b := prMap[t.Name]
		isProtectedBranch, _ := pbList.IsProtectedBranch(t.Name)
		ret := BranchCommitDTO{
			Name:              t.Name,
			LastCommit:        commit2Dto(t.Commit),
			IsProtectedBranch: isProtectedBranch,
		}
		if b {
			ret.LastPullRequest = &PullRequestDTO{
				Id:       pr.Id,
				PrStatus: pr.PrStatus,
				PrTitle:  pr.PrTitle,
				Created:  pr.Created,
			}
		}
		return ret, nil
	})
	return data, totalCount, nil
}

// PageTagCommits 分页获取tag+提交信息
func (*outerImpl) PageTagCommits(ctx context.Context, reqDTO PageRefCommitsReqDTO) ([]TagCommitDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, 0, err
	}
	tags, totalCount, err := client.PageTagAndCommit(ctx, reqvo.PageRefCommitsReq{
		RepoPath: repo.Path,
		PageNum:  reqDTO.PageNum,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data, _ := listutil.Map(tags, func(t reqvo.RefCommitVO) (TagCommitDTO, error) {
		return TagCommitDTO{
			Name:   t.Name,
			Commit: commit2Dto(t.Commit),
		}, nil
	})
	return data, totalCount, nil
}

// CreateArchive 下载代码
func (s *outerImpl) CreateArchive(ctx context.Context, reqDTO CreateArchiveReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b := Inner.GetByRepoId(ctx, reqDTO.RepoId)
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err := checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	if s.CreateArchiveLimiter.Borrow() {
		defer s.CreateArchiveLimiter.Return()
		err = client.CreateArchive(reqvo.CreateArchiveReq{
			RepoPath: repo.Path,
			FileName: reqDTO.FileName,
			C:        reqDTO.C,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		return nil
	}
	return util.NewBizErrWithMsg(apicode.TooManyOperationCode, i18n.GetByKey(i18n.SystemTooManyOperation))
}
