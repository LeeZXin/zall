package idapi

type InsertGeneratorReqVO struct {
	CurrentId int64  `json:"currentId"`
	BizName   string `json:"bizName"`
}
