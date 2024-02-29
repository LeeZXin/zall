package pullrequestapi

import "github.com/LeeZXin/zall/git/modules/model/pullrequestmd"

type SubmitPullRequestReqVO struct {
	RepoId int64  `json:"repoId"`
	Target string `json:"target"`
	Head   string `json:"head"`
}

type ClosePullRequestReqVO struct {
	Id int64 `json:"id"`
}

type MergePullRequestReqVO struct {
	Id int64 `json:"id"`
}

type ReviewPullRequestReqVO struct {
	Id        int64                      `json:"id"`
	Status    pullrequestmd.ReviewStatus `json:"status"`
	ReviewMsg string                     `json:"reviewMsg"`
}
