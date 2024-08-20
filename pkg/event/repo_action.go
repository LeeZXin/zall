package event

type PullRequestAction string

const (
	PrSubmitAction PullRequestAction = "submit"
	PrCloseAction  PullRequestAction = "close"
	PrMergeAction  PullRequestAction = "merge"
	PrReviewAction PullRequestAction = "review"
)

type GitRepoAction string

const (
	RepoCreateAction             GitRepoAction = "create"
	RepoDeleteTemporarilyAction  GitRepoAction = "deleteTemporarily"
	RepoDeletePermanentlyAction  GitRepoAction = "deletePermanently"
	RepoArchivedAction           GitRepoAction = "archived"
	RepoUnArchivedAction         GitRepoAction = "unArchived"
	RepoRecoverFromRecycleAction GitRepoAction = "recoverFromRecycle"
)

type ProtectedBranchAction string

const (
	PbCreateAction ProtectedBranchAction = "create"
	PbUpdateAction ProtectedBranchAction = "update"
	PbDeleteAction ProtectedBranchAction = "delete"
)
