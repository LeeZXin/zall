package reposrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/actiontaskmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/actionsrv"
	"github.com/LeeZXin/zall/git/modules/service/gpgkeysrv"
	"github.com/LeeZXin/zall/git/modules/service/sshkeysrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/opsrv"
	"github.com/LeeZXin/zall/meta/modules/service/teamsrv"
	"github.com/LeeZXin/zall/pkg/action"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/git"
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
	"gopkg.in/yaml.v3"
	"hash/crc32"
	"path"
	"path/filepath"
	"strconv"
	"time"
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

// GetById 通过id获取仓库信息
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

func (*innerImpl) CheckAccessToken(ctx context.Context, reqDTO CheckAccessTokenReqDTO) bool {
	if err := reqDTO.IsValid(); err != nil {
		return false
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	token, b, err := repomd.GetAccessToken(ctx, repomd.GetAccessTokenReqDTO{
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
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return TreeDTO{}, err
	}
	if p.GetRepoPerm(repo.Id).CanAccessRepo {
		return TreeDTO{}, util.UnauthorizedError()
	}
	resp, err := client.EntriesRepo(ctx, reqvo.EntriesRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
		Offset:   reqDTO.Offset,
	}, repo.NodeId)
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
	p, b, err := teammd.GetTeamUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.UnauthorizedError()
	}
	// 项目管理员可看到所有仓库或者应用所有仓库权限配置
	if p.IsAdmin || len(p.PermDetail.RepoPermList) == 0 {
		repoList, err := repomd.ListAllRepo(ctx, reqDTO.TeamId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		return repoList, nil
	}
	// 通过可访问仓库id查询
	permList := p.PermDetail.RepoPermList
	IdList, _ := listutil.Map(permList, func(t perm.RepoPermWithId) (int64, error) {
		return t.RepoId, nil
	})
	repoList, err := repomd.ListRepoByIdList(ctx, IdList)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
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
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return CatFileRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return CatFileRespDTO{}, util.UnauthorizedError()
	}
	resp, err := client.CatFile(ctx, reqvo.CatFileReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
		FileName: reqDTO.FileName,
	}, repo.NodeId)
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
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return TreeRepoRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return TreeRepoRespDTO{}, util.UnauthorizedError()
	}
	resp, err := client.TreeRepo(ctx, reqvo.TreeRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
	}, repo.NodeId)
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

