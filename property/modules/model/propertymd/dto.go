package propertymd

type InsertEtcdNodeReqDTO struct {
	AppId     string
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
	Endpoints string
	Username  string
	Password  string
	Creator   string
}

type PageHistoryReqDTO struct {
	FileId   int64
	PageNum  int
	PageSize int
}

type ListEtcdNodeReqDTO struct {
	AppId string   `json:"appId"`
	Env   string   `json:"env"`
	Cols  []string `json:"cols"`
}
