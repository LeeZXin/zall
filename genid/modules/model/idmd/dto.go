package idmd

type UpdateCurrentIdReqDTO struct {
	BizName   string
	CurrentId int64
	Version   int64
}

type InsertGeneratorReqDTO struct {
	BizName   string
	CurrentId int64
}
