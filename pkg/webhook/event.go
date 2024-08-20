package webhook

type Events struct {
	ProtectedBranch bool `json:"protectedBranch"`
	GitPush         bool `json:"gitPush"`
	PullRequest     bool `json:"pullRequest"`
	GitRepo         bool `json:"gitRepo"`
}
