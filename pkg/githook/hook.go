package githook

const (
	ApiPreReceiveUrl  = "api/v1/git/hook/pre-receive"
	ApiPostReceiveUrl = "api/v1/git/hook/post-receive"
)

type RevInfo struct {
	OldCommitId string `json:"oldCommitId"`
	NewCommitId string `json:"newCommitId"`
	Ref         string `json:"ref"`
}

type Opts struct {
	RevInfoList                  []RevInfo `json:"revInfoList"`
	RepoId                       int64     `json:"repoId"`
	PrId                         int64     `json:"prId"`
	PusherAccount                string    `json:"pusherAccount"`
	PusherEmail                  string    `json:"pusherEmail"`
	ObjectDirectory              string    `json:"objectDirectory"`
	AlternativeObjectDirectories string    `json:"alternativeObjectDirectories"`
	QuarantinePath               string    `json:"quarantinePath"`
}
