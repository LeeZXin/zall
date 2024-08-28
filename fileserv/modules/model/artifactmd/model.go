package artifactmd

import "time"

const (
	ArtifactTableName = "zfile_artifact"
)

type Artifact struct {
	Id      int64     `json:"id" xorm:"pk autoincr"`
	AppId   string    `json:"appId"`
	Env     string    `json:"env"`
	Name    string    `json:"name"`
	Creator string    `json:"creator"`
	Created time.Time `json:"created" xorm:"created"`
}

func (*Artifact) TableName() string {
	return ArtifactTableName
}
