package event

const (
	PullRequestTopic     = "pull-request"
	GitRepoTopic         = "git-repo"
	ProtectedBranchTopic = "protected-branch"
	GitPushTopic         = "git-push"
	GitWorkflowTopic     = "git-workflow"
	GitWorkflowVarsTopic = "git-workflow-vars"
	GitWebhookTopic      = "git-webhook"

	AppTopic       = "app"
	TeamTopic      = "team"
	TeamRoleTopic  = "team-role"
	TeamUserTopic  = "team-user"
	TimerTopic     = "timer"
	TimerTaskTopic = "timer-task"

	AppPropertyFileTopic       = "app-property-file"
	AppPropertyVersionTopic    = "app-property-version"
	AppSourceTopic             = "app-source"
	AppDiscoveryTopic          = "app-discovery"
	AppDeployServiceTopic      = "app-deploy-service"
	AppDeployPipelineTopic     = "app-deploy-pipeline"
	AppDeployPipelineVarsTopic = "app-deploy-pipeline-vars"
	AppDeployPlanTopic         = "app-deploy-plan"
	AppArtifactTopic           = "app-artifact"
	AppPromScrapeTopic         = "app-prom-scrape"
)
