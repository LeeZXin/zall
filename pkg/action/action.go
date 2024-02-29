package action

import (
	"context"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"strings"
	"time"
)

var (
	apiExecutor *executor.Executor
	httpClient  = httputil.NewHttpClient()
	authHeader  = map[string]string{
		"Authorization": static.GetString("actions.token"),
	}
)

func init() {
	apiExecutor, _ = executor.NewExecutor(10, 1024, time.Minute, executor.CallerRunsStrategy)
}

func TriggerActionHook(hook Webhook, instanceHost string) {
	apiExecutor.Execute(func() {
		url := strings.TrimSuffix(instanceHost, "/") + "/actions/git"
		err := httputil.Post(context.Background(), httpClient, url, authHeader, hook, nil)
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}

type Webhook struct {
	RepoId    int64    `json:"repoId"`
	RepoName  string   `json:"repoName"`
	Ref       string   `json:"ref"`
	EventTime int64    `json:"eventTime"`
	Operator  git.User `json:"operator"`
	// yaml配置信息
	YamlContent string `json:"yamlContent"`
	// 触发类型
	TriggerType int
}
