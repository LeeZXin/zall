package webhook

import (
	"github.com/LeeZXin/zall/pkg/git"
)

type GitReceiveHook struct {
	RepoId    int64    `json:"repoId"`
	RepoName  string   `json:"repoName"`
	IsCreated bool     `json:"isCreated"`
	IsDeleted bool     `json:"isDeleted"`
	Ref       string   `json:"ref"`
	EventTime int64    `json:"eventTime"`
	Operator  git.User `json:"operator"`
	IsTagPush bool     `json:"isTagPush"`
}
