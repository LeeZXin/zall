package prommd

type InsertScrapeReqDTO struct {
	ServerUrl  string
	AppId      string
	Target     string
	TargetType TargetType
	Env        string
}

type UpdateScrapeByIdReqDTO struct {
	Id         int64
	ServerUrl  string
	Target     string
	TargetType TargetType
	Env        string
}
