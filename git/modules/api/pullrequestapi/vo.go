package pullrequestapi

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/util"
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

type ListPullRequestReqVO struct {
	RepoId    int64                  `json:"repoId"`
	Status    pullrequestmd.PrStatus `json:"status"`
	SearchKey string                 `json:"searchKey"`
	PageNum   int                    `json:"pageNum"`
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

type TimelineVO struct {
	Id        int64                `json:"id"`
	PrId      int64                `json:"prId"`
	Action    pullrequestmd.Action `json:"action"`
	Account   string               `json:"account"`
	Created   string               `json:"created"`
	AvatarUrl string               `json:"avatarUrl"`
	Name      string               `json:"name"`
}

type ReviewVO struct {
	Id           int64  `json:"id"`
	Reviewer     string `json:"reviewer"`
	ReviewStatus string `json:"reviewStatus"`
	Updated      string `json:"updated"`
	AvatarUrl    string `json:"avatarUrl"`
	Name         string `json:"name"`
}

type AddCommentReqVO struct {
	PrId      int64  `json:"prId"`
	ReplyFrom int64  `json:"replyFrom"`
	Comment   string `json:"comment"`
	HasReply  bool   `json:"hasReply"`
}

type ProtectedBranchCfgVO struct {
	PushOption              branch.PushOption `json:"pushOption"`
	ReviewCountWhenCreatePr int               `json:"reviewCountWhenCreatePr"`
	ReviewerList            []util.User       `json:"reviewerList"`
}

type CanMergePullRequestRespVO struct {
	StatusChange       bool                 `json:"statusChange"`
	CanMerge           bool                 `json:"canMerge"`
	IsProtectedBranch  bool                 `json:"isProtectedBranch"`
	ProtectedBranchCfg ProtectedBranchCfgVO `json:"protectedBranchCfg"`
	ReviewCount        int                  `json:"reviewCount"`
	GitCanMerge        bool                 `json:"gitCanMerge"`
	GitConflictFiles   []string             `json:"gitConflictFiles"`
	GitCommitCount     int                  `json:"gitCommitCount"`
}

type CanReviewPullRequestRespVO struct {
	StatusChange      bool     `json:"statusChange"`
	CanReview         bool     `json:"canReview"`
	IsProtectedBranch bool     `json:"isProtectedBranch"`
	ReviewerList      []string `json:"reviewerList"`
	IsInReviewerList  bool     `json:"isInReviewerList"`
	HasAgree          bool     `json:"hasAgree"`
}

type StatusChangeVO struct {
	StatusChange bool `json:"statusChange"`
}
