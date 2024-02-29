package action

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zsf-utils/executor"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"strings"
	"time"
)

var (
	apiExecutor *executor.Executor
	httpClient  = httputil.NewHttpClient()
)

func init() {
	apiExecutor, _ = executor.NewExecutor(10, 1024, time.Minute, executor.CallerRunsStrategy)
}

func TriggerActionHook(hook Webhook, instanceHost string) {
	apiExecutor.Execute(func() {
		url := strings.TrimSuffix(instanceHost, "/") + "/actions/git"
		err := httputil.Post(context.Background(), httpClient, url, map[string]string{
			"Authorization": getActionToken(),
		}, hook, nil)
		if err != nil {
			logger.Logger.Error(err)
		}
	})
}

func getActionToken() string {
	cfg, b := cfgsrv.Inner.GetGitCfg(context.Background())
	if b {
		return cfg.ActionToken
	}
	return ""
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
