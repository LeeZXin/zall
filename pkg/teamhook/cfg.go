package teamhook

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
