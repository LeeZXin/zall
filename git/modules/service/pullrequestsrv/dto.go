package pullrequestsrv

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
)

type SubmitPullRequestReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Target   string              `json:"target"`
	Head     string              `json:"head"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *SubmitPullRequestReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	return nil
}

type ClosePullRequestReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ClosePullRequestReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type MergePullRequestReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *MergePullRequestReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ReviewPullRequestReqDTO struct {
	Id        int64                      `json:"id"`
	Status    pullrequestmd.ReviewStatus `json:"status"`
	ReviewMsg string                     `json:"reviewMsg"`
	Operator  apisession.UserInfo        `json:"operator"`
}

func (r *ReviewPullRequestReqDTO) IsValid() error {
	if len(r.ReviewMsg) > 255 {
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
