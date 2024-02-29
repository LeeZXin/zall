package actionsrv

type TriggerGitActionReqDTO struct {
	RepoId    string
	CorpId    string
	Ref       string
	IsTagPush bool
	IsDeleted bool
	IsCreated bool
}
