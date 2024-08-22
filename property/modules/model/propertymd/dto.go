package propertymd

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

type InsertFileReqDTO struct {
	AppId string
	Name  string
	Env   string
}

type InsertHistoryReqDTO struct {
	FileId      int64
	Content     string
	Version     string
	LastVersion string
	Creator     string
}

type InsertDeployReqDTO struct {
	NodeName  string
	HistoryId int64
	FileId    int64
	AppId     string
	Endpoints string
	Username  string
	Password  string
	Creator   string
}

type ListHistoryReqDTO struct {
	FileId   int64
	PageNum  int
	PageSize int
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