// InitRepo 初始化仓库
func (s *outerImpl) InitRepo(ctx context.Context, reqDTO InitRepoReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.InitRepo),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验项目信息
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
	if !b {
		err = util.UnauthorizedError()
		return
	}
	// 是否可创建项目
	if !p.PermDetail.TeamPerm.CanInitRepo {
		err = util.UnauthorizedError()
		return
	}
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
	// 选择git节点
	nodes, err := gitnodemd.GetAll(ctx)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if len(nodes) == 0 {
		logger.Logger.WithContext(ctx).Error("empty git nodes")
		err = util.NewBizErr(apicode.EmptyGitNodesErrCode, i18n.EmptyGitNodesError)
		return
	}
	nodeIndex := int(crc32.ChecksumIEEE([]byte(reqDTO.Name))) % len(nodes)
	// 添加数据
	insertReq := repomd.InsertRepoReqDTO{
		Name:          reqDTO.Name,
		Path:          relativePath,
		Author:        reqDTO.Operator.Account,
		TeamId:        reqDTO.TeamId,
		RepoDesc:      reqDTO.Desc,
		DefaultBranch: reqDTO.DefaultBranch,
		NodeId:        nodes[nodeIndex].NodeId,
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
			CreateReadme:  reqDTO.CreateReadme,
			GitIgnoreName: reqDTO.GitIgnoreName,
			DefaultBranch: reqDTO.DefaultBranch,
		}, nodes[nodeIndex].NodeId)
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
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 是否可删除权限
	if !p.TeamPerm.CanDeleteRepo {
		err = util.UnauthorizedError()
		return
	}
	// 检查特殊配置
	// 查询team角色是否包含该Id
	groups, err := teammd.ListTeamUserGroup(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 项目组仍有该仓库的特殊配置
	for _, group := range groups {
		for _, repoPerm := range group.GetPermDetail().RepoPermList {
			if repoPerm.RepoId == reqDTO.Id {
				err = util.NewBizErr(apicode.OperationFailedErrCode, i18n.RepoPermsContainsTargetRepoId)
				return
			}
		}
	}
	err = client.DeleteRepo(ctx, reqvo.DeleteRepoReq{
		RepoPath: repo.Path,
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	// 是否可访问
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return nil, util.UnauthorizedError()
	}
	branches, err := client.GetAllBranches(ctx, reqvo.GetAllBranchesReq{
		RepoPath: repo.Path,
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	// 是否可访问
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return nil, util.UnauthorizedError()
	}
	tags, err := client.GetAllTags(ctx, reqvo.GetAllTagsReq{
		RepoPath: repo.Path,
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 是否可访问
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return util.UnauthorizedError()
	}
	err = client.Gc(ctx, reqvo.GcReq{
		RepoPath: repo.Path,
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return DiffCommitsRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return DiffCommitsRespDTO{}, util.UnauthorizedError()
	}
	refs, err := client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath: repo.Path,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return DiffFileRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return DiffFileRespDTO{}, util.UnauthorizedError()
	}
	resp, err := client.DiffFile(ctx, reqvo.DiffFileReq{
		RepoPath: repo.Path,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		FileName: reqDTO.FileName,
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return nil, util.UnauthorizedError()
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
	}, repo.NodeId)
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
	// 校验权限
	repo, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return HistoryCommitsRespDTO{}, err
	}
	if !p.GetRepoPerm(repo.Id).CanAccessRepo {
		return HistoryCommitsRespDTO{}, util.UnauthorizedError()
	}
	resp, err := client.HistoryCommits(ctx, reqvo.HistoryCommitsReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Offset:   reqDTO.Cursor,
	}, repo.NodeId)
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

func (*outerImpl) InsertAccessToken(ctx context.Context, reqDTO InsertAccessTokenReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.InsertAccessToken),
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
	_, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 是否可编辑token
	if !p.GetRepoPerm(reqDTO.Id).CanUpdateToken {
		err = util.UnauthorizedError()
		return
	}
	for i := 0; i < 10; i++ {
		_, err = repomd.InsertAccessToken(ctx, repomd.InsertAccessTokenReqDTO{
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

func (*outerImpl) DeleteAccessToken(ctx context.Context, reqDTO DeleteAccessTokenReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.InsertAccessToken),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	accessToken, b, err := repomd.GetByTid(ctx, reqDTO.Id)
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
	_, p, err := getPerm(ctx, accessToken.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 是否可编辑token
	if !p.GetRepoPerm(accessToken.Id).CanUpdateToken {
		err = util.UnauthorizedError()
		return
	}
	_, err = repomd.DeleteAccessToken(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListAccessToken(ctx context.Context, reqDTO ListAccessTokenReqDTO) ([]AccessTokenDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	// 是否可访问token
	if !p.GetRepoPerm(reqDTO.Id).CanAccessToken {
		return nil, util.UnauthorizedError()
	}
	tokens, err := repomd.ListAccessToken(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tokens, func(t repomd.AccessToken) (AccessTokenDTO, error) {
		return AccessTokenDTO{
			Id:      t.Id,
			Account: t.Account,
			Token:   t.Token,
			Created: t.Created,
		}, nil
	})
}

func getPerm(ctx context.Context, Id int64, operator apisession.UserInfo) (repomd.RepoInfo, perm.Detail, error) {
	repo, b := Inner.GetByRepoId(ctx, Id)
	if !b {
		return repomd.RepoInfo{}, perm.Detail{}, util.InvalidArgsError()
	}
	p, b := teamsrv.Inner.GetTeamUserPermDetail(ctx, repo.TeamId, operator.Account)
	if !b {
		return repo, perm.Detail{}, util.UnauthorizedError()
	}
	return repo, p.PermDetail, nil
}

func (*outerImpl) InsertAction(ctx context.Context, reqDTO InsertActionReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.InsertAction),
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
	_, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 是否可编辑action
	if !p.GetRepoPerm(reqDTO.Id).CanUpdateAction {
		err = util.UnauthorizedError()
		return
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.ActionContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidActionContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	err = repomd.InsertAction(ctx, repomd.InsertActionReqDTO{
		RepoId:         reqDTO.Id,
		AssignInstance: reqDTO.AssignInstance,
		Content:        string(yamlOut),
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) DeleteAction(ctx context.Context, reqDTO DeleteActionReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.DeleteAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := repomd.GetByActionId(ctx, reqDTO.Id)
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
	_, p, err := getPerm(ctx, repoAction.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 是否可编辑action
	if !p.GetRepoPerm(repoAction.Id).CanUpdateAction {
		err = util.UnauthorizedError()
		return
	}
	_, err = repomd.DeleteAction(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

func (*outerImpl) ListAction(ctx context.Context, reqDTO ListActionReqDTO) ([]repomd.Action, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	_, p, err := getPerm(ctx, reqDTO.Id, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	// 是否可访问action
	if !p.GetRepoPerm(reqDTO.Id).CanAccessAction {
		return nil, util.UnauthorizedError()
	}
	ret, err := repomd.ListAction(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return ret, nil
}

func (*outerImpl) UpdateAction(ctx context.Context, reqDTO UpdateActionReqDTO) (err error) {
	// 插入日志
	defer func() {
		opsrv.Inner.InsertOpLog(ctx, opsrv.InsertOpLogReqDTO{
			Account:    reqDTO.Operator.Account,
			OpDesc:     i18n.GetByKey(i18n.RepoSrvKeysVO.UpdateAction),
			ReqContent: reqDTO,
			Err:        err,
		})
	}()
	if err = reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := repomd.GetByActionId(ctx, reqDTO.Id)
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
	_, p, err := getPerm(ctx, repoAction.Id, reqDTO.Operator)
	if err != nil {
		return
	}
	// 是否可编辑action
	if !p.GetRepoPerm(repoAction.Id).CanUpdateAction {
		err = util.UnauthorizedError()
		return
	}
	var graph action.GraphCfg
	err = yaml.Unmarshal([]byte(reqDTO.ActionContent), &graph)
	if err != nil || graph.IsValid() != nil {
		err = util.NewBizErr(apicode.InvalidArgsCode, i18n.InvalidActionContent)
		return
	}
	yamlOut, _ := yaml.Marshal(graph)
	_, err = repomd.UpdateAction(ctx, repomd.UpdateActionReqDTO{
		Id:             reqDTO.Id,
		Content:        string(yamlOut),
		AssignInstance: reqDTO.AssignInstance,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
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
			}, repo.NodeId)
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
		}
	}()
	return
}

func (*outerImpl) TriggerAction(ctx context.Context, reqDTO TriggerActionReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repoAction, b, err := repomd.GetByActionId(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	repo, p, err := getPerm(ctx, repoAction.Id, reqDTO.Operator)
	if err != nil {
		return err
	}
	// 是否可编辑action
	if !p.GetRepoPerm(repoAction.Id).CanTriggerAction {
		return util.UnauthorizedError()
	}
	req := action.Webhook{
		RepoId:    repo.Id,
		RepoName:  repo.Name,
		Ref:       reqDTO.Ref,
		EventTime: time.Now().UnixMilli(),
		Operator: git.User{
			Account: reqDTO.Operator.Account,
			Email:   reqDTO.Operator.Email,
		},
		TriggerType: actiontaskmd.ManualTriggerType.Int(),
		YamlContent: repoAction.Content,
	}
	// 负载均衡选择一个节点
	instance, b, err := actionsrv.SelectAndIncrJobCountInstances(context.Background(), repoAction.AssignInstance)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 没有可用节点
	if !b {
		return util.NewBizErr(apicode.ActionInstanceNotFoundCode, i18n.ActionInstanceNotFound)
	}
	action.TriggerActionHook(req, instance.InstanceHost)
	return nil
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
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.Id)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	if !b {
		err = util.InvalidArgsError()
		return
	}
	// 查询team角色是否包含该Id
	groups, err := teammd.ListTeamUserGroup(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	// 项目组仍有该仓库的特殊配置
	for _, group := range groups {
		for _, repoPerm := range group.GetPermDetail().RepoPermList {
			if repoPerm.RepoId == reqDTO.Id {
				err = util.NewBizErr(apicode.OperationFailedErrCode, i18n.RepoPermsContainsTargetRepoId)
				return
			}
		}
	}
	_, err = repomd.TransferTeam(ctx, reqDTO.Id, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}
