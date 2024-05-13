package webhook

import (
	"encoding/json"
)

type Event int

const (
	/*
		ProtectedBranchEvent
		保护分支创建、删除、编辑
	*/
	ProtectedBranchEvent Event = iota + 1
	/*
		GitPushEvent
		git ssh或http push
	*/
	GitPushEvent
	/*
		PullRequestEvent
		合并请求创建、关闭、合并
	*/
	PullRequestEvent
	/*
		RepoEvent
		仓库创建、删除
	*/
	RepoEvent
	/*
		PingEvent ping
	*/
	PingEvent
)

func (e Event) String() string {
	switch e {
	case ProtectedBranchEvent:
		return "protected-branch-event"
	case GitPushEvent:
		return "git-push-event"
	case PullRequestEvent:
		return "pull-request-event"
	case RepoEvent:
		return "repo-event"
	case PingEvent:
		return "ping-event"
	default:
		return "unknown-event"
	}
}

func (e Event) IsValid() bool {
	switch e {
	case ProtectedBranchEvent, GitPushEvent, PullRequestEvent, RepoEvent:
		return true
	default:
		return false
	}
}

type Events []Event

func (es *Events) IsValid() bool {
	if len(*es) == 0 {
		return false
	}
	for _, e := range *es {
		if !e.IsValid() {
			return false
		}
	}
	return true
}

func (es *Events) Has(event Event) bool {
	for _, e := range *es {
		if e == event {
			return true
		}
	}
	return false
}

func (es *Events) FromDB(content []byte) error {
	if es == nil {
		*es = make([]Event, 0)
	}
	return json.Unmarshal(content, es)
}

func (es *Events) ToDB() ([]byte, error) {
	return json.Marshal(es)
}
