package propertymd

type InsertEtcdNodeReqDTO struct {
	NodeId    string
	Endpoints string
	Username  string
	Password  string
	Env       string
}

type UpdateEtcdNodeReqDTO struct {
	NodeId    string
	Endpoints string
	Username  string
	Password  string
	Env       string
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
	Env         string
	Creator     string
}

type InsertDeployReqDTO struct {
	ContentId    int64
	Content      string
	Version      string
	NodeId       string
	ContentAppId string
	ContentName  string
	Endpoints    string
	Username     string
	Password     string
	Creator      string
	Env          string
}

type PageHistoryReqDTO struct {
	FileId   int64
	PageNum  int
	PageSize int
}

type ListDeployReqDTO struct {
	ContentId int64
	Version   string
	NodeId    string
	Cursor    int64
	Limit     int
	Env       string
}

type InsertAuthReqDTO struct {
	AppId    string
	Username string
	Password string
	Env      string
}
