package pullrequestapi

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type SubmitPullRequestReqVO struct {
	RepoId     int64       `json:"repoId"`
	Target     string      `json:"target"`
	TargetType git.RefType `json:"targetType"`
	Head       string      `json:"head"`
	HeadType   git.RefType `json:"headType"`
	Title      string      `json:"title"`
	Comment    string      `json:"comment"`
}

type ReviewPullRequestReqVO struct {
	Id        int64                      `json:"id"`
	Status    pullrequestmd.ReviewStatus `json:"status"`
	ReviewMsg string                     `json:"reviewMsg"`
}

type ListPullRequestReqVO struct {
	RepoId    int64                  `json:"repoId"`
	Status    pullrequestmd.PrStatus `json:"status"`
	SearchKey string                 `json:"searchKey"`
	ginutil.Page2Req
}

type PullRequestVO struct {
	Id             int64                  `json:"id"`
	RepoId         int64                  `json:"repoId"`
	Target         string                 `json:"target"`
	TargetType     git.RefType            `json:"targetType"`
	TargetCommitId string                 `json:"targetCommitId"`
	Head           string                 `json:"head"`
	HeadType       git.RefType            `json:"headType"`
	HeadCommitId   string                 `json:"headCommitId"`
	PrStatus       pullrequestmd.PrStatus `json:"prStatus"`
	CreateBy       string                 `json:"createBy"`
	CloseBy        string                 `json:"closeBy"`
	MergeBy        string                 `json:"mergeBy"`
	PrTitle        string                 `json:"prTitle"`
	CommentCount   int                    `json:"commentCount"`
	Created        string                 `json:"created"`
	Closed         string                 `json:"closed"`
	Merged         string                 `json:"merged"`
}

type StatsPullRequestVO struct {
	TotalCount  int64 `json:"totalCount"`
	OpenCount   int64 `json:"openCount"`
	ClosedCount int64 `json:"closedCount"`
	MergedCount int64 `json:"mergedCount"`
}

type TimelineVO struct {
	Id      int64                `json:"id"`
	PrId    int64                `json:"prId"`
	Action  pullrequestmd.Action `json:"action"`
	Account string               `json:"account"`
	Created string               `json:"created"`
}

type AddCommentReqVO struct {
	PrId      int64  `json:"prId"`
	ReplyFrom int64  `json:"replyFrom"`
	Comment   string `json:"comment"`
	HasReply  bool   `json:"hasReply"`
}

type CanMergePullRequestRespVO struct {
	CanMerge                bool     `json:"canMerge"`
	IsProtectedBranch       bool     `json:"isProtectedBranch"`
	ReviewCountWhenCreatePr int      `json:"reviewCountWhenCreatePr"`
	ReviewerList            []string `json:"reviewerList"`
	DirectPushList          []string `json:"directPushList"`
	ReviewCount             int      `json:"reviewCount"`
	GitCanMerge             bool     `json:"gitCanMerge"`
	GitConflictFiles        []string `json:"gitConflictFiles"`
	GitCommitCount          int      `json:"gitCommitCount"`
}
