package reposrv

import (
	"context"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
)

var (
	Inner InnerService
	Outer OuterService
)

func Init() {
	if Inner == nil {
		Inner = new(innerImpl)
		Outer = newOuterImpl()
	}
}

type InnerService interface {
	GetByRepoPath(context.Context, string) (repomd.Repo, bool)
}

type OuterService interface {
	// SimpleInfo 基本信息
	SimpleInfo(context.Context, SimpleInfoReqDTO) (SimpleInfoRespDTO, error)
	// GetRepo 获取仓库信息
	GetRepo(context.Context, GetRepoReqDTO) (RepoDTO, error)
	// UpdateRepo 更新仓库配置
	UpdateRepo(context.Context, UpdateRepoReqDTO) error
	// EntriesRepo 文件列表
	EntriesRepo(context.Context, EntriesRepoReqDTO) ([]BlobDTO, error)
	// ListRepo 获取仓库列表
	ListRepo(context.Context, ListRepoReqDTO) ([]RepoDTO, error)
	// CatFile 展示文件详细内容 仅展示文本信息
	CatFile(context.Context, CatFileReqDTO) (CatFileRespDTO, error)
	// IndexRepo 代码首页
	IndexRepo(context.Context, IndexRepoReqDTO) (IndexRepoRespDTO, error)
	// CreateRepo 创建仓库
	CreateRepo(context.Context, CreateRepoReqDTO) error
	// AllGitIgnoreTemplateList ignore模板
	AllGitIgnoreTemplateList() []string
	// DeleteRepo 删除仓库
	DeleteRepo(context.Context, DeleteRepoReqDTO) error
	// DeleteRepoPermanently 永久删除仓库
	DeleteRepoPermanently(context.Context, DeleteRepoReqDTO) error
	// RecoverFromRecycle 恢复仓库
	RecoverFromRecycle(context.Context, RecoverFromRecycleReqDTO) error
	// AllBranches 获取所有分支
	AllBranches(context.Context, AllBranchesReqDTO) ([]string, error)
	// AllTags 所有的标签
	AllTags(context.Context, AllTagsReqDTO) ([]string, error)
	// Gc gc仓库
	Gc(context.Context, GcReqDTO) error
	// DiffRefs 比较分支或tag的不同
	DiffRefs(context.Context, DiffRefsReqDTO) (DiffRefsRespDTO, error)
	// DiffCommits 比较提交的不同
	DiffCommits(context.Context, DiffCommitsReqDTO) (DiffCommitsRespDTO, error)
	// DiffFile 文件差异内容
	DiffFile(context.Context, DiffFileReqDTO) (DiffFileRespDTO, error)
	// HistoryCommits 分支提交历史
	HistoryCommits(context.Context, HistoryCommitsReqDTO) (HistoryCommitsRespDTO, error)
	// TransferTeam 迁移团队
	TransferTeam(context.Context, TransferTeamReqDTO) error
	// Blame 获取每一行提交信息
	Blame(context.Context, BlameReqDTO) ([]BlameLineDTO, error)
	// PageBranchCommits 分页获取分支+提交信息
	PageBranchCommits(context.Context, PageRefCommitsReqDTO) ([]BranchCommitDTO, int64, error)
	// PageTagCommits 分页获取tag+提交信息
	PageTagCommits(context.Context, PageRefCommitsReqDTO) ([]TagCommitDTO, int64, error)
	// DeleteBranch 删除分支
	DeleteBranch(context.Context, DeleteBranchReqDTO) error
	// CreateArchive 下载代码
	CreateArchive(context.Context, CreateArchiveReqDTO) error
	// DeleteTag 删除tag
	DeleteTag(context.Context, DeleteTagReqDTO) error
	// SetRepoArchivedStatus 归档或非归档仓库
	SetRepoArchivedStatus(context.Context, SetRepoArchivedStatusReqDTO) error
	// ListDeletedRepo 展示已删除仓库
	ListDeletedRepo(context.Context, ListDeletedRepoReqDTO) ([]DeletedRepoDTO, error)
	// ListRepoByAdmin 管理员展示仓库列表
	ListRepoByAdmin(context.Context, ListRepoByAdminReqDTO) ([]SimpleRepoDTO, error)
}
