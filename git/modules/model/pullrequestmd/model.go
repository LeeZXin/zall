package pullrequestmd

import (
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/i18n"
	"time"
)

const (
	PullRequestTableName = "zgit_pull_request"
	ReviewTableName      = "zgit_pull_request_review"
	TimelineTableName    = "zgit_pull_request_timeline"
)

type PrStatus int

const (
	PrAllStatus PrStatus = iota
	PrOpenStatus
	PrClosedStatus
	PrMergedStatus
)

func (s PrStatus) Int() int {
	return int(s)
}

func (s PrStatus) IsValid() bool {
	switch s {
	case PrAllStatus, PrOpenStatus, PrClosedStatus, PrMergedStatus:
		return true
	default:
		return false
	}
}

func (s PrStatus) Readable() string {
	switch s {
	case PrOpenStatus:
		return i18n.GetByKey(i18n.PullRequestOpenStatus)
	case PrClosedStatus:
		return i18n.GetByKey(i18n.PullRequestClosedStatus)
	case PrMergedStatus:
		return i18n.GetByKey(i18n.PullRequestMergedStatus)
	default:
		return i18n.GetByKey(i18n.PullRequestUnknownStatus)
	}
}

type ReviewStatus int

const (
	AgreeMergeStatus ReviewStatus = iota + 1
	DisagreeMergeStatus
)

func (s ReviewStatus) Int() int {
	return int(s)
}

func (s ReviewStatus) Readable() string {
	switch s {
	case AgreeMergeStatus:
		return i18n.GetByKey(i18n.PullRequestAgreeMergeStatus)
	case DisagreeMergeStatus:
		return i18n.GetByKey(i18n.PullRequestDisagreeMergeStatus)
	default:
		return i18n.GetByKey(i18n.PullRequestUnknownReviewStatus)
	}
}

func (s ReviewStatus) IsValid() bool {
	switch s {
	case AgreeMergeStatus, DisagreeMergeStatus:
		return true
	default:
		return false
	}
}

type PullRequest struct {
	Id             int64       `json:"id" xorm:"pk autoincr"`
	RepoId         int64       `json:"repoId"`
	Target         string      `json:"target"`
	TargetType     git.RefType `json:"targetType"`
	TargetCommitId string      `json:"targetCommitId"`
	Head           string      `json:"head"`
	HeadType       git.RefType `json:"headType"`
	HeadCommitId   string      `json:"headCommitId"`
	PrStatus       PrStatus    `json:"prStatus"`
	CreateBy       string      `json:"createBy"`
	CloseBy        string      `json:"closeBy"`
	MergeBy        string      `json:"mergeBy"`
	PrTitle        string      `json:"prTitle"`
	CommentCount   int         `json:"commentCount"`
	Created        time.Time   `json:"created" xorm:"created"`
	Closed         *time.Time  `json:"closed"`
	Merged         *time.Time  `json:"merged"`
	Updated        time.Time   `json:"updated" xorm:"updated"`
}

func (*PullRequest) TableName() string {
	return PullRequestTableName
}

type Review struct {
	Id           int64        `json:"id" xorm:"pk autoincr"`
	PrId         int64        `json:"prId"`
	Reviewer     string       `json:"reviewer"`
	ReviewMsg    string       `json:"reviewMsg"`
	ReviewStatus ReviewStatus `json:"reviewStatus"`
	Created      time.Time    `json:"created" xorm:"created"`
	Updated      time.Time    `json:"updated" xorm:"updated"`
}

func (*Review) TableName() string {
	return ReviewTableName
}

type Timeline struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	PrId    int64     `json:"prId"`
	Action  *Action   `json:"action"`
	Account string    `json:"account"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*Timeline) TableName() string {
	return TimelineTableName
}
