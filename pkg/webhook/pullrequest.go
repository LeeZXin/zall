package webhook

import "github.com/LeeZXin/zall/pkg/git"

type PullRequestActionHook struct {
	PrId       int64    `json:"prId"`
	RepoId     int64    `json:"repoId"`
	RepoName   string   `json:"repoName"`
	TargetRef  string   `json:"targetRef"`
	HeadRef    string   `json:"headRef"`
	EventTime  int64    `json:"eventTime"`
	ActionType string   `json:"actionType"`
	Operator   git.User `json:"operator"`
}
