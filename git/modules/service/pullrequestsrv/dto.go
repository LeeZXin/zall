package pullrequestsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/util"
	"time"
)

type SubmitPullRequestReqDTO struct {
	RepoId     int64               `json:"repoId"`
	Target     string              `json:"target"`
	TargetType git.RefType         `json:"targetType"`
	Head       string              `json:"head"`
	HeadType   git.RefType         `json:"headType"`
	Title      string              `json:"title"`
	Comment    string              `json:"comment"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *SubmitPullRequestReqDTO) IsValid() error {
	if !r.TargetType.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.HeadType.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	if !pullrequestmd.IsPrTitleValid(r.Title) {
		return util.InvalidArgsError()
	}
	return nil
}

type ClosePullRequestReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ClosePullRequestReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type MergePullRequestReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *MergePullRequestReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AgreeReviewPullRequestReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AgreeReviewPullRequestReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListPullRequestReqDTO struct {
	RepoId    int64                  `json:"repoId"`
	Status    pullrequestmd.PrStatus `json:"status"`
	SearchKey string                 `json:"searchKey"`
	Operator  apisession.UserInfo    `json:"operator"`
	PageNum   int                    `json:"pageNum"`
}

func (r *ListPullRequestReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.SearchKey) > 255 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Status.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type PullRequestDTO struct {
	Id             int64
	RepoId         int64
	Target         string
	TargetType     git.RefType
	TargetCommitId string
	Head           string
	HeadType       git.RefType
	HeadCommitId   string
	PrStatus       pullrequestmd.PrStatus
	CreateBy       util.User
	PrTitle        string
	PrIndex        int
	CommentCount   int
	Created        time.Time
}

type GetPullRequestReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Index    int                 `json:"index"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetPullRequestReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Index <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListTimelineReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListTimelineReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TimelineDTO struct {
	Id        int64
	PrId      int64
	Action    pullrequestmd.Action
	Account   string
	AvatarUrl string
	Name      string
	Created   time.Time
}

type AddCommentReqDTO struct {
	PrId      int64               `json:"prId"`
	ReplyFrom int64               `json:"replyFrom"`
	Comment   string              `json:"comment"`
	HasReply  bool                `json:"hasReply"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *AddCommentReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Comment == "" || len(r.Comment) > 1024 {
		return util.InvalidArgsError()
	}
	if r.HasReply && r.ReplyFrom <= 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteCommentReqDTO struct {
	CommentId int64               `json:"commentId"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *DeleteCommentReqDTO) IsValid() error {
	if r.CommentId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CanMergePullRequestReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CanMergePullRequestReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListReviewReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListReviewReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CanReviewPullRequestReqDTO struct {
	PrId     int64               `json:"prId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CanReviewPullRequestReqDTO) IsValid() error {
	if r.PrId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ProtectedBranchCfgDTO struct {
	PushOption              branch.PushOption
	ReviewCountWhenCreatePr int
	ReviewerList            []util.User
}

type CanMergePullRequestRespDTO struct {
	CanMerge           bool
	IsProtectedBranch  bool
	ProtectedBranchCfg ProtectedBranchCfgDTO
	ReviewCount        int
	GitCanMerge        bool
	GitConflictFiles   []string
	GitCommitCount     int
}

type CanReviewPullRequestRespDTO struct {
	CanReview         bool
	IsProtectedBranch bool
	ReviewerList      []string
	IsInReviewerList  bool
	HasAgree          bool
}

type ReviewDTO struct {
	Id           int64
	Reviewer     string
	AvatarUrl    string
	Name         string
	ReviewStatus pullrequestmd.ReviewStatus
	Updated      time.Time
}
