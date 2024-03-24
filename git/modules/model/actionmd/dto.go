package actionmd

type InsertTaskReqDTO struct {
	ActionId      int64
	TaskStatus    TaskStatus
	TriggerType   TriggerType
	ActionContent string
	Operator      string
}

type InsertStepReqDTO struct {
	TaskId     int64
	JobName    string
	StepName   string
	StepIndex  int
	StepStatus StepStatus
}

type InsertNodeReqDTO struct {
	Name     string
	HttpHost string
	Token    string
}

type UpdateNodeReqDTO struct {
	Id       int64
	Name     string
	Token    string
	HttpHost string
}

type InsertActionReqDTO struct {
	Aid     string
	TeamId  int64
	Name    string
	Content string
	NodeId  int64
}

type UpdateActionReqDTO struct {
	Id      int64
	Name    string
	Content string
	NodeId  int64
}

type GetTaskReqDTO struct {
	ActionId int64
	Cursor   int64
	Limit    int
}
