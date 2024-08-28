package event

type PullRequestEventAction string

func (a PullRequestEventAction) GetI18nValue() string {
	return "pullRequestEvent." + string(a)
}

const (
	PullRequestEventSubmitAction        PullRequestEventAction = "submit"
	PullRequestEventCloseAction         PullRequestEventAction = "close"
	PullRequestEventMergeAction         PullRequestEventAction = "merge"
	PullRequestEventReviewAction        PullRequestEventAction = "review"
	PullRequestEventAddCommentAction    PullRequestEventAction = "addComment"
	PullRequestEventDeleteCommentAction PullRequestEventAction = "deleteComment"
)

type RepoEventAction string

func (a RepoEventAction) GetI18nValue() string {
	return "repoEvent." + string(a)
}

const (
	RepoEventCreateAction             RepoEventAction = "create"
	RepoEventUpdateAction             RepoEventAction = "update"
	RepoEventDeleteTemporarilyAction  RepoEventAction = "deleteTemporarily"
	RepoEventDeletePermanentlyAction  RepoEventAction = "deletePermanently"
	RepoEventArchivedAction           RepoEventAction = "archived"
	RepoEventUnArchivedAction         RepoEventAction = "unArchived"
	RepoEventRecoverFromRecycleAction RepoEventAction = "recoverFromRecycle"
)

type ProtectedBranchEventAction string

func (a ProtectedBranchEventAction) GetI18nValue() string {
	return "protectedBranchEvent." + string(a)
}

const (
	ProtectedBranchEventCreateAction ProtectedBranchEventAction = "create"
	ProtectedBranchEventUpdateAction ProtectedBranchEventAction = "update"
	ProtectedBranchEventDeleteAction ProtectedBranchEventAction = "delete"
)

type GitWorkflowEventAction string

func (a GitWorkflowEventAction) GetI18nValue() string {
	return "gitWorkflowEvent." + string(a)
}

const (
	GitWorkflowEventCreateAction  GitWorkflowEventAction = "create"
	GitWorkflowEventUpdateAction  GitWorkflowEventAction = "update"
	GitWorkflowEventDeleteAction  GitWorkflowEventAction = "delete"
	GitWorkflowEventTriggerAction GitWorkflowEventAction = "trigger"
	GitWorkflowEventKillAction    GitWorkflowEventAction = "kill"
)

type GitWorkflowVarsEventAction string

func (a GitWorkflowVarsEventAction) GetI18nValue() string {
	return "gitWorkflowVarsEvent." + string(a)
}

const (
	GitWorkflowVarsEventCreateAction GitWorkflowVarsEventAction = "create"
	GitWorkflowVarsEventUpdateAction GitWorkflowVarsEventAction = "update"
	GitWorkflowVarsEventDeleteAction GitWorkflowVarsEventAction = "delete"
)

type GitWebhookEventAction string

func (a GitWebhookEventAction) GetI18nValue() string {
	return "gitWebhookEvent." + string(a)
}

const (
	GitWebhookEventCreateAction GitWebhookEventAction = "create"
	GitWebhookEventUpdateAction GitWebhookEventAction = "update"
	GitWebhookEventDeleteAction GitWebhookEventAction = "delete"
)

type GitPushEventAction string

func (a GitPushEventAction) GetI18nValue() string {
	return "gitPushEvent." + string(a)
}

const (
	GitPushEventCommitAction GitPushEventAction = "commit"
	GitPushEventDeleteAction GitPushEventAction = "delete"
)
