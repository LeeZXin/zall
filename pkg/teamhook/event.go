package teamhook

import "github.com/LeeZXin/zall/pkg/event"

type Events struct {
	ProtectedBranch event.ProtectedBranchEventCfg `json:"protectedBranch"`
	GitPush         event.GitPushEventCfg         `json:"gitPush"`
	PullRequest     event.PullRequestEventCfg     `json:"pullRequest"`
	GitRepo         event.GitRepoEventCfg         `json:"gitRepo"`
	GitWorkflow     event.GitWorkflowEventCfg     `json:"gitWorkflow"`
	Team            event.TeamEventCfg            `json:"team"`
	TeamRole        event.TeamRoleEventCfg        `json:"teamRole"`
	TeamUser        event.TeamUserEventCfg        `json:"teamUser"`
	App             event.AppEventCfg             `json:"app"`
	EnvRelated      map[string]EnvRelatedCfg      `json:"envRelated"`
}

type EnvRelatedCfg struct {
	AppSource             event.AppSourceEventCfg          `json:"appSource"`
	AppPropertyFile       event.AppPropertyFileEventCfg    `json:"appPropertyFile"`
	AppPropertyVersion    event.AppPropertyVersionEventCfg `json:"appPropertyVersion"`
	AppDeployPipeline     event.AppDeployPipelineEventCfg  `json:"appDeployPipeline"`
	AppDeployPipelineVars event.AppPipelineVarsEventCfg    `json:"appDeployPipelineVars"`
	AppDeployPlan         event.AppDeployPlanEventCfg      `json:"appDeployPlan"`
	AppProduct            event.AppProductEventCfg         `json:"appProduct"`
	AppDiscovery          event.AppDiscoveryEventCfg       `json:"appDiscovery"`
	AppDeployService      event.AppDeployServiceEventCfg   `json:"appDeployService"`
	AppPromScrape         event.AppPromScrapeEventCfg      `json:"appPromScrape"`
	Timer                 event.TimerEventCfg              `json:"timer"`
	TimerTask             event.TimerTaskEventCfg          `json:"timerTask"`
}
