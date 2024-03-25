package webhookmd

import (
	"encoding/json"
	"time"
)

const (
	WebhookTableName = "zgit_webhook"
)

type HookType int

const (
	PushHook HookType = iota + 1
	TagHook
	PullRequestHook
)

func (t HookType) Int() int {
	return int(t)
}

func (t HookType) IsValid() bool {
	switch t {
	case PushHook, TagHook, PullRequestHook:
		return true
	default:
		return false
	}
}

func (t HookType) Readable() string {
	switch t {
	case PushHook:
		return "push"
	case TagHook:
		return "tag"
	case PullRequestHook:
		return "pullRequest"
	default:
		return "unknown"
	}
}

type HttpHeaders map[string]string

func (h *HttpHeaders) FromDB(content []byte) error {
	if h == nil {
		*h = make(map[string]string)
	}
	return json.Unmarshal(content, h)
}

func (h *HttpHeaders) ToDB() ([]byte, error) {
	return json.Marshal(h)
}

type Webhook struct {
	Id          int64       `json:"id" xorm:"pk autoincr"`
	RepoId      int64       `json:"repoId"`
	HookUrl     string      `json:"hookUrl"`
	HttpHeaders HttpHeaders `json:"httpHeaders"`
	HookType    HookType    `json:"hookType"`
	WildBranch  string      `json:"wildBranch"`
	WildTag     string      `json:"wildTag"`
	Created     time.Time   `json:"created" xorm:"created"`
	Updated     time.Time   `json:"updated" xorm:"updated"`
}

func (*Webhook) TableName() string {
	return WebhookTableName
}
