package webhook

import (
	"context"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
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

func initWebhook() {
	initOnce.Do(func() {
		httpClient = httputil.NewHttpClient()
		apiExecutor, _ = executor.NewExecutor(20, 10240, time.Minute, executor.CallerRunsStrategy)
	})
}

func TriggerGitHook(url string, headers map[string]string, req GitReceiveHook) {
	initWebhook()
	_ = apiExecutor.Execute(func() {
		err := httputil.Post(context.Background(), httpClient, url, headers, req, nil)
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}

func TriggerPrHook(url string, headers map[string]string, req PullRequestActionHook) {
	initWebhook()
	_ = apiExecutor.Execute(func() {
		err := httputil.Post(context.Background(), httpClient, url, headers, req, nil)
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}
