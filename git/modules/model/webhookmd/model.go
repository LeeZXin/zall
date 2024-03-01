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

type Webhook struct {
	Id          int64     `json:"id" xorm:"pk autoincr"`
	RepoId      int64     `json:"repoId"`
	HookUrl     string    `json:"hookUrl"`
	HttpHeaders string    `json:"httpHeaders"`
	HookType    HookType  `json:"hookType"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

func (*Webhook) TableName() string {
	return WebhookTableName
}

func (h *Webhook) ToWebhookDTO() WebhookDTO {
	dto := WebhookDTO{
		Id:       h.Id,
		RepoId:   h.RepoId,
		HookUrl:  h.HookUrl,
		HookType: h.HookType,
	}
	if h.HttpHeaders != "" {
		header := make(map[string]string)
		_ = json.Unmarshal([]byte(h.HttpHeaders), &header)
		dto.HttpHeaders = header
	} else {
		dto.HttpHeaders = make(map[string]string)
	}
	return dto
}
