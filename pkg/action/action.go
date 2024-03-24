package action

import (
	"context"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/logger"
	"strings"
)

func TriggerAction(nodeHost, aid, operator string) {
	url := "http://" + strings.TrimSuffix(nodeHost, "/") +
		"/actionAgent/execute/" + aid + "?token=" + getActionToken() + "&operator=" + operator
	err := httputil.Get(context.Background(), httpClient, url, nil, nil)
	if err != nil {
		logger.Logger.Error(err)
	}
}

func getActionToken() string {
	cfg, b := cfgsrv.Inner.GetGitCfg(context.Background())
	if b {
		return cfg.ActionToken
	}
	return ""
}
