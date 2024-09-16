package event

import "github.com/LeeZXin/zall/pkg/branch"

type BaseRepo struct {
	RepoId   int64  `json:"repoId"`
	RepoPath string `json:"repoPath"`
	RepoName string `json:"repoName"`
}

type ProtectedBranchEvent struct {
	BaseTeam
	BaseRepo
	BaseEvent
	Action ProtectedBranchEventAction `json:"action"`
	Before *branch.ProtectedBranch    `json:"before,omitempty"`
	After  *branch.ProtectedBranch    `json:"after,omitempty"`
}

func (*ProtectedBranchEvent) EventType() string {
	return "protected-branch-event"
}

type GitRepoEvent struct {
	BaseTeam
	BaseRepo
	BaseEvent
	Action RepoEventAction `json:"action"`
}

func (*GitRepoEvent) EventType() string {
	return "git-repo-event"
}

type GitPushEvent struct {
	RefType     string             `json:"refType"`
	Ref         string             `json:"ref"`
	OldCommitId string             `json:"oldCommitId"`
	NewCommitId string             `json:"newCommitId"`
	Action      GitPushEventAction `json:"action"`
	BaseTeam
	BaseRepo
	BaseEvent
}

func (*GitPushEvent) EventType() string {
	return "git-push-event"
}

type PullRequestEvent struct {
	PrId          int64                  `json:"prId"`
	PrTitle       string                 `json:"prTitle"`
	PrIndex       int                    `json:"prIndex"`
	TargetRef     string                 `json:"targetRef"`
	TargetRefType string                 `json:"targetRefType"`
	HeadRef       string                 `json:"headRef"`
	HeadRefType   string                 `json:"headRefType"`
	Action        PullRequestEventAction `json:"action"`
	BaseTeam
	BaseRepo
	BaseEvent
}

func (*PullRequestEvent) EventType() string {
	return "pull-request-event"
}

type PingEvent struct {
	EventTime int64 `json:"eventTime"`
}

func (*PingEvent) EventType() string {
	return "ping-event"
}

type GitWorkflowEvent struct {
	BaseTeam
	BaseRepo
	BaseEvent
	Action       GitWorkflowEventAction `json:"action"`
	WorkflowId   int64                  `json:"workflowId"`
	WorkflowName string                 `json:"workflowName"`
}

func (*GitWorkflowEvent) EventType() string {
	return "git-workflow-event"
}

type GitWorkflowVarsEvent struct {
	BaseTeam
	BaseRepo
	BaseEvent
	Action   GitWorkflowVarsEventAction `json:"action"`
	VarsId   int64                      `json:"varsId"`
	VarsName string                     `json:"varsName"`
}

func (*GitWorkflowVarsEvent) EventType() string {
	return "git-workflow-vars-event"
}

type GitWebhookEvent struct {
	BaseTeam
	BaseRepo
	BaseEvent
	Action     GitWebhookEventAction `json:"action"`
	WebhookId  int64                 `json:"webhookId"`
	WebhookUrl string                `json:"webhookUrl"`
}

func (*GitWebhookEvent) EventType() string {
	return "git-webhook-event"
}
