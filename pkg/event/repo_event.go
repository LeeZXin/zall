package event

import "github.com/LeeZXin/zall/pkg/branch"

type BaseRepo struct {
	TeamId   int64  `json:"teamId"`
	RepoId   int64  `json:"repoId"`
	RepoPath string `json:"repoPath"`
	RepoName string `json:"repoName"`
}

type ProtectedBranchEvent struct {
	BaseRepo
	BaseEvent
	Action ProtectedBranchAction   `json:"action"`
	Before *branch.ProtectedBranch `json:"before,omitempty"`
	After  *branch.ProtectedBranch `json:"after,omitempty"`
}

func (*ProtectedBranchEvent) EventType() string {
	return "protected-branch-event"
}

type GitRepoEvent struct {
	BaseRepo
	BaseEvent
	Action GitRepoAction `json:"action"`
}

func (*GitRepoEvent) EventType() string {
	return "git-repo-event"
}

type GitPushEvent struct {
	RefType     string `json:"refType"`
	Ref         string `json:"ref"`
	OldCommitId string `json:"oldCommitId"`
	NewCommitId string `json:"newCommitId"`
	BaseRepo
	BaseEvent
}

func (*GitPushEvent) EventType() string {
	return "git-push-event"
}

type PullRequestEvent struct {
	PrId    int64             `json:"prId"`
	PrTitle string            `json:"prTitle"`
	Ref     string            `json:"ref"`
	RefType string            `json:"refType"`
	Action  PullRequestAction `json:"action"`
	BaseRepo
	BaseEvent
}

func (*PullRequestEvent) EventType() string {
	return "pull-request-event"
}

type PingEvent struct {
	EventTime int64 `json:"eventTime"`
}

func (*PingEvent) EventType() string {
	return "ping-event"
}
