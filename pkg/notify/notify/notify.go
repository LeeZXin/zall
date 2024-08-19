package notify

import (
	"net/url"
	"strings"
)

type Type string

const (
	Wework Type = "wework"
	Feishu Type = "feishu"
)

type Cfg struct {
	Url           string `json:"url"`
	NotifyType    Type   `json:"notifyType"`
	Template      string `json:"template"`
	FeishuSignKey string `json:"feishuSignKey"`
}

func (c *Cfg) IsValid() bool {
	parsedUrl, err := url.Parse(c.Url)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return false
	}
	switch c.NotifyType {
	case Wework, Feishu:
		return true
	default:
		return false
	}
}
