package eventbus

import "time"

const GitRepoEventTopic = "git-repo"

type GitRepoEvent struct {
	RepoId    int64     `json:"repoId"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Operator  string    `json:"operator"`
	Action    string    `json:"action"`
	EventTime time.Time `json:"eventTime"`
}
