package artifactmd

type InsertArtifactReqDTO struct {
	AppId   string
	Name    string
	Creator string
	Env     string
}

type GetArtifactReqDTO struct {
	AppId string
	Name  string
	Env   string
}

type ListArtifactReqDTO struct {
	AppId    string
	Env      string
	PageNum  int
	PageSize int
}
