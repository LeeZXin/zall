package pullrequestsrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	// GetStats 获取统计详情
	GetStats(context.Context, GetStatsReqDTO) (GetStatsRespDTO, error)
	// GetPullRequest 查询合并请求
	GetPullRequest(context.Context, GetPullRequestReqDTO) (PullRequestDTO, error)
	// ListPullRequest 查询合并请求列表
	ListPullRequest(context.Context, ListPullRequestReqDTO) ([]PullRequestDTO, int64, error)
	// SubmitPullRequest 创建合并请求
	SubmitPullRequest(context.Context, SubmitPullRequestReqDTO) error
	// ClosePullRequest 关闭合并请求
	ClosePullRequest(context.Context, ClosePullRequestReqDTO) (bool, error)
	// MergePullRequest 合并代码
	MergePullRequest(context.Context, MergePullRequestReqDTO) (bool, error)
	// CanMergePullRequest 是否可合并
	CanMergePullRequest(context.Context, CanMergePullRequestReqDTO) (CanMergePullRequestRespDTO, bool, error)
	// CanReviewPullRequest 是否可评审代码
	CanReviewPullRequest(context.Context, CanReviewPullRequestReqDTO) (CanReviewPullRequestRespDTO, bool, error)
	// AgreeReviewPullRequest 同意合并代码
	AgreeReviewPullRequest(context.Context, AgreeReviewPullRequestReqDTO) (bool, error)
	// ListTimeline 展示时间轴
	ListTimeline(context.Context, ListTimelineReqDTO) ([]TimelineDTO, error)
	// AddComment 添加评论
	AddComment(context.Context, AddCommentReqDTO) error
	// DeleteComment 删除评论
	DeleteComment(context.Context, DeleteCommentReqDTO) error
	// ListReview 评审记录
	ListReview(context.Context, ListReviewReqDTO) ([]ReviewDTO, error)
}
