package action

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"strings"
)

var (
	httpClient = httputil.NewHttpClient()
)

const (
	ManualTriggerType = iota + 1
	SysTriggerType
)

func SysTriggerAction(yamlContent, nodeHost string, args map[string]string, actionId int64) {
	url := "http://" + strings.TrimSuffix(nodeHost, "/") + "/actionAgent/execute"
	err := httputil.Post(context.Background(), httpClient, url, map[string]string{
		"Authorization": getActionToken(),
	}, Hook{
		ActionYaml:  yamlContent,
		TriggerType: SysTriggerType,
		ActionId:    actionId,
		Args:        args,
	}, nil)
	if err != nil {
		logger.Logger.Error(err)
	}
}

func ManualTriggerAction(yamlContent, nodeHost string, args map[string]string, actionId int64) {
	url := "http://" + strings.TrimSuffix(nodeHost, "/") + "/actionAgent/execute"
	err := httputil.Post(context.Background(), httpClient, url, map[string]string{
		"Authorization": getActionToken(),
	}, Hook{
		ActionYaml:  yamlContent,
		TriggerType: ManualTriggerType,
		ActionId:    actionId,
		Args:        args,
	}, nil)
	if err != nil {
		logger.Logger.Error(err)
	}
}

type Hook struct {
	ActionId    int64             `json:"actionId"`
	ActionYaml  string            `json:"actionYaml"`
	TriggerType int               `json:"triggerType"`
	Args        map[string]string `json:"args"`
}

func getActionToken() string {
	cfg, b := cfgsrv.Inner.GetGitCfg(context.Background())
	if b {
		return cfg.ActionToken
	}
	return ""
}
