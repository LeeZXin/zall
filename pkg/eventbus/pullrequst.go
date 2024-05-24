package eventbus

import "time"

const PullRequestEventTopic = "pull-request"

type PullRequestEvent struct {
	PrId      int64
	PrTitle   string
	Action    string
	RepoId    int64
	RepoName  string
	RepoPath  string
	Account   string
	EventTime time.Time
	Ref       string
}
