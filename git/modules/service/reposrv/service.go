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
	GetByRepoPath(context.Context, string) (repomd.RepoInfo, bool)
	GetByRepoId(context.Context, int64) (repomd.RepoInfo, bool)
	CheckRepoToken(context.Context, CheckRepoTokenReqDTO) bool
}

type OuterService interface {
	EntriesRepo(context.Context, EntriesRepoReqDTO) (TreeDTO, error)
	ListRepo(context.Context, ListRepoReqDTO) ([]repomd.Repo, error)
	CatFile(context.Context, CatFileReqDTO) (CatFileRespDTO, error)
	TreeRepo(context.Context, TreeRepoReqDTO) (TreeRepoRespDTO, error)
	InitRepo(context.Context, InitRepoReqDTO) error
	AllGitIgnoreTemplateList() []string
	DeleteRepo(context.Context, DeleteRepoReqDTO) error
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
}
