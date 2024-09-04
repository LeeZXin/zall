package teamhook

import "github.com/LeeZXin/zall/pkg/event"

type Events struct {
	ProtectedBranch   event.ProtectedBranchEventCfg   `json:"protectedBranch"`
	GitPush           event.GitPushEventCfg           `json:"gitPush"`
	PullRequest       event.PullRequestEventCfg       `json:"pullRequest"`
	GitRepo           event.GitRepoEventCfg           `json:"gitRepo"`
	GitWorkflow       event.GitWorkflowEventCfg       `json:"gitWorkflow"`
	GitWorkflowVars   event.GitWorkflowVarsEventCfg   `json:"gitWorkflowVars"`
	GitWebhook        event.GitWebhookEventCfg        `json:"gitWebhook"`
	Team              event.TeamEventCfg              `json:"team"`
	TeamRole          event.TeamRoleEventCfg          `json:"teamRole"`
	TeamUser          event.TeamUserEventCfg          `json:"teamUser"`
	App               event.AppEventCfg               `json:"app"`
	WeworkAccessToken event.WeworkAccessTokenEventCfg `json:"weworkAccessToken"`
	FeishuAccessToken event.FeishuAccessTokenEventCfg `json:"feishuAccessToken"`
	NotifyTpl         event.NotifyTplEventCfg         `json:"notifyTpl"`
	EnvRelated        map[string]EnvRelatedCfg        `json:"envRelated"`
}

type EnvRelatedCfg struct {
	AppSource             event.AppSourceEventCfg          `json:"appSource"`
	AppPropertyFile       event.AppPropertyFileEventCfg    `json:"appPropertyFile"`
	AppPropertyVersion    event.AppPropertyVersionEventCfg `json:"appPropertyVersion"`
	AppDeployPipeline     event.AppDeployPipelineEventCfg  `json:"appDeployPipeline"`
	AppDeployPipelineVars event.AppPipelineVarsEventCfg    `json:"appDeployPipelineVars"`
	AppDeployPlan         event.AppDeployPlanEventCfg      `json:"appDeployPlan"`
	AppArtifact           event.AppArtifactEventCfg        `json:"appArtifact"`
	AppDiscovery          event.AppDiscoveryEventCfg       `json:"appDiscovery"`
	AppDeployService      event.AppDeployServiceEventCfg   `json:"appDeployService"`
	AppPromScrape         event.AppPromScrapeEventCfg      `json:"appPromScrape"`
	AppAlertConfig        event.AppAlertConfigEventCfg     `json:"appAlertConfig"`
	Timer                 event.TimerEventCfg              `json:"timer"`
	TimerTask             event.TimerTaskEventCfg          `json:"timerTask"`
}
