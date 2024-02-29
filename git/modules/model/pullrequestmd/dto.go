package pullrequestmd

type InsertPullRequestReqDTO struct {
	RepoId   int64
	Target   string
	Head     string
	CreateBy string
	PrStatus PrStatus
}

type InsertReviewReqDTO struct {
	PrId      int64
	ReviewMsg string
	Status    ReviewStatus
	Reviewer  string
}

type UpdateReviewReqDTO struct {
	Id     int64
	Status ReviewStatus
}

type ExistsOpenStatusPrByRepoIdAndRefReqDTO struct {
	RepoId int64
	Head   string
	Target string
}
