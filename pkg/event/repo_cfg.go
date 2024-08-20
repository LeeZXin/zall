package event

type PullRequestEventCfg struct {
	Submit bool `json:"submit"`
	Close  bool `json:"close"`
	Merge  bool `json:"merge"`
	Review bool `json:"review"`
}

type GitPushEventCfg struct {
	Commit bool `json:"commit"`
	Delete bool `json:"delete"`
}

type ProtectedBranchEventCfg struct {
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type GitRepoEventCfg struct {
	Create             bool `json:"create"`
	DeleteTemporarily  bool `json:"deleteTemporarily"`
	DeletePermanently  bool `json:"deletePermanently"`
	Archived           bool `json:"archived"`
	UnArchived         bool `json:"unArchived"`
	RecoverFromRecycle bool `json:"recoverFromRecycle"`
}

type GitWorkflowEventCfg struct {
	Create  bool `json:"create"`
	Update  bool `json:"update"`
	Delete  bool `json:"delete"`
	Trigger bool `json:"trigger"`
}
