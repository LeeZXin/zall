package deploymd

type InsertPlanReqDTO struct {
	Name           string
	PlanStatus     PlanStatus
	AppId          string
	PipelineId     int64
	PipelineName   string
	ProductVersion string
	Creator        string
	Env            string
	PipelineConfig string
}

type InsertDeployLogReqDTO struct {
	ConfigId       int64
	AppId          string
	PlanId         int64
	PipelineConfig string
	ProductVersion string
	Env            string
	DeployOutput   string
	Operator       string
}

type InsertOpLogReqDTO struct {
	ConfigId       int64
	Operator       string
	ScriptOutput   string
	Env            string
	Op             Op
	ProductVersion string
}

type ListDeployLogReqDTO struct {
	ConfigId int64
	Cursor   int64
	Limit    int
	Env      string
}

type ListOpLogReqDTO struct {
	ConfigId int64
	Cursor   int64
	Limit    int
	Env      string
}

type ListPlanReqDTO struct {
	AppId    string
	PageNum  int
	PageSize int
	Env      string
}

type InsertDeployStageReqDTO struct {
	PlanId      int64
	AppId       string
	StageIndex  int
	Agent       string
	TaskId      string
	StageStatus StageStatus
}

type InsertPipelineReqDTO struct {
	AppId  string
	Name   string
	Config string
	Env    string
}

type UpdatePipelineReqDTO struct {
	PipelineId int64
	Name       string
	Config     string
}

type ListPipelineReqDTO struct {
	AppId string
	Env   string
	Cols  []string
}

type InsertServiceSourceReqDTO struct {
	Name   string
	Env    string
	Host   string
	ApiKey string
}

type UpdateServiceSourceReqDTO struct {
	Id     int64
	Name   string
	Host   string
	ApiKey string
}

type ListServiceSourceReqDTO struct {
	Env  string
	Cols []string
}

type InsertPipelineVarsReqDTO struct {
	AppId   string
	Env     string
	Name    string
	Content string
}

type UpdatePipelineVarsReqDTO struct {
	Id      int64
	Content string
}

type InsertAppServiceSourceBindReqDTO struct {
	SourceId int64
	AppId    string
	Env      string
}
