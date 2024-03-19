package gitactionmd

type InsertTaskReqDTO struct {
	ActionId      int64
	TaskName      string
	TaskStatus    TaskStatus
	TriggerType   int
	ActionContent string
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
}

type UpdateNodeReqDTO struct {
	Id       int64
	Name     string
	HttpHost string
}

type InsertActionReqDTO struct {
	RepoId     int64
	Name       string
	Content    string
	NodeId     int64
	PushBranch string
}

type UpdateActionReqDTO struct {
	Id         int64
	Name       string
	Content    string
	NodeId     int64
	PushBranch string
}
