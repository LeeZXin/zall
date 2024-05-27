package oplogmd

import "time"

type InsertOpLogReqDTO struct {
	RepoId   int64
	Operator string
	Content  string
	ReqBody  string
	Created  time.Time
}

type PageOpLogReqDTO struct {
	RepoId             int64
	PageNum            int
	PageSize           int
	Account            string
	BeginTime, EndTime time.Time
}
