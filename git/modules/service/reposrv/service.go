package reposrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/util"
)

var (
	Inner InnerService = &innerImpl{
		pathCache: util.NewGoCache(),
		idCache:   util.NewGoCache(),
	}
	Outer OuterService = new(outerImpl)
)

type InnerService interface {
	GetByRepoPath(context.Context, string) (repomd.Repo, bool)
	GetByRepoId(context.Context, int64) (repomd.Repo, bool)
	CheckRepoToken(context.Context, CheckRepoTokenReqDTO) bool
}

type OuterService interface {
	// SimpleInfo 基本信息
	SimpleInfo(context.Context, SimpleInfoReqDTO) (SimpleInfoRespDTO, error)
	// GetRepo 获取仓库信息
	GetRepo(context.Context, GetRepoReqDTO) (repomd.Repo, bool, error)
	// EntriesRepo 文件列表
	EntriesRepo(context.Context, EntriesRepoReqDTO) ([]BlobDTO, error)
	// ListRepo 获取仓库列表
	ListRepo(context.Context, ListRepoReqDTO) ([]repomd.Repo, error)
	// CatFile 展示文件详细内容 仅展示文本信息
	CatFile(context.Context, CatFileReqDTO) (CatFileRespDTO, error)
	// IndexRepo 代码首页
	IndexRepo(context.Context, IndexRepoReqDTO) (IndexRepoRespDTO, error)
	// CreateRepo 创建仓库
	CreateRepo(context.Context, CreateRepoReqDTO) error
	AllGitIgnoreTemplateList() []string
	DeleteRepo(context.Context, DeleteRepoReqDTO) error
	// AllBranches 获取所有分支
	AllBranches(context.Context, AllBranchesReqDTO) ([]string, error)
	AllTags(context.Context, AllTagsReqDTO) ([]string, error)
	Gc(context.Context, GcReqDTO) error
	DiffCommits(context.Context, DiffCommitsReqDTO) (DiffCommitsRespDTO, error)
	DiffFile(context.Context, DiffFileReqDTO) (DiffFileRespDTO, error)
	ShowDiffTextContent(context.Context, ShowDiffTextContentReqDTO) ([]DiffLineDTO, error)
	HistoryCommits(context.Context, HistoryCommitsReqDTO) (HistoryCommitsRespDTO, error)
	InsertRepoToken(context.Context, InsertRepoTokenReqDTO) error
	DeleteRepoToken(context.Context, DeleteRepoTokenReqDTO) error
	ListRepoToken(context.Context, ListRepoTokenReqDTO) ([]RepoTokenDTO, error)
	RefreshAllGitHooks(context.Context, RefreshAllGitHooksReqDTO) error
	TransferTeam(context.Context, TransferTeamReqDTO) error
	// Blame 获取每一行提交信息
	Blame(context.Context, BlameReqDTO) ([]BlameLineDTO, error)
}
