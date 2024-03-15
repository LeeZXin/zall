package propmd

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

type InsertPropContentReqDTO struct {
	AppId string
	Name  string
	Env   string
}

type InsertHistoryReqDTO struct {
	ContentId int64
	Content   string
	Version   string
	Env       string
	Creator   string
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

type ListHistoryReqDTO struct {
	ContentId int64
	Version   string
	Cursor    int64
	Limit     int
	Env       string
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
