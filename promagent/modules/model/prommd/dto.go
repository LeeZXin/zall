package prommd

type InsertScrapeReqDTO struct {
	Endpoint   string
	AppId      string
	Target     string
	TargetType TargetType
	Env        string
}

type UpdateScrapeByIdReqDTO struct {
	Id         int64
	Endpoint   string
	Target     string
	TargetType TargetType
}

type GetAllScrapeReqDTO struct {
	Endpoint string
	Env      string
	Cols     []string
}

type ListScrapeReqDTO struct {
	AppId    string
	Env      string
	Endpoint string
	PageNum  int
	PageSize int
}
