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
	TeamMember      event.TeamMemberEventCfg      `json:"teamMember"`
	App             event.AppEventCfg             `json:"app"`
	EnvRelated      map[string]EnvRelatedCfg      `json:"envRelated"`
}

type EnvRelatedCfg struct {
	AppSource          event.AppSourceEventCfg          `json:"appSource"`
	AppPropertyFile    event.AppPropertyFileEventCfg    `json:"appPropertyFile"`
	AppPropertyVersion event.AppPropertyVersionEventCfg `json:"appPropertyVersion"`
	AppPipeline        event.AppPipelineEventCfg        `json:"appPipeline"`
	AppPipelineVars    event.AppPipelineVarsEventCfg    `json:"appPipelineVars"`
	AppDeployPlan      event.AppDeployPlanEventCfg      `json:"appDeployPlan"`
	AppProduct         event.AppProductEventCfg         `json:"appProduct"`
	Timer              event.TimerEventCfg              `json:"timer"`
	TimerTask          event.TimerTaskEventCfg          `json:"timerTask"`
}
