package pullrequestsrv

import (
	"context"
)

var (
	Outer OuterService = new(outerImpl)
)

type OuterService interface {
	SubmitPullRequest(context.Context, SubmitPullRequestReqDTO) error
	ClosePullRequest(context.Context, ClosePullRequestReqDTO) error
	MergePullRequest(context.Context, MergePullRequestReqDTO) error
	ReviewPullRequest(context.Context, ReviewPullRequestReqDTO) error
}
