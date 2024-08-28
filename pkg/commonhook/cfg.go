package commonhook

import (
	"net/url"
	"strings"
)

type HookType int

func (t HookType) IsValid() bool {
	switch t {
	case WebhookType, NotifyType:
		return true
	default:
		return false
	}
}

const (
	WebhookType HookType = iota + 1
	NotifyType
)

type Cfg struct {
	HookUrl     string `json:"hookUrl"`
	Secret      string `json:"secret"`
	NotifyTplId int64  `json:"notifyTplId"`
}

type TypeAndCfg struct {
	HookCfg  Cfg      `json:"hookCfg"`
	HookType HookType `json:"hookType"`
}

func (c *TypeAndCfg) IsValid() bool {
	switch c.HookType {
	case NotifyType:
		if c.HookCfg.NotifyTplId <= 0 {
			return false
		}
	case WebhookType:
		parsedUrl, err := url.Parse(c.HookCfg.HookUrl)
		if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
			return false
		}
		if c.HookCfg.Secret == "" {
			return false
		}
	}
	return true
}
