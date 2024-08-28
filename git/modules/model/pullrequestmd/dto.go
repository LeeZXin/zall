package pullrequestmd

import (
	"github.com/LeeZXin/zall/pkg/git"
)

type InsertPullRequestReqDTO struct {
	RepoId       int64
	Target       string
	TargetType   git.RefType
	Head         string
	HeadType     git.RefType
	CreateBy     string
	Title        string
	PrStatus     PrStatus
	CommentCount int
}

type InsertReviewReqDTO struct {
	PrId     int64
	Status   ReviewStatus
	Reviewer string
}

type UpdateReviewReqDTO struct {
	Id     int64
	Status ReviewStatus
}

type ExistsPrByRepoIdAndRefReqDTO struct {
	RepoId     int64
	Head       string
	HeadType   git.RefType
	Target     string
	TargetType git.RefType
	Status     PrStatus
}

type ListPullRequestReqDTO struct {
	RepoId    int64
	SearchKey string
	Status    PrStatus
	PageNum   int
	PageSize  int
}

type GroupByPrStatusDTO struct {
	PrStatus   PrStatus
	TotalCount int64
}

type InsertTimelineReqDTO struct {
	PrId    int64
	Action  Action
	Account string
}
