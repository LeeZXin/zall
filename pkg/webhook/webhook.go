package webhook

import (
	"context"
	"crypto/tls"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf/logger"
	"net/http"
	"sync"
	"time"
)

var (
	apiExecutor *executor.Executor
	httpClient  *http.Client
	initOnce    = sync.Once{}
)

func initTrigger() {
	initOnce.Do(func() {
		httpClient = &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:   true,
				MaxIdleConns:        0, // 禁用连接池
				MaxIdleConnsPerHost: 0, // 禁用连接池
				IdleConnTimeout:     0, // 禁用连接池
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 3 * time.Second,
		}
		apiExecutor, _ = executor.NewExecutor(20, 1024, time.Minute, executor.AbortStrategy)
	})
}

func TriggerWebhook(url, secret string, req event.Event) {
	initTrigger()
	apiExecutor.Execute(func() {
		err := Post(context.Background(), url, secret, req)
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}
