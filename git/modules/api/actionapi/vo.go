package actionapi

type InsertActionReqVO struct {
	Name          string `json:"name"`
	TeamId        int64  `json:"teamId"`
	NodeId        int64  `json:"nodeId"`
	ActionContent string `json:"actionContent"`
}

type DeleteActionReqVO struct {
	Id int64 `json:"id"`
}

type ListActionReqVO struct {
	TeamId int64 `json:"teamId"`
}

type UpdateActionReqVO struct {
	Id            int64  `json:"id"`
	NodeId        int64  `json:"nodeId"`
	ActionContent string `json:"actionContent"`
}

type ActionVO struct {
	Id            int64  `json:"id"`
	Aid           string `json:"aid"`
	NodeId        int64  `json:"nodeId"`
	ActionContent string `json:"actionContent"`
	Created       string `json:"created"`
}

type TriggerActionReqVO struct {
	Id int64 `json:"id"`
}

type InsertNodeReqVO struct {
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
	Token    string `json:"token"`
}

type UpdateNodeReqVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
	Token    string `json:"token"`
}

type DeleteNodeReqVO struct {
	Id int64 `json:"id"`
}

type NodeVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	HttpHost string `json:"httpHost"`
	Token    string `json:"token"`
}

type SimpleNodeVO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ListTaskReqVO struct {
	ActionId int64 `json:"actionId"`
	Cursor   int64 `json:"cursor"`
	Limit    int   `json:"limit"`
}

type TaskVO struct {
	TaskStatus    string `json:"taskStatus"`
	TriggerType   string `json:"triggerType"`
	ActionContent string `json:"actionContent"`
	Operator      string `json:"operator"`
	Created       string `json:"created"`
}

type StepVO struct {
	JobName    string `json:"jobName"`
	StepName   string `json:"stepName"`
	StepIndex  int    `json:"stepIndex"`
	LogContent string `json:"logContent"`
	StepStatus string `json:"stepStatus"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}

type GetTaskStepsReqVO struct {
	TaskId int64 `json:"taskId"`
}
