package webhookmd

import (
	"github.com/LeeZXin/zall/pkg/webhook"
	"time"
)

const (
	WebhookTableName = "zgit_webhook"
)

type Webhook struct {
	Id      int64           `json:"id" xorm:"pk autoincr"`
	RepoId  int64           `json:"repoId"`
	HookUrl string          `json:"hookUrl"`
	Secret  string          `json:"secret"`
	Events  *webhook.Events `json:"events"`
	Created time.Time       `json:"created" xorm:"created"`
	Updated time.Time       `json:"updated" xorm:"updated"`
}

func (*Webhook) TableName() string {
	return WebhookTableName
}
