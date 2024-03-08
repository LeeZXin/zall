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
	CheckAccessToken(context.Context, CheckAccessTokenReqDTO) bool
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
	InsertAccessToken(context.Context, InsertAccessTokenReqDTO) error
	DeleteAccessToken(context.Context, DeleteAccessTokenReqDTO) error
	ListAccessToken(context.Context, ListAccessTokenReqDTO) ([]AccessTokenDTO, error)
	InsertAction(context.Context, InsertActionReqDTO) error
	UpdateAction(context.Context, UpdateActionReqDTO) error
	DeleteAction(context.Context, DeleteActionReqDTO) error
	ListAction(context.Context, ListActionReqDTO) ([]repomd.Action, error)
	RefreshAllGitHooks(context.Context, RefreshAllGitHooksReqDTO) error
	TriggerAction(context.Context, TriggerActionReqDTO) error
	TransferTeam(context.Context, TransferTeamReqDTO) error
}
