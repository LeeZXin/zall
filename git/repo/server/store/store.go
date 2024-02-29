package store

import (
	"context"
	"github.com/LeeZXin/zall/git/repo/reqvo"
)

type Store interface {
	// InitRepo 初始化仓库
	InitRepo(context.Context, reqvo.InitRepoReq) error
	// DeleteRepo 删除仓库
	DeleteRepo(context.Context, reqvo.DeleteRepoReq) error
	// GetAllBranches 获取所有分支
	GetAllBranches(context.Context, reqvo.GetAllBranchesReq) ([]string, error)
	// GetAllTags 获取所有tags
	GetAllTags(context.Context, reqvo.GetAllTagsReq) ([]string, error)
	// Gc 触发仓库gc
	Gc(context.Context, reqvo.GcReq) error
	// DiffRefs 对比两个ref差异
	DiffRefs(context.Context, reqvo.DiffRefsReq) (reqvo.DiffRefsResp, error)
	// DiffFile 对比两个分支单个文件差异
	DiffFile(context.Context, reqvo.DiffFileReq) (reqvo.DiffFileResp, error)
	// GetRepoSize 获取仓库大小
	GetRepoSize(context.Context, reqvo.GetRepoSizeReq) (int64, error)
	// ShowDiffTextContent 获取某个commitId文件内容
	ShowDiffTextContent(context.Context, reqvo.ShowDiffTextContentReq) ([]reqvo.DiffLineVO, error)
	// HistoryCommits 获取历史提交记录
	HistoryCommits(context.Context, reqvo.HistoryCommitsReq) (reqvo.HistoryCommitsResp, error)
	// InitRepoHook 重建仓库hook
	InitRepoHook(context.Context, reqvo.InitRepoHookReq) error
	// EntriesRepo 仓库文件列表
	EntriesRepo(context.Context, reqvo.EntriesRepoReq) (reqvo.TreeVO, error)
	// CatFile 展示文件内容
	CatFile(context.Context, reqvo.CatFileReq) (reqvo.CatFileResp, error)
	// TreeRepo 仓库首页
	TreeRepo(context.Context, reqvo.TreeRepoReq) (reqvo.TreeRepoResp, error)
	// UploadPack git-upload-pack
	UploadPack(reqvo.UploadPackReq)
	// ReceivePack git-receive-pack
	ReceivePack(reqvo.ReceivePackReq)
	// InfoRefs smart http infoRefs
	InfoRefs(context.Context, reqvo.InfoRefsReq)
	// Merge 合并两个分支
	Merge(context.Context, reqvo.MergeReq) error
}
