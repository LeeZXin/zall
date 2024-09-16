package discoverymd

type InsertEtcdNodeReqDTO struct {
	Name      string
	Endpoints string
	Username  string
	Password  string
	Env       string
}

type UpdateEtcdNodeReqDTO struct {
	Id        int64
	Name      string
	Endpoints string
	Username  string
	Password  string
}

type ListEtcdNodeReqDTO struct {
	Env  string
	Cols []string
}

type InsertAppEtcdNodeBindReqDTO struct {
	NodeId int64
	AppId  string
	Env    string
}
