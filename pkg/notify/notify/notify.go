package notify

import (
	"bytes"
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/notify/feishu"
	"github.com/LeeZXin/zall/pkg/notify/wework"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/pingcap/errors"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	sendExecutor     *executor.Executor
	initExecutorOnce = sync.Once{}
)

func initSendExecutor() {
	initExecutorOnce.Do(func() {
		sendExecutor, _ = executor.NewExecutor(10, 1024, time.Minute, executor.AbortStrategy)
	})
}

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

// SendNotification 发送通知
func SendNotification(cfg Cfg, params any) error {
	if !cfg.IsValid() {
		return errors.New("invalid cfg")
	}
	notificationTpl, err := template.New("").Parse(cfg.Template)
	if err != nil {
		// 模板错误
		return err
	}
	mb := new(bytes.Buffer)
	err = notificationTpl.Execute(mb, params)
	if err != nil {
		// 参数错误
		return err
	}
	initSendExecutor()
	switch cfg.NotifyType {
	case Wework:
		var msg wework.Message
		err := json.Unmarshal(mb.Bytes(), &msg)
		if err != nil {
			return err
		}
		if err = msg.IsValid(); err != nil {
			return err
		}
		return sendExecutor.Execute(func() {
			wework.SendMessage(cfg.Url, msg)
		})
	case Feishu:
		var msg feishu.Message
		err := json.Unmarshal(mb.Bytes(), &msg)
		if err != nil {
			return err
		}
		if err = msg.IsValid(); err != nil {
			return err
		}
		return sendExecutor.Execute(func() {
			feishu.SendMessage(cfg.Url, cfg.FeishuSignKey, msg)
		})
	}
	return nil
}
