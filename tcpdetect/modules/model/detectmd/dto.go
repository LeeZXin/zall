package detectmd

type InsertLogReqDTO struct {
	DetectId int64
	Ip       string
	Port     int
	Valid    bool
}

type InsertDetectReqDTO struct {
	Ip            string
	Port          int
	Name          string
	HeartbeatTime int64
	Enabled       bool
}

type UpdateDetectReqDTO struct {
	Id   int64
	Ip   string
	Port int
	Name string
}

type ListDetectReqDTO struct {
	Name   string
	Cursor int64
	Limit  int
}

type ListLogReqDTO struct {
	Id     int64
	Cursor int64
	Limit  int
}
