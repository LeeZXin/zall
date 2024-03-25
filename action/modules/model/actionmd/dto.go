package actionmd

type InsertTaskReqDTO struct {
	ActionId      int64
	TaskStatus    TaskStatus
	TriggerType   TriggerType
	ActionContent string
	Operator      string
	AgentUrl      string
	AgentToken    string
}

type InsertStepReqDTO struct {
	TaskId     int64
	JobName    string
	StepName   string
	StepIndex  int
	StepStatus StepStatus
}

type InsertActionReqDTO struct {
	Aid        string
	TeamId     int64
	Name       string
	Content    string
	AgentUrl   string
	AgentToken string
}

type UpdateActionReqDTO struct {
	Id         int64
	Name       string
	Content    string
	AgentUrl   string
	AgentToken string
}

type GetTaskReqDTO struct {
	ActionId int64
	Cursor   int64
	Limit    int
}
