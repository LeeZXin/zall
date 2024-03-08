package appmd

type InsertAppReqDTO struct {
	AppId  string
	TeamId int64
	Name   string
}

type UpdateAppReqDTO struct {
	AppId string
	Name  string
}

type ListAppReqDTO struct {
	AppId  string
	TeamId int64
	Cursor int64
	Limit  int
}
